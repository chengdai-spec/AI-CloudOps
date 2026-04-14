# 关键设计决策说明

## 1) 选择 Google Wire（编译期 DI）

- 采用 Wire 作为依赖注入方案，避免运行时容器引入隐藏依赖与全局状态。
- 优点：依赖图可编译期检查，启动路径更清晰，便于测试替换。

## 2) `internal/app` 作为生命周期编排入口

- 将 HTTP Server、Asynq Server/Scheduler、Cron Manager 等启动/退出逻辑集中在一个位置，减少散落的 goroutine 与全局关闭顺序问题。

## 3) Pod 的传输能力回归 Handler

问题：日志 SSE、终端 WebSocket、文件下载响应头与数据写入属于传输层行为，如果放在 Service 会导致 Gin/HTTP 细节向下泄漏。

解决：
- Service：只返回 `io.ReadCloser` / 执行 `PodTerminalSession` / 触发上传下载。
- Handler：负责 SSE/WebSocket 升级、写响应头、按连接状态终止流。

## 4) WebSocket Upgrader 去全局化

- 取消 `ssh.UpGrader` 全局变量的使用，改为在 DI 层构造 `*websocket.Upgrader` 并注入到需要的 Handler。
- 目的：避免全局单例导致测试困难与隐式依赖。

## 5) Cron 载荷下沉到 `internal/cron/task`

- Service/Scheduler 只关心任务载荷结构，不应依赖 Handler 包（否则形成反向依赖，破坏分层）。

