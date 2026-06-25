-- +goose Up
CREATE TABLE IF NOT EXISTS community_categories (
    id VARCHAR(96) PRIMARY KEY,
    slug VARCHAR(96) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    description VARCHAR(320) NULL,
    accent_color VARCHAR(32) NULL,
    parent_slug VARCHAR(96) NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_community_categories_parent_slug ON community_categories (parent_slug);

CREATE TABLE IF NOT EXISTS community_creators (
    id VARCHAR(96) PRIMARY KEY,
    handle VARCHAR(96) NOT NULL UNIQUE,
    display_name VARCHAR(120) NOT NULL,
    avatar_url VARCHAR(512) NULL,
    bio VARCHAR(640) NULL,
    follower_count BIGINT NOT NULL DEFAULT 0,
    joined_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS community_videos (
    id VARCHAR(96) PRIMARY KEY,
    slug VARCHAR(160) NOT NULL UNIQUE,
    title VARCHAR(240) NOT NULL,
    description VARCHAR(720) NULL,
    thumbnail_url VARCHAR(512) NOT NULL,
    duration_seconds INTEGER NOT NULL,
    view_count BIGINT NOT NULL DEFAULT 0,
    comment_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    source_url VARCHAR(512) NOT NULL,
    published_at TIMESTAMP NOT NULL,
    uploader_id VARCHAR(96) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_community_videos_published_at ON community_videos (published_at);
CREATE INDEX idx_community_videos_uploader_id ON community_videos (uploader_id);

CREATE TABLE IF NOT EXISTS community_video_categories (
    video_id VARCHAR(96) NOT NULL,
    category_slug VARCHAR(96) NOT NULL,
    PRIMARY KEY (video_id, category_slug)
);

CREATE TABLE IF NOT EXISTS community_video_tags (
    video_id VARCHAR(96) NOT NULL,
    tag VARCHAR(96) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (video_id, tag)
);

CREATE TABLE IF NOT EXISTS community_video_sources (
    id VARCHAR(96) PRIMARY KEY,
    video_id VARCHAR(96) NOT NULL,
    src VARCHAR(512) NOT NULL,
    kind VARCHAR(32) NOT NULL,
    label VARCHAR(120) NOT NULL,
    mime_type VARCHAR(120) NULL,
    quality_label VARCHAR(64) NULL,
    bitrate_kbps INTEGER NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_community_video_sources_video_id ON community_video_sources (video_id);

CREATE TABLE IF NOT EXISTS community_video_danmaku (
    id VARCHAR(96) PRIMARY KEY,
    video_id VARCHAR(96) NOT NULL,
    body VARCHAR(280) NOT NULL,
    time_seconds INTEGER NOT NULL,
    mode VARCHAR(24) NOT NULL,
    color VARCHAR(32) NOT NULL,
    author_name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_community_video_danmaku_video_id ON community_video_danmaku (video_id);

INSERT INTO community_categories (id, slug, name, description, accent_color, parent_slug, display_order, created_at, updated_at) VALUES
('cat-home', 'home', '首页', '全部精选内容', '#111111', NULL, 0, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-creative', 'creative', '创作', '动画、音乐、MAD 与视觉表达', '#f2709c', NULL, 10, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-animation', 'animation', '动画', '动画与短片', '#22b8cf', 'creative', 10, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-music', 'music', '音乐', '音乐视频与现场', '#5b8def', 'creative', 20, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-mad', 'mad', '音 MAD', '混剪与再创作', '#f2709c', 'creative', 30, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-design', 'design', '设计', '视觉与交互', '#111111', 'creative', 40, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-knowledge', 'knowledge', '知识', '技术、开发与硬件知识', '#0f9fb7', NULL, 20, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-tech', 'tech', '科技', '技术、开发与硬件', '#0f9fb7', 'knowledge', 10, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-play', 'play', '游玩', '游戏、实况与互动娱乐', '#7a68f0', NULL, 30, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-games', 'games', '游戏', '游戏与实况', '#7a68f0', 'play', 10, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-community', 'community', '社区', '社区综合内容与站内动态', '#64757b', NULL, 40, '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('cat-general', 'general', '综合', '社区综合内容', '#64757b', 'community', 10, '2026-06-26 00:00:00', '2026-06-26 00:00:00');

INSERT INTO community_creators (id, handle, display_name, avatar_url, bio, follower_count, joined_at, created_at, updated_at) VALUES
('user-rin', 'rin721', 'Rin721', NULL, 'Banyao 发起者，关注清晰的视频社区体验、可演进前端和可靠后端契约。', 3420, '2026-04-10 00:00:00', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('user-lab', 'aoi-lab', 'Aoi Lab', NULL, '原型实验室，持续发布视觉、交互和内容社区探索。', 2180, '2026-04-18 00:00:00', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('user-backend', 'aoi-backend', 'Aoi Backend', NULL, '记录 Go API、route contract、数据流和工程化取舍。', 1480, '2026-05-08 00:00:00', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('user-frontend', 'frontend-memo', 'Frontend Memo', NULL, '前端布局、响应式细节和交互体验的实现备忘录。', 1320, '2026-05-01 00:00:00', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('user-design', 'design-note', 'Design Note', NULL, '把复杂视觉系统拆成可复用、可维护的产品界面记录。', 960, '2026-04-25 00:00:00', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('user-motion', 'aoi-motion', 'Aoi Motion', NULL, '关注柔和动效、状态反馈和媒体界面的情绪表达。', 760, '2026-05-12 00:00:00', '2026-06-26 00:00:00', '2026-06-26 00:00:00');

INSERT INTO community_videos (id, slug, title, description, thumbnail_url, duration_seconds, view_count, comment_count, like_count, source_url, published_at, uploader_id, created_at, updated_at) VALUES
('video-aoi-alpha', 'aoi-alpha', 'Banyao Alpha：几何品牌层与社区首页的第一次试映', '从 Nuxt mock 迁移到 Go 社区模块，并吸收 kirakira 的黑白几何视觉。', 'gradient:aoi-alpha', 300, 1200, 36, 100, 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', '2026-06-03 10:00:00', 'user-rin', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-token-array', 'token-array', 'Colorful Array：青粉主题 token 的几何折线实验', 'Aoi token 与硬边线框背景的映射记录。', 'gradient:token-array', 198, 856, 18, 71, 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', '2026-06-02 10:00:00', 'user-lab', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-dark-mode', 'dark-mode', '夜间模式预览：暗色界面里的高可读内容网格', '暗色主题与媒体卡片可读性演示。', 'gradient:dark-mode', 252, 420, 9, 35, 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', '2026-06-01 10:00:00', 'user-design', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-mobile-grid', 'mobile-grid', '移动端双列卡片：避免裁切的响应式布局', '390px 移动视口下的网格约束与底部导航。', 'gradient:mobile-grid', 384, 638, 12, 53, 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', '2026-05-30 10:00:00', 'user-frontend', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-go-api', 'go-api-ready', 'Go API Ready：前端数据如何切换后端社区模块', 'DTO、runtime config 与公开读取接口的接入记录。', 'gradient:go-api', 496, 1800, 44, 150, 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', '2026-05-28 10:00:00', 'user-backend', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-sakura-accent', 'sakura-accent', 'Sakura Accent：柔和粉色在状态反馈中的使用方式', '状态、动效与品牌点缀。', 'gradient:sakura-accent', 165, 512, 15, 42, 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', '2026-05-26 10:00:00', 'user-motion', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-music-stream', 'music-stream', 'Banyao Session：轻音乐频道的首个社区播放清单', '用内容 feed 模拟音乐频道的沉浸式浏览体验。', 'gradient:music-stream', 276, 980, 22, 82, 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', '2026-05-24 10:00:00', 'user-lab', '2026-06-26 00:00:00', '2026-06-26 00:00:00'),
('video-game-room', 'game-room', 'Game Room：游戏分区的信息密度与卡片状态', '游戏内容如何在两列移动卡片里保持可读。', 'gradient:game-room', 612, 2310, 67, 192, 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', '2026-05-21 10:00:00', 'user-frontend', '2026-06-26 00:00:00', '2026-06-26 00:00:00');

INSERT INTO community_video_categories (video_id, category_slug) VALUES
('video-aoi-alpha', 'design'),
('video-aoi-alpha', 'tech'),
('video-token-array', 'design'),
('video-dark-mode', 'design'),
('video-mobile-grid', 'tech'),
('video-mobile-grid', 'design'),
('video-go-api', 'tech'),
('video-sakura-accent', 'animation'),
('video-sakura-accent', 'design'),
('video-music-stream', 'music'),
('video-music-stream', 'general'),
('video-game-room', 'games'),
('video-game-room', 'tech');

INSERT INTO community_video_tags (video_id, tag, display_order) VALUES
('video-aoi-alpha', 'Banyao', 10),
('video-aoi-alpha', '设计', 20),
('video-aoi-alpha', 'Go API', 30),
('video-token-array', '视觉 token', 10),
('video-token-array', '几何线框', 20),
('video-dark-mode', '暗色模式', 10),
('video-mobile-grid', '移动端', 10),
('video-mobile-grid', '响应式', 20),
('video-go-api', '后端接入', 10),
('video-go-api', 'route contract', 20),
('video-sakura-accent', '动效', 10),
('video-music-stream', '音乐', 10),
('video-game-room', '游戏', 10);

INSERT INTO community_video_sources (id, video_id, src, kind, label, mime_type, quality_label, bitrate_kbps, is_default, display_order) VALUES
('source-aoi-alpha-primary', 'video-aoi-alpha', 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', 'native', 'R2 示例 1', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-aoi-alpha-secondary', 'video-aoi-alpha', 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', 'native', 'R2 示例 2', 'video/mp4', 'Alt', NULL, FALSE, 20),
('source-token-array-primary', 'video-token-array', 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', 'native', 'R2 示例 1', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-dark-mode-primary', 'video-dark-mode', 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', 'native', 'R2 示例 1', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-mobile-grid-primary', 'video-mobile-grid', 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', 'native', 'R2 示例 1', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-go-api-primary', 'video-go-api', 'https://r2-store.kobayashi.eu.org/aoi/video/1e32a269-bde5-4eb6-9c7e-c35add52b482.mp4', 'native', 'R2 示例 1', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-sakura-primary', 'video-sakura-accent', 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', 'native', 'R2 示例 2', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-music-primary', 'video-music-stream', 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', 'native', 'R2 示例 2', 'video/mp4', 'Auto', NULL, TRUE, 10),
('source-game-primary', 'video-game-room', 'https://r2-store.kobayashi.eu.org/aoi/video/BV1EF3uzeETo.mp4', 'native', 'R2 示例 2', 'video/mp4', 'Auto', NULL, TRUE, 10);

INSERT INTO community_video_danmaku (id, video_id, body, time_seconds, mode, color, author_name, created_at) VALUES
('danmaku-aoi-alpha-1', 'video-aoi-alpha', 'Go 后端数据来了', 2, 'scroll', '#ffffff', 'Aoi Viewer', '2026-06-03 10:01:00'),
('danmaku-aoi-alpha-2', 'video-aoi-alpha', '几何线框很 KIRA', 7, 'scroll', '#7ee7ff', 'Design Note', '2026-06-03 10:02:00'),
('danmaku-aoi-alpha-3', 'video-aoi-alpha', '移动端也要稳', 13, 'top', '#ffe58a', 'Frontend Memo', '2026-06-03 10:03:00'),
('danmaku-aoi-alpha-4', 'video-aoi-alpha', 'route contract 对齐', 20, 'bottom', '#ffb4d8', 'Aoi Backend', '2026-06-03 10:04:00'),
('danmaku-go-api-1', 'video-go-api', '这个接口结构清楚', 4, 'scroll', '#ffffff', 'Aoi Viewer', '2026-05-28 10:01:00'),
('danmaku-go-api-2', 'video-go-api', '前端不用再猜数据了', 10, 'scroll', '#7ee7ff', 'Frontend Memo', '2026-05-28 10:02:00');

-- +goose Down
DROP TABLE IF EXISTS community_video_danmaku;
DROP TABLE IF EXISTS community_video_sources;
DROP TABLE IF EXISTS community_video_tags;
DROP TABLE IF EXISTS community_video_categories;
DROP TABLE IF EXISTS community_videos;
DROP TABLE IF EXISTS community_creators;
DROP TABLE IF EXISTS community_categories;
