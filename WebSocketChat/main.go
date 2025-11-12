package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// For demo purposes: simple token -> username map. In real apps replace with
// proper auth (JWT validation, session store, etc.).
var validTokens = map[string]string{
	"token-alice": "Alice",
	"token-bob":   "Bob",
}

// validateToken returns mapped username if token is valid.
func validateToken(token string) (string, bool) {
	if token == "" {
		return "", false
	}
	if u, ok := validTokens[token]; ok {
		return u, true
	}
	return "", false
}

// Simple gorilla/websocket chat server inspired by the canonical hub example.

var addr = flag.String("addr", ":8081", "http service address")

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	// rooms maps room name to set of clients
	rooms map[string]map[*Client]bool
	// users maps username to a set of clients (support multiple connections per username)
	users map[string]map[*Client]bool
	// inbound parsed messages
	inbound chan *Inbound
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		users:      make(map[string]map[*Client]bool),
		inbound:    make(chan *Inbound),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				// remove from rooms
				if c.room != "" {
					if clients, ok := h.rooms[c.room]; ok {
						delete(clients, c)
						if len(clients) == 0 {
							delete(h.rooms, c.room)
						}
					}
				}
				// remove client from username mapping set
				if c.username != "" {
					if set, ok := h.users[c.username]; ok {
						delete(set, c)
						if len(set) == 0 {
							delete(h.users, c.username)
						}
					}
				}
				close(c.send)
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// client unresponsive; remove
					delete(h.clients, c)
					close(c.send)
				}
			}
		case in := <-h.inbound:
			// handle parsed inbound message
			m := in.message
			sender := in.client
			// ensure timestamp
			if m.Timestamp == "" {
				m.Timestamp = time.Now().Format(time.RFC3339)
			}
			// Listing requests
			if m.Type == "list_rooms" {
				// collect room names
				var rooms []string
				for rn := range h.rooms {
					rooms = append(rooms, rn)
				}
				resp := Message{Type: "rooms", Rooms: rooms, Timestamp: m.Timestamp}
				if b, err := json.Marshal(resp); err == nil {
					select {
					case sender.send <- b:
					default:
					}
				}
				continue
			}
			if m.Type == "list_users" {
				var users []string
				for uname := range h.users {
					users = append(users, uname)
				}
				resp := Message{Type: "users", Users: users, Timestamp: m.Timestamp}
				if b, err := json.Marshal(resp); err == nil {
					select {
					case sender.send <- b:
					default:
					}
				}
				continue
			}
			if m.Type == "list_room_users" {
				if m.Room == "" {
					// room not provided -> send empty
					resp := Message{Type: "room_users", Users: []string{}, Timestamp: m.Timestamp}
					if b, err := json.Marshal(resp); err == nil {
						select {
						case sender.send <- b:
						default:
						}
					}
					continue
				}
				var users []string
				if clients, ok := h.rooms[m.Room]; ok {
					seen := make(map[string]bool)
					for c := range clients {
						if c.username != "" && !seen[c.username] {
							seen[c.username] = true
							users = append(users, c.username)
						}
					}
				}
				resp := Message{Type: "room_users", Users: users, Timestamp: m.Timestamp}
				if b, err := json.Marshal(resp); err == nil {
					select {
					case sender.send <- b:
					default:
					}
				}
				continue
			}

			// If this is a join message, register username/room
			if m.Type == "join" {
				// If token provided in join, validate and require it. If invalid -> disconnect.
				if m.Token != "" {
					if uname, ok := validateToken(m.Token); ok {
						sender.username = uname
						if _, ok := h.users[uname]; !ok {
							h.users[uname] = make(map[*Client]bool)
						}
						h.users[uname][sender] = true
					} else {
						// auth failed: notify and unregister
						nack := Message{Type: "system", Msg: "auth failed", Timestamp: m.Timestamp}
						if b, err := json.Marshal(nack); err == nil {
							select {
							case sender.send <- b:
							default:
							}
						}
						h.unregister <- sender
						continue
					}
				} else if m.Username != "" {
					// No token provided, accept username (backward compatible)
					sender.username = m.Username
					if _, ok := h.users[m.Username]; !ok {
						h.users[m.Username] = make(map[*Client]bool)
					}
					h.users[m.Username][sender] = true
				}

				if m.Room != "" {
					sender.room = m.Room
					if _, ok := h.rooms[m.Room]; !ok {
						h.rooms[m.Room] = make(map[*Client]bool)
					}
					h.rooms[m.Room][sender] = true
				}
				// ack join to the sender
				ack := Message{Type: "system", Msg: "joined", Username: sender.username, Room: sender.room, Timestamp: m.Timestamp}
				if b, err := json.Marshal(ack); err == nil {
					select {
					case sender.send <- b:
					default:
					}
				}
				continue
			}
			// build output JSON
			out, err := json.Marshal(m)
			if err != nil {
				// fallback: ignore
				continue
			}
			// private to username
			if m.Target != "" {
				if set, ok := h.users[m.Target]; ok {
					for dest := range set {
						select {
						case dest.send <- out:
						default:
						}
					}
				}
				continue
			}
			// room broadcast
			if m.Room != "" {
				if clients, ok := h.rooms[m.Room]; ok {
					for c := range clients {
						select {
						case c.send <- out:
						default:
						}
					}
				}
				continue
			}
			// global broadcast
			for c := range h.clients {
				select {
				case c.send <- out:
				default:
				}
			}
		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
	// username chosen by client (optional until join)
	username string
	// room joined by client (optional)
	room string
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for example purposes. In production lock this down.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Try to parse JSON message into our message structure.
		var m Message
		if err := json.Unmarshal(message, &m); err != nil {
			// fallback: treat as simple text message
			m = Message{Type: "message", Msg: string(message)}
		}
		// send parsed inbound to the hub for processing (join, room, private, broadcast)
		c.hub.inbound <- &Inbound{client: c, message: m}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Message represents the JSON message structure used by clients and server.
type Message struct {
	Type      string `json:"type"` // "join", "leave", "message", "system"
	Username  string `json:"username,omitempty"`
	Room      string `json:"room,omitempty"`
	Target    string `json:"target,omitempty"` // username for private
	Token     string `json:"token,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	// helper response fields
	Users []string `json:"users,omitempty"`
	Rooms []string `json:"rooms,omitempty"`
}

// Inbound is a wrapper for a parsed message and its sender client
type Inbound struct {
	client  *Client
	message Message
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Check token on handshake (query param). If provided, validate before upgrade.
	token := r.URL.Query().Get("token")
	var tokenUser string
	if token != "" {
		if u, ok := validateToken(token); ok {
			tokenUser = u
		} else {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	// register client first (adds to clients map in hub.run)
	client.hub.register <- client
	// if handshake carried a valid token, enqueue a join inbound so hub.run can map the username
	if tokenUser != "" {
		join := Message{Type: "join", Username: tokenUser}
		client.hub.inbound <- &Inbound{client: client, message: join}
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	addr := *addr
	fmt.Printf("server started at %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
