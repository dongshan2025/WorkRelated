# Viper2 配置与日志说明

这是 `Viper2` 模块的说明，包含如何使用 `config` 封装读取配置以及如何通过配置控制 zap 日志行为的详细示例。

目录结构（相关）：
- `config/config.go`：viper 封装（Init、WatchConfig、GetString/Int/Bool、Unmarshal 等）。
- `config/zap_adapter.go`：从 `log.*` 配置构建 zap logger 的适配器（支持 encoder 配置、时间格式、采样等）。
- `config.yaml`：示例配置文件。

## 使用方式

在程序启动时先初始化配置，然后用 `SetupZapFromConfig()` 从配置创建 zap logger 并注入：

```go
if err := config.Init("./config.yaml", "MYAPP"); err != nil { /* handle */ }
l, err := config.SetupZapFromConfig()
if err == nil && l != nil {
    defer l.Sync()
}
```

或者如果你已有一个 zap logger，可以直接注入：

```go
l, _ := zap.NewProduction()
defer l.Sync()
config.SetZapLogger(l)
```

`config` 包会暴露 `GetString/GetInt/GetBool/Unmarshal` 等方法，方便读取配置。

## 日志配置详解（`log` 节点）

你可以在 `config.yaml`（或环境变量）中配置 `log` 下的字段，下面列出受支持的键与说明：

- `log.level` (string)
  - 日志级别：`debug` `info` `warn` `error` `dpanic` `panic` `fatal`。

- `log.development` (bool)
  - 是否使用开发模式的默认 zap 配置（会更偏向可读输出）。

- `log.encoding` (string)
  - 输出编码：`json` 或 `console`。

- `log.output_paths` ([]string)
  - 输出路径列表，例如 `["stdout", "logs/app.log"]`。

- `log.error_output_paths` ([]string)
  - 错误输出路径列表，例如 `["stderr"]`。

### Encoder 配置（`log.encoder.*`）
这些项会映射到 zap 的 `EncoderConfig`，以便微调输出字段和编码器：

- `log.encoder.time_key` (string)
  - 时间字段名，默认 `ts` 或 zap 默认值。

- `log.encoder.level_key`, `log.encoder.name_key`, `log.encoder.caller_key`, `log.encoder.message_key`, `log.encoder.stacktrace_key`, `log.encoder.line_ending`
  - 对应 zap 的 `LevelKey, NameKey, CallerKey, MessageKey, StacktraceKey, LineEnding`。

- `log.encoder.level_encode` (string)
  - 级别编码器：`lowercase`, `lowercase_color`, `capital`, `capital_color`。

- `log.encoder.caller_encode` (string)
  - Caller 编码：`short` 或 `full`（映射到 zap 的 Short/Full CallerEncoder）。

- `log.encoder.name_encode` (string)
  - Name 编码：`short` 或 `full`。
  - `short` 实现为取点分隔名称的最后一段（例如 `service.handler` -> `handler`）。

- `log.encoder.duration_encode` (string)
  - Duration 编码：`string`, `nanos`, `millis`, `seconds`（对应 zapcore 的 String/Nanos/Millis/Seconds Duration 编码器）。

- `log.encoder.time_encode` (string) 或 `log.encoder.time_format` (string)
  - `time_encode` 支持：`iso8601`, `millis`, `nanos`, `seconds`（使用 zapcore 的内置编码器）。
  - `time_format`：如果设置，会以该 Go 时间格式（例如 `2006-01-02T15:04:05.000Z07:00`）格式化时间（优先于 `time_encode`）。

### 采样（Sampling）
- `log.sampling.enabled` (bool)
- `log.sampling.initial` (int)
- `log.sampling.thereafter` (int)

当 `sampling.enabled` 为 `true` 时，会使用 zap 的采样配置（`Initial`、`Thereafter`）来降低高吞吐量日志的输出频率。

## 配置示例（摘自 `config.yaml`）

```yaml
log:
  level: info
  development: false
  encoding: console
  output_paths:
    - stdout
  error_output_paths:
    - stderr
  encoder:
    time_key: ts
    level_key: level
    name_key: logger
    caller_key: caller
    message_key: msg
    stacktrace_key: stack
    line_ending: "\n"
    time_format: "2006-01-02T15:04:05.000Z07:00"
    level_encode: lowercase
    caller_encode: short
    name_encode: short
    duration_encode: string
  sampling:
    enabled: true
    initial: 100
    thereafter: 100
```

## Server ID 示例

如果你想为每个实例提供一个可识别的 `server.id`，可以在 `config.yaml` 中设置默认值：

```yaml
server:
  id: srv-01
```

注意优先级关系（从高到低）：
- 命令行 flag（优先）
- 环境变量（BindEnv，示例：`SERVER_ID=srv-01`）
- 配置文件（`config.yaml` 中的 `server.id`）
- 程序默认值

在本项目中我们已经通过 `vv.BindEnv("server.id", "SERVER_ID")` 将环境变量 `SERVER_ID` 绑定到 viper 的 `server.id`，因此你可以通过运行时环境注入（例如 Docker 环境或 CI 注入）来覆盖 `config.yaml` 中的默认值。

`validate` 子命令会打印出 `server.id` 的当前生效值，`serve` 启动时也会在启动日志中包含 `server.id` 字段，便于运维识别不同实例。


## 环境变量覆盖

`config` 支持环境变量覆盖。默认会把 `.` 替换为 `_`，并可用前缀（`Init(configFile, envPrefix)`）来限制。例：

```
TEST_APP_NAME=envapp
TEST_LOG_LEVEL=debug
```

上面会覆盖 `app.name` 和 `log.level`。

## 进阶说明与扩展

- `name_encode` 的 `short` 是一个简单实现（取最后段）。如果你需要更复杂的策略（如截断、前缀或正则抓取），我可以替换为更复杂的实现并将策略暴露为配置。
- `duration_encode` 目前支持 zapcore 内置编码器；如果你需要自定义格式（例如 `ms` 并保留 3 位小数），我可以添加 `log.encoder.duration_format` 支持。

如果你希望我把 README 的这些段落同步到仓库顶层 README 或生成更完整的文档（例如 markdown 页面），告诉我要把文档放哪里即可。

## 环境变量示例（PowerShell）

下面的示例展示如何在没有 `config.yaml` 的情况下，通过环境变量完整配置日志与应用设置。注意：`Init` 中传入的 `envPrefix` 会作为前缀（下面示例使用 `MYAPP`），并且键中的 `.` 会被替换为 `_`。

设置示例：

```powershell
# 应用配置
$env:MYAPP_APP_NAME = 'env-app'
$env:MYAPP_APP_PORT = '9090'
$env:MYAPP_DEBUG = 'true'

# 日志配置示例
$env:MYAPP_LOG_LEVEL = 'debug'
$env:MYAPP_LOG_ENCODING = 'console'
$env:MYAPP_LOG_ENCODER_TIME_FORMAT = '2006-01-02T15:04:05.000Z07:00'
$env:MYAPP_LOG_SAMPLING_ENABLED = 'true'
$env:MYAPP_LOG_SAMPLING_INITIAL = '100'
$env:MYAPP_LOG_SAMPLING_THEREAFTER = '100'

# 运行程序（程序会检测到没有 config.yaml 并使用环境变量）
go run .
```

环境变量命名规则映射示例：

- `app.name` -> `MYAPP_APP_NAME`
- `log.encoder.time_format` -> `MYAPP_LOG_ENCODER_TIME_FORMAT`

如果你使用其他 shell（例如 bash），请按 shell 语法设置环境变量，例如：

```bash
export MYAPP_LOG_LEVEL=debug
export MYAPP_LOG_ENCODER_TIME_FORMAT='2006-01-02T15:04:05.000Z07:00'
./your-program
```

## 命令行参数（flags）

`Viper2/main.go` 支持两个命令行 flag，用于灵活选择配置来源：

- `-config`（string）
  - 指定配置文件路径（默认 `./config.yaml`）。
  - 如果设置为空字符串 `""`，程序将不会尝试读取文件，只使用环境变量。
  - 行为：如果指定文件存在，会读取该文件并与环境变量合并；如果指定文件不存在，会回退到仅使用环境变量并打印提示。

- `-env-prefix`（string）
  - 指定环境变量前缀（默认 `MYAPP`）。程序会把配置 key 中的 `.` 替换为 `_`，并在前面加上这个前缀。例如 `app.name` 映射为 `MYAPP_APP_NAME`。

示例：在没有 `config.yaml` 的情况下，使用自定义前缀并运行：

```powershell
$env:MYAPP2_APP_NAME = 'env-app'
$env:MYAPP2_LOG_LEVEL = 'info'
go run . -config "" -env-prefix MYAPP2
```

或者指定一个配置文件（优先使用该文件）：

```powershell
go run . -config ./custom-config.yaml -env-prefix MYAPP
```

注意：命令行 flag 优先级高于程序内默认值，但配置文件本身会与环境变量合并（环境变量会覆盖文件内相同字段）。
