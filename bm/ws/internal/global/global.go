package global

import "github.com/prometheus/client_golang/prometheus"

const (
	WS_MACHINEID   = "WS:WS_MACHINEID"
	WS_WSIDS       = "WS:WSIDS"
	WS_NOTIFYTOPIC = "NOTIFY-TOPIC"
)

var MachineID string
var RegisterIP string
var ApiPort int
var RpcPort int
var KafkaIsOpen bool
var CallbackIsOpen bool

var WebSocketCount prometheus.Gauge

type NotifyMessage struct {
	Wsid    []string
	Message string
}
