# Community 模块

`internal/modules/community` 为 `frontend/` Nuxt 视频社区提供公开只读数据接口。当前阶段覆盖首页、分类、视频列表、视频详情、弹幕、搜索、创作者资料和关注推荐预览。

## 职责

- 提供 `/api/v1/public/community/*` 公开读取 API，供 Nuxt 前端在关闭 mock 后直接接入。
- 维护社区分类、创作者、视频、视频源、标签和弹幕的持久化读取模型。
- 保持 DTO、筛选条件、错误和仓储 contract 在模块内部，不污染根 `types`。

## 非职责

- 不处理投稿审核、真实关注关系、评论写入、点赞收藏写入或创作者后台管理。
- 不恢复插件系统或 `/api/v1/plugins`。
- 不让 Nuxt 页面凭空展示后端不存在的生产写入能力；当前评论和关注仍是前端本地体验层。

## 分层

| 目录 | 职责 |
| --- | --- |
| `model` | 社区公开 DTO、持久化模型、分页结果和稳定状态值 |
| `service` | 公开读取用例、分类树构造、搜索聚合和错误归一 |
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
- `GET /search`
- `GET /users/:handle`
- `GET /feed/following`

## 数据

迁移 `20260626000100_create_community_tables.sql` 创建社区表，并写入一组与 Nuxt 现有体验相同的初始内容。后续若增加真实投稿、审核或关系能力，应追加迁移，不修改历史迁移。

## 验证

```powershell
go test ./internal/modules/community/... -count=1 -mod=readonly
go test ./internal/transport/http -count=1 -mod=readonly
go run ./cmd/console api openapi --output docs/api/openapi.yaml
```
