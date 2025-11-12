// 分布式唯一 ID 示例，使用 sony/sonyflake
package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"time"

	"github.com/sony/sonyflake"
)

// getMachineID 尝试从环境变量 MACHINE_ID 中读取机器ID（uint16），如果未提供则使用本机第一个非回环 IPv4 地址
// 的哈希值的低16位作为回退。sonyflake 需要返回一个 uint16 指定机器ID。
func getMachineIDFromEnvOrIP() (uint16, error) {
	if v := os.Getenv("MACHINE_ID"); v != "" {
		// 直接解析为整数
		var id uint64
		_, err := fmt.Sscanf(v, "%d", &id)
		if err == nil {
			return uint16(id), nil
		}
		// 如果失败，继续使用 IP 回退
	}

	// 查找本地网络地址
	ifaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if (iface.Flags & net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not ipv4
			}
			// 使用简单哈希将 IPv4 映射到 uint16
			h := fnv.New32a()
			_, _ = h.Write([]byte(ip))
			return uint16(h.Sum32() & 0xffff), nil
		}
	}
	// 最后回退为 0
	return 0, nil
}

func main() {
	var count int
	var msid uint
	var useEnvMachine bool

	flag.IntVar(&count, "n", 5, "生成 ID 的数量")
	flag.UintVar(&msid, "machine", 0, "可选的机器 ID（会覆盖 MACHINE_ID 环境变量）")
	flag.BoolVar(&useEnvMachine, "fromenv", false, "优先使用 MACHINE_ID 环境变量或根据本机 IP 计算机器 ID")
	flag.Parse()

	var machineIDFunc func() (uint16, error)
	if useEnvMachine {
		machineIDFunc = func() (uint16, error) {
			id, err := getMachineIDFromEnvOrIP()
			return id, err
		}
	} else {
		// 如果用户在命令行提供了 machine 参数，则使用该值
		machineIDFunc = func() (uint16, error) {
			return uint16(msid), nil
		}
	}

	st := sonyflake.Settings{
		StartTime: time.Now().Add(-time.Hour * 24), // 可选：设置较近的起始时间，默认2014-09-01
		MachineID: machineIDFunc,
	}
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		fmt.Println("failed to create sonyflake")
		os.Exit(2)
	}

	for i := 0; i < count; i++ {
		id, err := sf.NextID()
		if err != nil {
			fmt.Printf("generate id error: %v\n", err)
			continue
		}
		// 打印十进制和十六进制表示
		fmt.Printf("%3d: %d  (0x%x)\n", i+1, id, id)
	}
}
