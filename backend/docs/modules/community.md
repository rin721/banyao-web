# Community 模块

`internal/modules/community` 为 `frontend/` Nuxt 视频社区提供公开社区数据接口。当前模块覆盖首页聚合、分类树、视频列表、视频详情、弹幕、视频评论、搜索、创作者资料、匿名创作者关注和关注动态。

## 当前能力

| 能力 | 路由 | 说明 |
| --- | --- | --- |
| API 状态 | `GET /api/v1/public/community/status` | 返回数据源模式、base path、端点清单和更新时间 |
| 首页聚合 | `GET /api/v1/public/community/home` | 返回公告、分类树和最新视频 |
| 分类树 | `GET /api/v1/public/community/categories` | 返回顶层分类及子分类 |
| 视频列表 | `GET /api/v1/public/community/videos` | 支持 `category`、`q`、`limit`、`cursor` 查询 |
| 视频详情 | `GET /api/v1/public/community/videos/:idOrSlug` | 返回播放源、标签和相关推荐 |
| 弹幕 | `GET /api/v1/public/community/videos/:idOrSlug/danmaku` | 返回视频初始弹幕列表 |
| 评论列表 | `GET /api/v1/public/community/videos/:idOrSlug/comments` | 支持 `sort`、`limit` 查询，返回公开可见评论 |
| 发布评论 | `POST /api/v1/public/community/videos/:idOrSlug/comments` | 公开发布轻量评论，并同步视频评论数 |
| 搜索 | `GET /api/v1/public/community/search` | 聚合视频、创作者和分类搜索结果 |
| 创作者资料 | `GET /api/v1/public/community/users/:handle` | 返回创作者资料、分类和最新视频 |
| 关注状态 | `GET /api/v1/public/community/users/:handle/follow-state` | 按匿名 `clientId` 返回当前创作者关注状态和关注者数 |
| 关注 / 取消关注 | `POST/DELETE /api/v1/public/community/users/:handle/follow` | 以匿名 `clientId` 写入或软删除创作者关注关系，并同步 `follower_count` |
| 关注动态 | `GET /api/v1/public/community/feed/following` | 传入 `clientId` 时返回已关注创作者和更新；无关注时返回推荐预览 |

## 边界

- 本模块的社区接口是公开接口，不提供 IAM 权限码，也不写入 `system_apis` 的受保护权限目录。
- 当前评论发布与创作者关注不绑定登录态或审核流；评论编辑删除、举报处理、投稿审核、登录态关系归并、点赞收藏写入、通知和创作者后台管理仍未实现。
- 前端的历史、收藏和投稿草稿仍是浏览器本地体验层；评论和关注只在后端 API 不可用时降级到本地浏览器状态。

## 数据与装配

- 迁移 `internal/migrations/20260626000100_create_community_tables.sql` 创建社区分类、创作者、视频、视频分类、标签、播放源和弹幕表，并写入初始内容。
- 迁移 `internal/migrations/20260626000200_create_community_video_comments.sql` 创建社区视频评论表、写入初始评论，并将视频 `comment_count` 收敛为评论表统计值。
- 迁移 `internal/migrations/20260626000300_create_community_creator_follows.sql` 创建匿名创作者关注关系表。
- 应用装配位于 `internal/app/initapp`，HTTP contract 位于 `internal/transport/http/contracts.go`，真实路由注册位于 `internal/transport/http/router.go`。
- OpenAPI 由 route contract 生成到 `docs/api/openapi.yaml`，不得手写维护。

## 验证

```powershell
go test ./internal/modules/community/... -count=1 -mod=readonly
go test ./internal/transport/http -count=1 -mod=readonly
go run ./cmd/console api openapi --output docs/api/openapi.yaml
```
