# 深度重构 - TaskList

## Overview

本任务用于推进 AI-CloudOps-backend 的分层纯度、上下文透传与依赖注入规范化，提升可测试性与可维护性。

## Tasks

- [x] **收敛配置加载**
  Description: 将配置加载集中到 `internal/config`，隔离 viper 的散落使用
  Priority: High
  Category: Architecture
  Dependencies: None
  Estimated Effort: L

- [x] **重构启动生命周期**
  Description: 引入 `internal/app` 统一管理 HTTP/Asynq/Cron 的启动与优雅退出
  Priority: High
  Category: Backend
  Dependencies: 收敛配置加载
  Estimated Effort: L

- [x] **Handler 透传 request ctx**
  Description: 将 `ctx.Request.Context()` 作为 Service 调用的上下文入参，杜绝 `*gin.Context` 下传
  Priority: High
  Category: Backend
  Dependencies: None
  Estimated Effort: M

- [x] **Pod SSE/WebSocket/File 解耦**
  Description: 将日志/终端/文件传输的 HTTP 细节从 Service 层移回 Handler，Service 仅保留业务与调用
  Priority: High
  Category: Backend
  Dependencies: Handler 透传 request ctx
  Estimated Effort: L

- [x] **Cron 载荷解耦**
  Description: 抽离 `CronTaskPayload` 到 `internal/cron/task`，消除 Service/Scheduler 对 Handler 的依赖
  Priority: High
  Category: Architecture
  Dependencies: None
  Estimated Effort: S

- [ ] **清理 viper.Get* 散落使用**
  Description: 将散落的 viper 调用逐步迁移到 `internal/config`，并通过 `config.Config` 注入
  Priority: Medium
  Category: Architecture
  Dependencies: 收敛配置加载
  Estimated Effort: XL

- [ ] **统一数据库事务闭包模式**
  Description: 引入 tx helper，减少显式 Begin/Commit/Rollback 与泄漏风险
  Priority: Medium
  Category: Backend
  Dependencies: None
  Estimated Effort: L

- [ ] **补齐关键接口测试**
  Description: 为 SSE/WebSocket/文件流接口补充 `httptest` 与取消场景测试
  Priority: Medium
  Category: Testing
  Dependencies: Pod SSE/WebSocket/File 解耦
  Estimated Effort: M

## Progress Tracking

- Total Tasks: 8
- Completed: 5
- In Progress: 0
- Remaining: 3

## Next Steps

1. 优先清理最核心链路的 `viper.Get*`（cron/tree/prometheus）
2. 引入事务闭包 helper 并逐模块迁移
3. 为 Pod 日志与文件流接口补齐测试，覆盖取消与错误分支

## Notes

- 后续如需要更强的约束，可以引入 lint 规则禁止在 `internal/**/service` 与 `internal/**/dao` 引用 `gin` 包。

