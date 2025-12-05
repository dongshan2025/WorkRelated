package wsgrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"ws/internal/config"
	"ws/internal/global"
	wsserver "ws/internal/rpc/server/proto"
	"ws/internal/wskafka/wsproducer"

	"ws/internal/svc"
	"ws/pkg/wscache"
	"ws/pkg/wsuuid"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

// 主要用于转发通知时使用
var svcCtx *svc.ServiceContext

type WsServer struct {
	wsserver.UnimplementedWsServerServer
}

func NewRpcWsServer(svc *svc.ServiceContext, config *config.Config) {
	// 赋值
	svcCtx = svc
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.RPC.Host, config.RPC.Port))
	if err != nil {
		logx.Errorf("启动WS-gRPC服务失败，错误信息为：%v", err)
		panic(err)
	}

	s := grpc.NewServer()
	wsserver.RegisterWsServerServer(s, &WsServer{})

	logx.Infof("WS-gRPC服务在[%s:%d]上启动中...", config.RPC.Host, config.RPC.Port)

	if err := s.Serve(lis); err != nil {
		logx.Errorf("启动WS-gRPC服务失败，错误信息为：%v", err)
		panic(err)
	}
}

func (s *WsServer) Preconnect(ctx context.Context, req *wsserver.PreconnectRequest) (*wsserver.PreconnectResponse, error) {
	logx.Infof("gRPC服务端[GetWsid]接收到来自[%s]的请求", req.Token)
	uuid, err := wsuuid.GetUUID()
	if err != nil {
		logx.Errorf("gRPC服务端[GetWsid]生成UUID失败，错误信息为：%v", err)
		return &wsserver.PreconnectResponse{Status: wsserver.WsServerStatus_WSSERVER_FAILED}, errors.New("生成WSID失败")
	}

	wsid := fmt.Sprintf("%s-%d-%s", global.MachineID, uuid, req.Token)
	// 将wsid写入Redis
	_, err = wscache.GetInstance().SAdd(global.WS_WSIDS, wsid)
	if err != nil {
		logx.Errorf("gRPC服务端[GetWsid]将WSID写入Redis失败，错误信息为：%v", err)
		return &wsserver.PreconnectResponse{Status: wsserver.WsServerStatus_WSSERVER_FAILED}, errors.New("将WSID写入Redis失败")
	}

	logx.Infof("gRPC服务端[GetWsid]生成WSID[%s]成功", wsid)

	return &wsserver.PreconnectResponse{
		Status:    wsserver.WsServerStatus_WSSERVER_SUCCESS,
		Wsid:      wsid,
		MachineId: global.MachineID,
		WsAddr:    fmt.Sprintf("%s/v1/ws?wsid=%s", global.RegisterWsAddr, wsid),
		GrpcAddr:  fmt.Sprintf("%s:%d", global.RegisterRpcIp, global.RpcPort),
	}, nil
}

func (s *WsServer) Send(ctx context.Context, req *wsserver.SendRequest) (*wsserver.SendResponse, error) {
	logx.Infof("gRPC服务端[SendNotify]接收到来自[%v]的请求，消息为：%s", req.Wsid, req.Message)
	if global.KafkaIsOpen {
		message := global.NotifyMessage{
			Wsid:    req.Wsid,
			Message: req.Message,
		}

		msgData, err := json.Marshal(message)
		if err != nil {
			logx.Errorf("[Send]JSON序列化时发生错误，错误信息为：%v", err)
			return &wsserver.SendResponse{Status: wsserver.WsServerStatus_WSSERVER_FAILED}, err
		}

		wsproducer.SyncSendMessage(global.WS_NOTIFYTOPIC, string(msgData))
	} else {
		for _, wsid := range req.Wsid {
			cli := svcCtx.Manager.GetClient(wsid)
			if cli != nil {
				cli.Egress <- []byte(req.Message)
				logx.Infof("gRPC服务端[SendNotify]处理了[%s]发送的消息[%s]", wsid, req.Message)
			} else {
				logx.Errorf("gRPC服务端[SendNotify]没有发现[%s]的客户端，消息[%s]发送失败", wsid, req.Message)
			}
		}

		logx.Infof("gRPC服务端[SendNotify]处理[%v]的消息[%s]完成", req.Wsid, req.Message)
	}

	resp := &wsserver.SendResponse{Status: wsserver.WsServerStatus_WSSERVER_SUCCESS, SenderWsid: req.SenderWsid, SenderStatus: "1"}
	cli := svcCtx.Manager.GetClient(req.SenderWsid)
	if cli == nil {
		resp.SenderStatus = "2"
	}
	offlineWsid := make([]string, 0, len(req.Wsid))
	for _, wsid := range req.Wsid {
		c := svcCtx.Manager.GetClient(wsid)
		if c == nil {
			offlineWsid = append(offlineWsid, wsid)
		}
	}
	resp.OfflineWisd = offlineWsid

	return resp, nil
}
