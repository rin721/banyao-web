-- +goose Up
DELETE FROM community_notifications
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
)
   OR creator_id IN (
    'user-rin',
    'user-lab',
    'user-backend',
    'user-frontend',
    'user-design',
    'user-motion'
)
   OR target_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room',
    'dynamic-rin-alpha',
    'dynamic-design-sakura',
    'dynamic-backend-contract',
    'dynamic-frontend-mobile',
    'comment-aoi-alpha-1',
    'comment-aoi-alpha-2',
    'comment-aoi-alpha-3',
    'comment-go-api-1',
    'comment-mobile-grid-1'
);

DELETE FROM community_reports
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
)
   OR target_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room',
    'dynamic-rin-alpha',
    'dynamic-design-sakura',
    'dynamic-backend-contract',
    'dynamic-frontend-mobile',
    'comment-aoi-alpha-1',
    'comment-aoi-alpha-2',
    'comment-aoi-alpha-3',
    'comment-go-api-1',
    'comment-mobile-grid-1',
    'danmaku-aoi-alpha-1',
    'danmaku-aoi-alpha-2',
    'danmaku-aoi-alpha-3',
    'danmaku-aoi-alpha-4',
    'danmaku-go-api-1',
    'danmaku-go-api-2'
);

DELETE FROM community_video_history
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
);

DELETE FROM community_video_interactions
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
);

DELETE FROM community_creator_follows
WHERE creator_id IN (
    'user-rin',
    'user-lab',
    'user-backend',
    'user-frontend',
    'user-design',
    'user-motion'
);

DELETE FROM community_dynamics
WHERE id IN (
    'dynamic-rin-alpha',
    'dynamic-design-sakura',
    'dynamic-backend-contract',
    'dynamic-frontend-mobile'
)
   OR video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
)
   OR creator_id IN (
    'user-rin',
    'user-lab',
    'user-backend',
    'user-frontend',
    'user-design',
    'user-motion'
);

DELETE FROM community_video_comments
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
)
   OR id IN (
    'comment-aoi-alpha-1',
    'comment-aoi-alpha-2',
    'comment-aoi-alpha-3',
    'comment-go-api-1',
    'comment-mobile-grid-1'
);

DELETE FROM community_video_danmaku
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
)
   OR id IN (
    'danmaku-aoi-alpha-1',
    'danmaku-aoi-alpha-2',
    'danmaku-aoi-alpha-3',
    'danmaku-aoi-alpha-4',
    'danmaku-go-api-1',
    'danmaku-go-api-2'
);

DELETE FROM community_video_sources
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
)
   OR id IN (
    'source-aoi-alpha-primary',
    'source-aoi-alpha-secondary',
    'source-token-array-primary',
    'source-dark-mode-primary',
    'source-mobile-grid-primary',
    'source-go-api-primary',
    'source-sakura-primary',
    'source-music-primary',
    'source-game-primary'
);

DELETE FROM community_video_tags
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
);

DELETE FROM community_video_categories
WHERE video_id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
);

DELETE FROM community_videos
WHERE id IN (
    'video-aoi-alpha',
    'video-token-array',
    'video-dark-mode',
    'video-mobile-grid',
    'video-go-api',
    'video-sakura-accent',
    'video-music-stream',
    'video-game-room'
);

DELETE FROM community_creators
WHERE id IN (
    'user-rin',
    'user-lab',
    'user-backend',
    'user-frontend',
    'user-design',
    'user-motion'
);

-- +goose Down
SELECT 1;
