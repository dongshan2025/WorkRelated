// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package notify

import (
	"net/http"

	"ws/internal/logic/notify"
	"ws/internal/svc"
	"ws/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetWsConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetWsConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("客户端端创建WebSocket连接时解析请求参数时报错，错误信息为：%v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notify.NewGetWsConnectLogic(r.Context(), svcCtx)
		l.GetWsConnect(w, r, req, svcCtx)
	}
}
