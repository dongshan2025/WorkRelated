# UUID (sonyflake) — 分布式唯一 ID 生成示例

这是一个使用 `github.com/sony/sonyflake` 的简单示例程序，演示如何在分布式环境中生成全局唯一且可排序的 ID。

主要特点：
- 基于 Sonyflake（Snowflake 风格）：包含时间戳、机器ID、序列号，按时间可排序。
- 支持通过环境变量或命令行指定机器 ID（以避免冲突）。

## 要求
- Go 1.24+

## 文件
- `main.go`：示例实现（使用 sony/sonyflake，支持机器ID配置与 IP 作为回退）。

## 用法

在 PowerShell（Windows）中，进入 `UUID` 目录后：

```powershell
cd 'd:\Workspace\Golang\WorkRelated\UUID'
go get github.com/sony/sonyflake@latest
go mod tidy
go run . -n 3 -fromenv
```

命令行参数：
- `-n` 要生成的 ID 数量（默认为 5）。
- `-machine` 手动指定机器 ID（整数，最终使用 uint16）。
- `-fromenv` 优先使用 `MACHINE_ID` 环境变量；若未设置则尝试使用本机第一个非回环 IPv4 地址的哈希值（低16位）作为机器 ID 回退。

示例：
- 使用环境变量指定机器 ID 并生成 10 个 ID：

```powershell
$env:MACHINE_ID="12"
go run . -n 10 -fromenv
```

- 使用命令行指定机器 ID：

```powershell
go run . -n 5 -machine 3
```

示例输出（十进制与十六进制）：

```
  1: 144955196593518  (0x83d60300556e)
  2: 144955196659054  (0x83d60301556e)
  3: 144955196724590  (0x83d60302556e)
```

## 部署注意事项（重要）

1. 机器 ID 必须在你的集群/部署中唯一（范围 0..65535）。建议在容器/VM/实例创建时由部署系统显式注入 `MACHINE_ID` 环境变量或通过配置管理分配。依赖自动生成（如 IP 哈希）会存在冲突风险，尤其在 NAT、多网卡或容器短暂重启场景下。

2. 时间回拨（clock drift）可能导致序列重复或生成错误的时间顺序 ID。确保主机时间同步（NTP）。如果有严格的时间一致性需求，可考虑对时间回拨做额外保护或使用中心化分配器。

3. Sonyflake 的 StartTime 设置会影响可表示的时间范围：在极端长寿命系统中请选择合适起始时间或使用扩展方案。

4. 如果打算在多个数据中心部署，仍要保证机器 ID 的全局唯一性；可以将高位做数据中心编码，或在部署阶段分配不同区间的机器 ID。

## 进阶建议

- 将生成器封装成一个小型 HTTP 服务（例如 `GET /id` 返回一个新 ID），其他应用通过 HTTP 调用来获取 ID，从而避免在每个服务里自行处理机器 ID 的分配与存活问题。我可以为你添加一个简单的 HTTP 包装示例（需要你确认）。

- 如果需要更短的可读 ID，可将生成的 uint64 转为 Base62 编码输出。

## 下一步（可选）
- 我可以为你：
  - 添加一个 `README`（已完成）并提交；
  - 添加单元测试验证机器 ID 解析与 ID 唯一性（小规模）；
  - 实现一个简单的 HTTP API 服务封装 sonyflake，带 health-check 与 metrics（Prometheus）。

如果你想继续，请回复你希望我实现的下一项（例如：`添加 HTTP 服务示例`）。

---
作者：仓库自动示例（基于 sonyflake）
