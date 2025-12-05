package wsconn

import (
	"fmt"
	"time"
	"ws/internal/global"
	"ws/pkg/wscache"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	pongWait     = 20 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[string]*Client

type Client struct {
	wsid    string
	conn    *websocket.Conn
	manager *Manager

	Egress chan []byte
}

func NewClient(wsid string, conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		wsid:    wsid,
		conn:    conn,
		manager: manager,
		Egress:  make(chan []byte),
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		logx.Infof("[ReadMessage]准备关闭[%s]的WebSocket连接", c.wsid)
		c.manager.RemoveClient(c.wsid)
	}()

	c.conn.SetReadLimit(1024)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logx.Errorf("[ReadMessage]为[%s]的WebSocket连接设置超时时间时发生错误，错误信息为：%v", c.wsid, err)
		return
	}

	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("[%s]的WebSocket连接发生错误，错误信息为：%v", c.wsid, err)
			}
			break
		}

		fmt.Printf("收到客户端信息为：%s\n", string(msg))
	}
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		logx.Infof("[WriteMessage]准备关闭[%s]的WebSocket连接", c.wsid)
		c.manager.RemoveClient(c.wsid)
	}()

	for {
		select {
		case message, ok := <-c.Egress:
			// egress关闭情况
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// 记录关闭原因
					logx.Errorf("[WriteMessage][%s]的egress通道关闭，WebSocket发送消息的错误信息为：%v", c.wsid, err)
				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logx.Errorf("[WriteMessage][%s]写入WebSocket消息时发生错误，错误信息为：%v", c.wsid, err)
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// logx.Errorf("[WriteMessage][%s]发送Ping消息时发生错误，错误信息为：%v", c.wsid, err)
				// 删除Redis缓存
				_, err := wscache.GetInstance().SRem(global.WS_WSIDS, c.wsid)
				if err != nil {
					logx.Errorf("删除[%s]的[%s]缓存失败", global.WS_WSIDS, c.wsid)
				}
				return
			}
			fmt.Printf("[%s]-%s\n", c.wsid, "Ping")
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	fmt.Printf("[%s]-%s\n", c.wsid, "Pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
