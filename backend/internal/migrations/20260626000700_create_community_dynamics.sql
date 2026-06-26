-- +goose Up
CREATE TABLE IF NOT EXISTS community_dynamics (
    id VARCHAR(96) PRIMARY KEY,
    client_id VARCHAR(96) NOT NULL DEFAULT '',
    creator_id VARCHAR(96) NOT NULL DEFAULT '',
    video_id VARCHAR(96) NOT NULL DEFAULT '',
    author_name VARCHAR(120) NOT NULL,
    body VARCHAR(500) NOT NULL,
    kind VARCHAR(32) NOT NULL DEFAULT 'text',
    status VARCHAR(32) NOT NULL DEFAULT 'visible',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_community_dynamics_status_created
    ON community_dynamics (status, created_at);
CREATE INDEX idx_community_dynamics_creator_created
    ON community_dynamics (creator_id, created_at);
CREATE INDEX idx_community_dynamics_client_created
    ON community_dynamics (client_id, created_at);
CREATE INDEX idx_community_dynamics_video_created
    ON community_dynamics (video_id, created_at);

INSERT INTO community_dynamics (id, client_id, creator_id, video_id, author_name, body, kind, status, created_at, updated_at) VALUES
('dynamic-rin-alpha', '', 'user-rin', 'video-aoi-alpha', 'Rin721', '今天把首页、评论、弹幕和通知都切到后端社区模块了。接下来会继续把动态、投稿和审核能力从界面草图推进到真实数据。', 'video_update', 'visible', '2026-06-26 10:00:00', '2026-06-26 10:00:00'),
('dynamic-design-sakura', '', 'user-design', 'video-sakura-accent', 'Design Note', '这版视觉会保留 kirakira 式的细线、粉色强调和轻量动效，但内容密度会更适合长期视频社区，不会变成一次性的营销页。', 'video_update', 'visible', '2026-06-26 09:20:00', '2026-06-26 09:20:00'),
('dynamic-backend-contract', '', 'user-backend', 'video-go-api', 'Aoi Backend', '社区公开接口继续以 route contract 生成 OpenAPI，前端页面只通过 API client 消费数据，避免 mock 与真实接口双轨漂移。', 'video_update', 'visible', '2026-06-26 08:40:00', '2026-06-26 08:40:00'),
('dynamic-frontend-mobile', '', 'user-frontend', 'video-mobile-grid', 'Frontend Memo', '移动端这一轮重点看文字换行、触控间距和底部导航遮挡，动态卡片会先保持短句节奏和稳定媒体比例。', 'video_update', 'visible', '2026-06-26 08:10:00', '2026-06-26 08:10:00');

-- +goose Down
DROP TABLE IF EXISTS community_dynamics;
