# AI-CloudOps-backend 深度重构计划（Gin + Gorm + DI）

## 目标

1. 严格分层：Handler → Service → Repository（DAO），禁止跨层泄漏。
2. 统一上下文：所有 Service/Repository 方法 `ctx context.Context` 必须作为第一个参数，且从 HTTP 请求透传。
3. 依赖注入：依赖全部通过构造函数注入（本项目采用 Google Wire）。
4. 可观测性：日志消息使用中文、结构化字段清晰，错误返回可定位、可追踪。

## 范围

- 应用启动链路与生命周期管理
- 配置加载与校验
- 关键模块 Handler → Service 传参与解耦（重点：K8s Pod、Cron、Tree 终端）

## 已完成

### 1) 启动链路与配置收敛

- 将配置加载收敛到 `internal/config`，隔离 viper 的使用范围。
- 引入 `internal/app` 统一管理 HTTP/异步任务/定时任务的生命周期与优雅退出。
- Wire 注入入口调整为 `ProvideApp(*config.Config)`，移除旧的全局配置注入方式。

### 2) Handler 透传 request ctx

- 统一将 `ctx.Request.Context()` 作为服务层调用的上下文入参，避免把 `*gin.Context` 传入 Service/Repository。

### 3) Pod 相关传输层泄漏治理

- 将 K8s Pod 的日志/终端/文件上传下载从 Service 层剥离传输细节：
  - Service 仅负责参数校验、调用 Manager 获取流或执行会话。
  - Handler 负责 SSE/WebSocket/文件响应头与 `io.Copy`。

### 4) Cron 反向依赖解除

- 抽离 Cron 任务载荷结构体到 `internal/cron/task`，避免 Service/Scheduler 依赖 Handler 包。

## 未完成（建议后续继续推进）

1. 逐步清理 `viper.Get*` 的散落使用，继续向 `internal/config` 迁移。
2. 统一数据库事务处理为闭包模式（`WithTx(ctx, func(tx *gorm.DB) error { ... })`），减少显式 `Begin/Commit/Rollback` 泄漏。
3. 补齐关键 Handler 的 `httptest` 测试（特别是 SSE/WebSocket/文件流接口的错误分支与取消场景）。
4. 统一 zap 字段 key（建议中文 key 或约定清单），避免同类字段在不同模块命名不一致。

## 验证

- 全量执行：`go test ./...`
- 若引入 Wire 变更：执行 `wire ./pkg/di` 后再测试

