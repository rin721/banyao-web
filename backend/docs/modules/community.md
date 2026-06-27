# Community 模块

`backend/internal/modules/community` 为 `frontend/` Nuxt 视频社区提供公开社区数据接口。当前模块覆盖社区账号注册与会话、首页聚合、分类树、视频列表与详情、弹幕、评论发布与本人编辑删除、搜索、创作者资料、关注动态、动态发布与本人编辑删除、视频互动、收藏 / 稍后看、观看历史、投稿元数据、受控审核队列、发布视频 ID 回写、审核发布时关联 system media 资产并生成社区视频记录、举报和通知收件箱。

## 当前能力

| 能力 | 路由 | 说明 |
| --- | --- | --- |
| API 状态 | `GET /api/v1/public/community/status` | 始终可读，返回数据源模式、base path、端点清单、响应耗时、更新时间和 `setup.required/completed/currentStep` |
| 社区账号 | `POST /auth/signup`、`POST /auth/login`、`GET /auth/session`、`POST /auth/logout` | 使用普通社区账号语义，只向前端暴露 `userId`、`sessionId`、过期时间和 `account.id/handle/displayName` |
| 首页聚合 | `GET /home` | 返回可选公告、分类树、最新视频和社区动态；无真实公告来源时 `announcement` 为 `null` |
| 分类与视频 | `GET /categories`、`GET /videos`、`GET /videos/:idOrSlug` | 支持分类、关键词、游标和数量限制，视频详情包含播放源、标签和相关推荐 |
| 弹幕 | `GET /videos/:idOrSlug/danmaku`、`POST /videos/:idOrSlug/danmaku` | 读取初始弹幕并支持轻量发布；发布失败时前端可降级到浏览器体验层 |
| 评论 | `GET /videos/:idOrSlug/comments`、`POST /videos/:idOrSlug/comments`、`PATCH/DELETE /videos/:idOrSlug/comments/:commentId`、`PATCH/DELETE /account/videos/:idOrSlug/comments/:commentId` | 支持 `sort=newest/oldest` 与数量限制；发布、本人编辑和本人删除按匿名或账号 `clientId` 归属校验，删除后同步视频评论数；列表可通过 `clientId` 返回 `ownedByCurrentClient` |
| 视频互动 | `GET /videos/:idOrSlug/interaction-state`、`POST/DELETE /videos/:idOrSlug/interactions/:kind` | 匿名 `clientId` 或社区账号范围写入点赞、收藏、稍后看状态 |
| 观看历史 | `GET /history`、`POST /history/clear`、`POST /videos/:idOrSlug/history` | 匿名或社区账号范围记录最近观看时间和播放进度 |
| 资料库 | `GET /library`、`GET /account/library` | 读取收藏和稍后看列表 |
| 投稿元数据 | `GET/POST /submissions`、`GET/POST /account/submissions` | 保存作者、标题、简介、分类、标签、可见性和文件元数据，不保存文件字节 |
| 投稿审核 | `GET /api/v1/community/submissions`、`PATCH /api/v1/community/submissions/:submissionId/review` | 主系统 IAM 权限 `community_submission:review` 保护；支持 `approved`、`rejected`、`published` 状态，写入审核备注、审核人、审核时间、受控 `mediaAssetId` 和发布视频 ID；`published` 可绑定既有视频 ID，也可在请求携带 `mediaAssetId` 与 `durationSeconds` 时从 `system_media_assets` 生成 `community_videos`、默认播放源、分类关联和标签记录；`sourceUrl` / `durationSeconds` 仍保留为过渡发布路径 |
| 举报 | `POST /videos/:idOrSlug/reports` | 保存匿名 `clientId`、原因、补充说明和待处理状态 |
| 通知 | `GET /notifications`、`POST /notifications/read`、`GET /account/notifications`、`POST /account/notifications/read` | 评论、弹幕、举报、关注和视频互动写入轻量消息流，支持标记已读 |
| 动态流 | `GET /dynamics`、`POST /dynamics`、`PATCH/DELETE /dynamics/:dynamicId`、`POST /account/dynamics`、`PATCH/DELETE /account/dynamics/:dynamicId` | 查询公开动态时间线，匿名或账号发布、本人编辑和本人删除轻量动态，可绑定视频摘要；列表通过 `clientId` 返回 `ownedByCurrentClient` |
| 创作者关注 | `GET /users/:handle`、`GET /users/:handle/follow-state`、`POST/DELETE /users/:handle/follow`、`GET /feed/following` | 匿名或账号关注创作者，并返回关注创作者的视频与动态 |
| 账号范围接口 | `/account/**` | 需要社区账号会话与 CSRF token，数据范围由账号派生的 `clientId` 控制 |

## Setup Gate

- `GET /api/v1/public/community/status` 不受 setup gate 阻断，用于前端判断真实 API 是否可用。
- 平台首次初始化未完成时，社区内容接口、社区账号认证接口和 `/account/**` 接口统一返回 HTTP `503` 的 `types/result` envelope。
- 初始化阻断响应使用稳定 `messageKey=api.setup.required`，`data` 为 `CommunitySetupStatus`，包含 `required`、`completed` 和 `currentStep`。
- 前端必须把该状态展示为“真实后端尚未初始化”的引导态，不得回退成伪造真实数据；Nuxt mock 可以继续提供演示数据，但必须显示为 mock 边界。

## 边界

- 社区公开接口不提供 IAM 权限码，也不写入 `system_apis` 的受保护权限目录；它们仍通过 route contract 生成 OpenAPI。
- 社区账号用于普通观看、创作者互动和投稿流程，不暴露后台组织、角色、权限或控制台身份字段。
- 当前评论、弹幕、动态、投稿和举报是轻量生产接口；评论和动态支持本人编辑 / 删除，投稿支持主系统权限保护的审核状态流转、发布视频 ID 回写，并能在审核请求携带受控 `mediaAssetId` 时从 system media 资产生成社区视频记录。普通社区投稿创建仍只保存文件元数据，真实文件字节通过主系统 system media 上传入口进入平台；转码、公开播放源生成、后台可视化审核页、举报处理、登录态与匿名关系归并、外部通知投递和创作者后台管理仍属于后续任务。
- 浏览器本地状态只保存匿名 `clientId`、显示偏好、必要降级缓存和上传草稿文件元数据；不得保存文件字节、后台权限 payload 或不可恢复的大对象。

## 数据与装配

- 迁移 `internal/migrations/20260626000100_create_community_tables.sql` 创建社区分类、创作者、视频、视频分类、标签、播放源和弹幕表；真实运行路径保留中性分类 taxonomy，早期演示创作者、视频、播放源、标签和弹幕由后续清理迁移删除。
- 迁移 `internal/migrations/20260626000200_create_community_video_comments.sql` 创建社区视频评论表；早期演示评论由后续清理迁移删除，真实评论来自公开或账号评论接口写入。
- 迁移 `internal/migrations/20260626000300_create_community_creator_follows.sql` 创建匿名创作者关注关系表。
- 迁移 `internal/migrations/20260626000400_create_community_video_interactions.sql` 创建匿名视频点赞、收藏和稍后看关系表。
- 迁移 `internal/migrations/20260626000500_create_community_reports.sql` 创建社区举报记录表。
- 迁移 `internal/migrations/20260626000600_create_community_notifications.sql` 创建匿名通知表。
- 迁移 `internal/migrations/20260626000700_create_community_dynamics.sql` 创建社区动态表；动态表已有 `client_id` 归属列，用于匿名或账号本人编辑 / 删除，早期演示时间线由后续清理迁移删除。
- 迁移 `internal/migrations/20260626000800_create_community_submissions.sql` 创建社区投稿元数据表。
- 迁移 `internal/migrations/20260626000900_create_community_video_history.sql` 创建观看历史表。
- 迁移 `internal/migrations/20260626001000_refine_community_demo_copy.sql` 收敛演示内容文案。
- 迁移 `internal/migrations/20260626001100_add_community_video_comment_client_id.sql` 为视频评论补充 `client_id` 归属列和查询索引，用于本人评论编辑 / 删除；历史种子评论保持空 `client_id`，因此默认只读。
- 迁移 `internal/migrations/20260627000100_add_community_submission_review_state.sql` 为投稿元数据补充审核备注、审核人、审核时间、发布视频 ID 和发布时间，用于主系统审核队列与发布状态回写；审核发布生成视频复用既有 `community_videos`、`community_video_sources`、`community_video_categories` 和 `community_video_tags` 表，不为发布链路新增专门表结构。
- 迁移 `internal/migrations/20260627000200_remove_community_demo_content.sql` 只按固定 demo ID 清理早期演示视频、创作者、动态、评论、弹幕、播放源、标签、视频分类和相关派生记录；不删除社区基础分类 taxonomy。
- 迁移 `internal/migrations/20260627000300_add_community_submission_media_asset.sql` 为投稿元数据补充 `media_asset_id`，用于记录审核发布时关联的 `system_media_assets` 资产；社区 service 只读取最小媒体资产投影，不复制 system media 上传逻辑。
- 应用装配位于 `internal/app/initapp`，setup 状态由 `internal/app/initcenter` 适配为社区 `SetupStatus`，HTTP contract 位于 `internal/transport/http/contracts.go`，真实路由注册位于 `internal/transport/http/router.go`。
- OpenAPI 由 route contract 生成到 `docs/api/openapi.yaml`，不得手写维护。

## 前后端协作

- 后端新增或修改社区 API 时，先更新 route contract、稳定 DTO、handler/service/repository，再生成 OpenAPI。
- 前端只通过 `frontend/app/composables/useAoiApi.ts`、`frontend/app/composables/useAoiAuthApi.ts` 和 `frontend/shared/types/api.ts` 消费社区契约。
- Mock API 位于 `frontend/server/api/mock`，用于演示和调试；真实联调以 `GET /status` 的 `mode=go`、setup 状态和端点清单为准，真实 API 模式不得回退展示 mock/demo fixture。

## 验证

```powershell
go test ./internal/modules/community/... -count=1 -mod=readonly
go test ./internal/transport/http -count=1 -mod=readonly
go run ./cmd/console api openapi --output docs/api/openapi.yaml
```

聚合仓库根目录还提供社区前后端联调烟测：

```powershell
powershell -ExecutionPolicy Bypass -File scripts/check-frontend-community-api-smoke.ps1
```
