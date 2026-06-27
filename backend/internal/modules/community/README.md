# Community 模块

`internal/modules/community` 为 `frontend/` Nuxt 视频社区提供公开社区数据接口。当前阶段覆盖社区账号注册入口、首页、分类、视频列表、视频详情、弹幕读取与轻量发布、视频评论读取 / 轻量发布 / 本人编辑删除、视频级举报提交、投稿元数据待审核池、主系统审核队列、发布视频 ID 回写、审核发布时关联 system media 资产并生成社区视频记录、搜索、创作者资料、创作者关注、视频点赞 / 收藏 / 稍后看、观看历史、社区动态流查询 / 发布 / 本人编辑删除、关注动态和通知收件箱。首页公告来自公告模块的真实已发布公告；没有已发布公告时返回 `announcement=null`。

## 职责

- 提供 `/api/v1/public/community/*` 公开 API，供 Nuxt 前端在关闭 mock 后直接接入；`/status` 的 endpoint 清单由真实 route contract 注册结果派生。
- 提供社区账号注册入口，浏览器提交用户名、显示名、邮箱和密码，响应只暴露社区会话和账号展示需要的最小字段。
- 维护社区分类、创作者、视频、视频源、标签、弹幕、公开评论、社区动态、关注关系、视频互动关系、观看历史和通知收件箱的持久化模型。
- 支持视频弹幕列表读取和轻量公开发布。
- 支持视频评论列表按 `sort=newest/oldest` 读取、轻量公开发布和本人编辑 / 删除；评论归属来自匿名或账号 `clientId`，列表可返回 `ownedByCurrentClient`，删除后同步视频 `comment_count`。
- 支持视频级轻量举报提交，保存匿名 `clientId`、原因、补充说明和待处理状态。
- 支持投稿元数据提交和查询，保存社区账号或匿名 `clientId`、作者名、标题、简介、分类、标签、可见性和文件元数据，不保存文件字节。
- 支持主系统权限保护的投稿审核队列和审核状态流转：`community_submission:review` 可读取全部投稿、写入 `approved` / `rejected` / `published` 状态、审核备注、审核时间、受控 `mediaAssetId` 和发布视频 ID；`published` 可绑定既有视频 ID，也可在请求携带 `mediaAssetId` 与 `durationSeconds` 时从 `system_media_assets` 生成 `community_videos`、默认播放源、分类关联和标签记录；`sourceUrl` / `durationSeconds` 仍作为过渡发布路径保留。
- 支持社区账号或匿名 `clientId` 关注 / 取消关注创作者，并同步创作者 `follower_count`。
- 支持社区账号或匿名 `clientId` 点赞、收藏、取消互动、查询互动状态和读取收藏 / 稍后看列表。
- 支持社区账号或匿名 `clientId` 记录、查询和清空视频观看历史，保存最近观看时间与播放进度。
- 支持社区账号或匿名 `clientId` 查询通知、标记已读，并把评论、弹幕、举报、关注和视频互动写入轻量消息流。
- 支持查询公开社区动态时间线，并以社区账号或匿名 `clientId` 发布、本人编辑和本人删除轻量动态；动态可绑定视频并装饰创作者与视频摘要，前端只消费 `ownedByCurrentClient` 判断操作入口。
- 保持 DTO、筛选条件、错误和仓储 contract 在模块内部，不污染根 `types`。

社区账号会话面向普通观看与创作者流程，响应字段为 `userId`、`sessionId`、过期时间和 `account.id / account.handle / account.displayName`，前台页面只基于这些社区字段展示账号状态。账号动态和账号投稿使用 `account.displayName` 作为作者展示名，数据范围仍由登录态派生的账号 `clientId` 控制。

## 非职责

- 社区用户侧真实文件上传、视频转码、评论 / 弹幕审核、举报审核处理、外部通知投递、后台可视化审核页和创作者后台管理属于独立能力范围；当前普通投稿创建仍只保存浏览器草稿文件元数据，真实文件字节通过主系统 system media 上传入口进入平台，并由审核发布请求用 `mediaAssetId` 受控关联。
- 生产交付面聚焦 `/api/v1/public/community/*` 与模块化扩展链路。
- Nuxt 页面展示的生产写入能力以当前社区公开 API 为准；评论写入失败必须显式反馈，弹幕写入失败只能降级到浏览器缓存，浏览器本地状态只保存显示名称、匿名 clientId、观看历史和必要降级缓存。

## 分层

| 目录 | 职责 |
| --- | --- |
| `model` | 社区公开 DTO、持久化模型、分页结果和稳定状态值 |
| `service` | 公开社区用例、分类树构造、评论校验、搜索聚合和错误归一 |
| `repository` | 使用数据库端口读取社区表，隔离 SQL/ORM 细节 |
| `handler` | HTTP 输入输出适配，统一返回 `types/result` 响应 |

## API

公开接口统一位于 `/api/v1/public/community`：

`GET /status` 始终可读，响应包含 `setup.required`、`setup.completed`、`setup.currentStep` 和由真实路由契约注册结果派生的 endpoint 清单。社区公开内容、账号接口和社区认证接口在平台初始化未完成时返回 503 result envelope，`messageKey` 为 `api.setup.required`，前端据此显示初始化状态并保留真实数据入口。

- `GET /status`
- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/session`
- `POST /auth/signup`
- `POST /account/dynamics`
- `PATCH /account/dynamics/:dynamicId`
- `DELETE /account/dynamics/:dynamicId`
- `GET /account/feed/following`
- `GET /account/history`
- `POST /account/history/clear`
- `GET /account/library`
- `GET /account/notifications`
- `POST /account/notifications/read`
- `GET /account/submissions`
- `POST /account/submissions`
- `GET /account/users/:handle/follow-state`
- `POST /account/users/:handle/follow`
- `DELETE /account/users/:handle/follow`
- `GET /account/videos/:idOrSlug/interaction-state`
- `POST /account/videos/:idOrSlug/interactions/:kind`
- `DELETE /account/videos/:idOrSlug/interactions/:kind`
- `POST /account/videos/:idOrSlug/history`
- `GET /home`
- `GET /categories`
- `GET /videos`
- `GET /videos/:idOrSlug`
- `GET /videos/:idOrSlug/danmaku`
- `POST /videos/:idOrSlug/danmaku`
- `GET /videos/:idOrSlug/comments`
- `POST /videos/:idOrSlug/comments`
- `PATCH /videos/:idOrSlug/comments/:commentId`
- `DELETE /videos/:idOrSlug/comments/:commentId`
- `PATCH /account/videos/:idOrSlug/comments/:commentId`
- `DELETE /account/videos/:idOrSlug/comments/:commentId`
- `GET /videos/:idOrSlug/interaction-state`
- `POST /videos/:idOrSlug/interactions/:kind`
- `DELETE /videos/:idOrSlug/interactions/:kind`
- `POST /videos/:idOrSlug/reports`
- `GET /submissions`
- `POST /submissions`
- `GET /api/v1/community/submissions`（主系统权限 `community_submission:review`）
- `PATCH /api/v1/community/submissions/:submissionId/review`（主系统权限 `community_submission:review`）
- `GET /notifications`
- `POST /notifications/read`
- `GET /dynamics`
- `POST /dynamics`
- `PATCH /dynamics/:dynamicId`
- `DELETE /dynamics/:dynamicId`
- `GET /search`
- `GET /users/:handle`
- `GET /users/:handle/follow-state`
- `POST /users/:handle/follow`
- `DELETE /users/:handle/follow`
- `GET /feed/following`
- `GET /library`
- `GET /history`
- `POST /history/clear`
- `POST /videos/:idOrSlug/history`

## 数据

社区数据表覆盖分类、创作者、视频、视频源、标签、弹幕、公开评论、关注关系、视频互动、举报、通知、动态、投稿元数据和观看历史。平台初始化完成后，公开接口读取这些表并返回真实社区内容；平台初始化未完成时，`/status` 只返回接口清单和 setup 状态，内容接口保持统一的初始化响应。真实运行路径只保留中性的社区分类 taxonomy；早期迁移写入的演示视频、创作者、动态、评论、弹幕、播放源、标签和相关派生记录由 `internal/migrations/20260627000200_remove_community_demo_content.sql` 按固定 demo ID 清理，Nuxt mock 继续保留演示内容。社区账号与匿名客户端通过 `client_id` 统一区分数据范围；视频评论和社区动态的 `client_id` 用于本人编辑 / 删除。投稿元数据额外保存 `review_note`、`reviewer_id`、`reviewed_at`、`media_asset_id`、`published_video_id` 和 `published_at`，用于主系统审核队列、system media 资产关联和发布状态回写；审核发布生成视频时会写入 `community_videos`、`community_video_sources`、`community_video_categories`、`community_video_tags` 和基于投稿作者名生成的社区创作者资料，不写入演示型 bio 或默认展示名。

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
