-- +goose Up
UPDATE community_creators
SET handle = 'aoi-curator',
    display_name = 'Aoi Curator',
    bio = '整理 Aoi 的创作者投稿、频道主题和社区活动笔记。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'user-backend';

UPDATE community_creators
SET handle = 'layout-notes',
    display_name = 'Layout Notes',
    bio = '记录社区页面的阅读节奏、响应式细节和互动体验。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'user-frontend';

UPDATE community_creators
SET handle = 'color-note',
    display_name = 'Color Note',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'user-design';

UPDATE community_creators
SET bio = 'Aoi 发起者，关注清透、高信息可读性和可演进的社区产品体验。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'user-rin';

UPDATE community_videos
SET title = 'Aoi Alpha：清透社区首页的第一次视觉试映',
    description = '一次首页试映，观察清透标题、分类入口和动态流如何自然连在一起。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'video-aoi-alpha';

UPDATE community_videos
SET title = 'Colorful Array：青粉主题的几何折线实验',
    description = '用柔和色彩和几何折线做一组频道封面的延展实验。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'video-token-array';

UPDATE community_videos
SET title = 'Community Notes：从投稿到互动的顺畅动线',
    description = '围绕评论、收藏和关注整理一条轻快的社区浏览路径。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'video-go-api';

UPDATE community_video_tags
SET tag = '视觉试映'
WHERE video_id = 'video-aoi-alpha' AND display_order = 30;

UPDATE community_video_tags
SET tag = '视觉实验'
WHERE video_id = 'video-token-array' AND display_order = 10;

UPDATE community_video_tags
SET tag = '社区动线'
WHERE video_id = 'video-go-api' AND display_order = 10;

UPDATE community_video_tags
SET tag = '互动体验'
WHERE video_id = 'video-go-api' AND display_order = 20;

UPDATE community_video_danmaku
SET body = '开场好清爽'
WHERE id = 'danmaku-aoi-alpha-1';

UPDATE community_video_danmaku
SET author_name = 'Color Note'
WHERE id = 'danmaku-aoi-alpha-2';

UPDATE community_video_danmaku
SET author_name = 'Layout Notes'
WHERE id = 'danmaku-aoi-alpha-3';

UPDATE community_video_danmaku
SET body = '互动反馈很顺',
    author_name = 'Aoi Curator'
WHERE id = 'danmaku-aoi-alpha-4';

UPDATE community_video_danmaku
SET body = '这条路线很清楚'
WHERE id = 'danmaku-go-api-1';

UPDATE community_video_danmaku
SET body = '收藏入口好找',
    author_name = 'Layout Notes'
WHERE id = 'danmaku-go-api-2';

UPDATE community_video_comments
SET body = '首页信息密度轻了很多，动态、分类和最新投稿之间的节奏更顺了。',
    author_name = 'Layout Notes',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'comment-aoi-alpha-1';

UPDATE community_video_comments
SET author_name = 'Color Note',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'comment-aoi-alpha-2';

UPDATE community_video_comments
SET body = '评论、弹幕和收藏入口都很轻，适合先让大家自然交流。',
    author_name = 'Aoi Curator',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'comment-aoi-alpha-3';

UPDATE community_video_comments
SET body = '从投稿到收藏这一段路径很顺，适合作为新用户的第一条浏览路线。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'comment-go-api-1';

UPDATE community_dynamics
SET body = '今天把首页动态整理成更轻的阅读节奏，短更新和关联视频会自然连在一起。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'dynamic-rin-alpha';

UPDATE community_dynamics
SET author_name = 'Color Note',
    body = '这版视觉会保留 kirakira 式的细线、粉色强调和轻量动效，但内容密度会更适合长期视频社区。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'dynamic-design-sakura';

UPDATE community_dynamics
SET author_name = 'Aoi Curator',
    body = '从投稿到收藏这一段路径很顺，适合作为新用户的第一条浏览路线。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'dynamic-backend-contract';

UPDATE community_dynamics
SET author_name = 'Layout Notes',
    body = '手机上也能轻松阅读动态卡片，长句会自然换行，视频入口保持稳定比例。',
    updated_at = '2026-06-26 18:20:00'
WHERE id = 'dynamic-frontend-mobile';

-- +goose Down
SELECT 1;
