package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/valyala/fasthttp"
)

func main() {
	httpServer := &fasthttp.Server{
		Handler: fastHTTPHandler,
		Name:    "kingkong",
	}

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = httpServer.Serve(ln)
	if err != nil {
		fmt.Println(err)
	}
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/createEventSub":
		createEventSubHandler(ctx)
	// case "/postJson":
	// 	postJsonHandler(ctx)
	// case "/postForm":
	// 	postFormHandler(ctx)
	// case "/upload":
	// 	uploadHandler(ctx)
	// case "/uploadMulti":
	// 	uploadMultiHandler(ctx)
	// case "/bar":
	// 	barHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func createEventSubHandler(ctx *fasthttp.RequestCtx) {
	req := CallbackInfo{}
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ChannelId: %s, MsgId: %s, MsgTimestamp: %d, SubscribeId: %s\n\n", req.ChannelId, req.MsgId, req.MsgTimestamp, req.SubscribeId)

	var userEvent UserEvent
	for _, con := range req.Contents {
		if con.Event == "ChannelEvent" {
			fmt.Printf("ChannelEvent ChannelId: %s, EventTag: %s, Timestamp: %d, ChannelProfile: %s\n\n", con.ChannelEvent.ChannelId, con.ChannelEvent.EventTag, con.ChannelEvent.Timestamp, con.ChannelEvent.ChannelProfile)
		} else if con.Event == "UserEvent" {
			fmt.Printf("UserEvent UserId: %s, EventTag: %s, SessionId: %s, Timestamp: %d, ChannelProfile: %s, Role: %d, TerminalType: %d, UserType: %d, CurrentMedias: %s, Reason: %d\n\n", con.UserEvent.UserId, con.UserEvent.EventTag, con.UserEvent.SessionId, con.UserEvent.Timestamp, con.UserEvent.ChannelProfile, con.UserEvent.Role, con.UserEvent.TerminalType, con.UserEvent.UserType, con.UserEvent.CurrentMedias, con.UserEvent.Reason)
			userEvent = con.UserEvent
		}
	}
	fmt.Println("-------------------------------------------------------------------------------------------")

	str, err := json.Marshal(userEvent)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(string(str))

	if userEvent.EventTag == "Join" {
		fmt.Println(userEvent.EventTag)
	} else if userEvent.EventTag == "Leave" {
		fmt.Println(userEvent.EventTag)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	fmt.Fprintf(ctx, "Create Event Sub")
}

type ChannelEvent struct {
	ChannelId      string
	EventTag       string // 事件类型，Open:会议开始 Close:会议结束
	Timestamp      int64
	ChannelProfile string
}

type UserEvent struct {
	UserId         string // 用户ID
	EventTag       string // 事件类型，Join:入会 Leave:离会 PublishVideo:开始推视频流 PublishAudio:开始推音频流 PublishScreen:开始屏幕共享 UnpublishVideo:停止推视频流 UnpublishAudio:停止推音频流 UnpublishScreen:停止屏幕共享 Roleupdate:角色切换
	SessionId      string // 产生该事件的SessionID
	Timestamp      int64  // 事件发生时的Unix时间戳
	ChannelProfile string
	US             int64
	Reason         int // 入会、离会原因（仅Join事件有），1:正常入会、离会 2:重连入会 3:跨频道转推 4:超时离会 5:用户启用新的会话，当前会话被挤下线 6:被踢出 7:频道解散
	Role           int // 角色类型，1:主播 2:观众
	TerminalType   int
	UserType       int
	CurrentMedias  string // 推流类型，1:音频 2:视频 3:屏幕共享
}

type Content struct {
	Event        string // 订阅的事件，频道内用户事件
	ChannelEvent ChannelEvent
	UserEvent    UserEvent
}

type CallbackInfo struct {
	MsgId        string // 消息ID
	MsgTimestamp int64  // 消息发送时的Unix时间戳
	SubscribeId  string // 订阅ID
	AppId        string // 产生该消息的appid
	ChannelId    string // 产生该消息的频道
	Contents     []Content
}
