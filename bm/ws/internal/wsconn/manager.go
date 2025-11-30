package wsconn

import (
	"net/http"
	"sync"

	"ws/internal/global"
	callback "ws/internal/rpc/client"
	"ws/pkg/wscache"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	WebsocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) AddClient(wsid string, client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[wsid] = client

	// 通知业务系统该wsid对应的WebSocket连接开启
	if global.CallbackIsOpen {
		_, err := callback.Open(wsid)
		if err != nil {
			logx.Errorf("[%s]WebSocket连接开启后，回调业务系统的Open服务出错，错误信息为：%v", wsid, err)
		}
	}

	global.WebSocketCount.Inc()

	logx.Infof("新增一个[%s]WebSocket连接，目前连接数为[%d]", wsid, len(m.clients))
}

func (m *Manager) RemoveClient(wsid string) {
	m.Lock()
	defer m.Unlock()

	if cli, ok := m.clients[wsid]; ok {
		cli.conn.Close()
		delete(m.clients, wsid)

		// 删除Redis缓存
		_, err := wscache.GetInstance().SRem(global.WS_WSIDS, wsid)
		if err != nil {
			logx.Errorf("删除[%s]的[%s]缓存失败", global.WS_WSIDS, wsid)
		}

		global.WebSocketCount.Dec()
		logx.Infof("减少一个[%s]WebSocket连接，目前连接数为[%d]", wsid, len(m.clients))

		// 回调业务系统，通知业务系统该wsid对应的连接已经关闭
		if global.CallbackIsOpen {
			_, err = callback.Close(wsid)
			if err != nil {
				logx.Errorf("[%s]WebSocket连接关闭后，回调业务系统的Close服务出错，错误信息为：%v", wsid, err)
			}
		}
	}
}

func (m *Manager) GetClient(wsid string) *Client {
	client, ok := m.clients[wsid]
	if !ok {
		return nil
	}

	return client
}
