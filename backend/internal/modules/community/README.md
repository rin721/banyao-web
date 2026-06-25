# Community 模块

`internal/modules/community` 为 `frontend/` Nuxt 视频社区提供公开社区数据接口。当前阶段覆盖首页、分类、视频列表、视频详情、弹幕、视频评论、搜索、创作者资料、匿名创作者关注和关注动态。

## 职责

- 提供 `/api/v1/public/community/*` 公开 API，供 Nuxt 前端在关闭 mock 后直接接入。
- 维护社区分类、创作者、视频、视频源、标签、弹幕、公开评论和匿名关注关系的持久化模型。
- 支持视频评论列表读取和轻量公开发布，并同步视频 `comment_count`。
- 支持以匿名 `clientId` 关注 / 取消关注创作者，并同步创作者 `follower_count`。
- 保持 DTO、筛选条件、错误和仓储 contract 在模块内部，不污染根 `types`。

## 非职责

- 不处理投稿审核、登录态用户关系归并、评论审核、评论编辑删除、点赞收藏写入、通知投递或创作者后台管理。
- 不恢复插件系统或 `/api/v1/plugins`。
- 不让 Nuxt 页面凭空展示后端不存在的生产写入能力；评论写入失败必须显式反馈，浏览器本地状态只能保存显示名称、匿名 clientId、关注缓存等前端体验辅助数据。

## 分层

| 目录 | 职责 |
| --- | --- |
| `model` | 社区公开 DTO、持久化模型、分页结果和稳定状态值 |
| `service` | 公开社区用例、分类树构造、评论校验、搜索聚合和错误归一 |
| `repository` | 使用数据库端口读取社区表，隔离 SQL/ORM 细节 |
| `handler` | HTTP 输入输出适配，统一返回 `types/result` 响应 |

## API

公开接口统一位于 `/api/v1/public/community`：

- `GET /status`
- `GET /home`
- `GET /categories`
- `GET /videos`
- `GET /videos/:idOrSlug`
- `GET /videos/:idOrSlug/danmaku`
- `GET /videos/:idOrSlug/comments`
- `POST /videos/:idOrSlug/comments`
- `GET /search`
- `GET /users/:handle`
- `GET /users/:handle/follow-state`
- `POST /users/:handle/follow`
- `DELETE /users/:handle/follow`
- `GET /feed/following`

## 数据

迁移 `20260626000100_create_community_tables.sql` 创建社区读取表，并写入一组与 Nuxt 现有体验相同的初始内容。迁移 `20260626000200_create_community_video_comments.sql` 追加公开评论表、初始评论和 `comment_count` 收敛。迁移 `20260626000300_create_community_creator_follows.sql` 追加匿名创作者关注关系表。后续若增加真实投稿、审核、登录态归并或通知能力，应追加迁移，不修改历史迁移。

## 验证

```powershell
go test ./internal/modules/community/... -count=1 -mod=readonly
go test ./internal/transport/http -count=1 -mod=readonly
go run ./cmd/console api openapi --output docs/api/openapi.yaml
```
