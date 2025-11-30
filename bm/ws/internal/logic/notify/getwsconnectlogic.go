// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package notify

import (
	"context"
	"net/http"

	"ws/internal/global"
	"ws/internal/svc"
	"ws/internal/types"
	"ws/internal/wsconn"
	"ws/pkg/wscache"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWsConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWsConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWsConnectLogic {
	return &GetWsConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWsConnectLogic) GetWsConnect(w http.ResponseWriter, r *http.Request, req types.GetWsConnectReq, svcCtx *svc.ServiceContext) {
	logx.Infof("API服务端[GetWsConnect]接收到来自[%s]的创建WebSocket请求", req.Wsid)

	isExist, err := wscache.GetInstance().SIsMember(global.WS_WSIDS, req.Wsid)
	if err != nil {
		logx.Errorf("API服务端[GetWsConnect]查询WSID[%s]信息时失败，错误信息为：%v", req.Wsid, err)
		return
	}

	if !isExist {
		logx.Errorf("API服务端[GetWsConnect]查询WSID[%s]不存在", req.Wsid)
		return
	}

	conn, err := wsconn.WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Errorf("API服务端[GetWsConnect]为WSID[%s]升级WebSocket时失败，错误信息为：%v", req.Wsid, err)
		return
	}

	client := wsconn.NewClient(req.Wsid, conn, svcCtx.Manager)
	svcCtx.Manager.AddClient(req.Wsid, client)
	logx.Infof("API服务端[GetWsConnect]为WSID[%s]升级WebSocket成功", req.Wsid)

	go client.ReadMessage()
	go client.WriteMessage()
}
