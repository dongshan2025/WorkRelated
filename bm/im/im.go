package main

import callback "im/rpc/server"

func main() {
	// wsid, err := notify.GetWsid("xxx-yyy-zzz")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("收到服务端生成的wsid为：%s", wsid)

	// wsids := []string{"xxx-yyy-001", "xxx-yyy-002"}

	// status, err := notify.SendNotify(wsids, "通知")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("收到服务端发送通知的响应：%d", status)

	callback.NewCallbackServer()
}
