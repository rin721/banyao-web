package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/open-console/console-platform/internal/modules/community/model"
	authtypes "github.com/open-console/console-platform/types/auth"
)

type MediaStorage interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte, os.FileMode) error
	MkdirAll(string, os.FileMode) error
}

type VideoService interface {
	UploadSource(context.Context, authtypes.Principal, UploadSourceInput) (model.CommunitySubmissionUploadResult, error)
	CreateTranscodeJob(context.Context, authtypes.Principal, string, model.CreateCommunityVideoJobRequest) (model.CommunityVideoJobItem, error)
	ListJobs(context.Context, model.CommunityVideoJobFilter) (model.CommunityVideoJobPayload, error)
	GetJob(context.Context, string) (model.CommunityVideoJobItem, error)
	RetryJob(context.Context, authtypes.Principal, string) (model.CommunityVideoJobItem, error)
	GetAsset(context.Context, string) (VideoAsset, error)
	GetSourceAsset(context.Context, string) (VideoAsset, error)
}

type UploadSourceInput struct {
	Filename    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type VideoAsset struct {
	ContentType string
	Data        []byte
}

type VideoConfig struct {
	Mode          string
	LocalBasePath string
	LocalFSType   string
	Local         VideoLocalConfig
	HLS           VideoHLSConfig
	Cloud         VideoCloudConfig
}

type VideoLocalConfig struct {
	FFmpegPath    string
	FFprobePath   string
	OutputRoot    string
	SourceRoot    string
	PublicBaseURL string
}

type VideoHLSConfig struct {
	SegmentSeconds int
	Renditions     []VideoRenditionConfig
}

type VideoRenditionConfig struct {
	Label     string
	Width     int
	Height    int
	VideoKbps int
	AudioKbps int
}

type VideoCloudConfig struct {
	Provider       string
	ObjectStorage  string
	Bucket         string
	CDNBaseURL     string
	CallbackSecret string
}

type configuredVideoService struct {
	app      *service
	cfg      VideoConfig
	provider videoProvider
}

type videoProvider interface {
	Name() string
	Transcode(context.Context, model.CommunityVideoJob, model.CommunitySubmission, model.CommunityMediaAsset, model.CreateCommunityVideoJobRequest) (videoTranscodeResult, error)
}

type videoTranscodeResult struct {
	DurationSeconds  int
	ThumbnailURL     string
	MasterURL        string
	OutputStorageKey string
	Renditions       []model.CommunityVideoRendition
}

func newConfiguredVideoService(app *service, cfg VideoConfig) VideoService {
	cfg = normalizeVideoConfig(cfg)
	v := &configuredVideoService{app: app, cfg: cfg}
	switch cfg.Mode {
	case model.CommunityVideoProviderCloud:
		v.provider = cloudVideoProvider{cfg: cfg}
	default:
		v.provider = localVideoProvider{cfg: cfg, storage: app.cfg.Storage}
	}
	return v
}

func normalizeVideoConfig(cfg VideoConfig) VideoConfig {
	cfg.Mode = strings.ToLower(strings.TrimSpace(cfg.Mode))
	if cfg.Mode == "" {
		cfg.Mode = model.CommunityVideoProviderLocal
	}
	if strings.TrimSpace(cfg.Local.FFmpegPath) == "" {
		cfg.Local.FFmpegPath = "ffmpeg"
	}
	if strings.TrimSpace(cfg.Local.FFprobePath) == "" {
		cfg.Local.FFprobePath = "ffprobe"
	}
	if strings.TrimSpace(cfg.Local.SourceRoot) == "" {
		cfg.Local.SourceRoot = "community/sources"
	}
	if strings.TrimSpace(cfg.Local.OutputRoot) == "" {
		cfg.Local.OutputRoot = "community/hls"
	}
	if strings.TrimSpace(cfg.Local.PublicBaseURL) == "" {
		cfg.Local.PublicBaseURL = "/api/v1/public/community/hls"
	}
	if cfg.HLS.SegmentSeconds <= 0 {
		cfg.HLS.SegmentSeconds = 6
	}
	if len(cfg.HLS.Renditions) == 0 {
		cfg.HLS.Renditions = []VideoRenditionConfig{
			{Label: "360p", Width: 640, Height: 360, VideoKbps: 800, AudioKbps: 96},
			{Label: "720p", Width: 1280, Height: 720, VideoKbps: 2800, AudioKbps: 128},
			{Label: "1080p", Width: 1920, Height: 1080, VideoKbps: 5000, AudioKbps: 160},
		}
	}
	return cfg
}

func (v *configuredVideoService) UploadSource(ctx context.Context, principal authtypes.Principal, input UploadSourceInput) (model.CommunitySubmissionUploadResult, error) {
	if v.app == nil || v.app.repo == nil || v.app.cfg.Storage == nil {
		return model.CommunitySubmissionUploadResult{}, ErrStorageUnavailable
	}
	if _, err := communityAccountClientID(principal); err != nil {
		return model.CommunitySubmissionUploadResult{}, err
	}
	filename := cleanUploadFilename(input.Filename)
	if filename == "" || input.Reader == nil {
		return model.CommunitySubmissionUploadResult{}, ErrInvalidInput
	}
	data, err := io.ReadAll(input.Reader)
	if err != nil {
		return model.CommunitySubmissionUploadResult{}, err
	}
	if len(data) == 0 {
		return model.CommunitySubmissionUploadResult{}, ErrInvalidInput
	}
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := normalizeUploadedVideoMIME(input.ContentType, data)
	if !strings.HasPrefix(contentType, "video/") {
		contentType = videoMIMEFromExtension(ext)
	}
	if !strings.HasPrefix(contentType, "video/") {
		return model.CommunitySubmissionUploadResult{}, ErrInvalidInput
	}
	now := v.app.now()
	id := v.app.newMediaAssetID()
	if ext == "" {
		exts, _ := mime.ExtensionsByType(contentType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}
	if ext == "" {
		ext = ".bin"
	}
	storageKey := cleanStorageKey(v.cfg.Local.SourceRoot, strconv.FormatInt(id, 10)+ext)
	if err := v.app.cfg.Storage.MkdirAll(path.Dir(storageKey), 0755); err != nil {
		return model.CommunitySubmissionUploadResult{}, mapStorageError(err)
	}
	if err := v.app.cfg.Storage.WriteFile(storageKey, data, 0644); err != nil {
		return model.CommunitySubmissionUploadResult{}, mapStorageError(err)
	}
	asset := model.CommunityMediaAsset{
		ID:                 id,
		DisplayName:        filename,
		OriginalName:       filename,
		StorageKey:         storageKey,
		URL:                "/api/v1/public/community/source-assets/" + strconv.FormatInt(id, 10),
		MIMEType:           contentType,
		Extension:          strings.TrimPrefix(ext, "."),
		SizeBytes:          int64(len(data)),
		Source:             "upload",
		External:           false,
		UploadedByUsername: communityAccountAuthorName(principal),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := v.app.repo.CreateMediaAsset(ctx, asset); err != nil {
		return model.CommunitySubmissionUploadResult{}, mapStorageError(err)
	}
	return model.CommunitySubmissionUploadResult{
		MediaAssetID: asset.ID,
		DisplayName:  asset.DisplayName,
		OriginalName: asset.OriginalName,
		URL:          asset.URL,
		MIMEType:     asset.MIMEType,
		SizeBytes:    asset.SizeBytes,
	}, nil
}

func (v *configuredVideoService) CreateTranscodeJob(ctx context.Context, principal authtypes.Principal, submissionID string, req model.CreateCommunityVideoJobRequest) (model.CommunityVideoJobItem, error) {
	if v.app == nil || v.app.repo == nil {
		return model.CommunityVideoJobItem{}, ErrStorageUnavailable
	}
	reviewerID, err := communityReviewPrincipalID(principal)
	if err != nil {
		return model.CommunityVideoJobItem{}, err
	}
	submissionID = strings.TrimSpace(submissionID)
	if submissionID == "" {
		return model.CommunityVideoJobItem{}, ErrInvalidInput
	}
	submission, err := v.app.repo.FindCommunitySubmission(ctx, submissionID)
	if err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	if submission.Status != model.CommunitySubmissionStatusApproved {
		return model.CommunityVideoJobItem{}, ErrInvalidInput
	}
	if submission.MediaAssetID <= 0 {
		return model.CommunityVideoJobItem{}, ErrInvalidInput
	}
	asset, err := v.app.repo.FindMediaAssetByID(ctx, submission.MediaAssetID)
	if err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	now := v.app.now()
	job := model.CommunityVideoJob{
		ID:               v.app.newVideoJobID(),
		SubmissionID:     submission.ID,
		MediaAssetID:     asset.ID,
		Provider:         v.provider.Name(),
		Status:           model.CommunityVideoJobStatusQueued,
		Progress:         0,
		InputStorageKey:  asset.StorageKey,
		OutputStorageKey: cleanStorageKey(v.cfg.Local.OutputRoot, submissionVideoOutputDir(*submission)),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := v.app.repo.CreateCommunityVideoJob(ctx, job); err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	started := v.app.now()
	job.Status = model.CommunityVideoJobStatusRunning
	job.Progress = 5
	job.StartedAt = &started
	job.UpdatedAt = started
	if err := v.app.repo.UpdateCommunityVideoJob(ctx, job); err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	transcode, err := v.provider.Transcode(ctx, job, *submission, *asset, req)
	if err != nil {
		finished := v.app.now()
		job.Status = model.CommunityVideoJobStatusFailed
		job.Progress = 100
		job.ErrorMessage = trimRunes(err.Error(), 1200)
		job.FinishedAt = &finished
		job.UpdatedAt = finished
		if updateErr := v.app.repo.UpdateCommunityVideoJob(ctx, job); updateErr != nil {
			return model.CommunityVideoJobItem{}, mapStorageError(updateErr)
		}
		return v.decorateJob(ctx, job)
	}
	videoID, err := v.publishTranscodedSubmission(ctx, *submission, *asset, transcode, req, reviewerID)
	if err != nil {
		finished := v.app.now()
		job.Status = model.CommunityVideoJobStatusFailed
		job.Progress = 100
		job.ErrorMessage = trimRunes(err.Error(), 1200)
		job.FinishedAt = &finished
		job.UpdatedAt = finished
		if updateErr := v.app.repo.UpdateCommunityVideoJob(ctx, job); updateErr != nil {
			return model.CommunityVideoJobItem{}, mapStorageError(updateErr)
		}
		return model.CommunityVideoJobItem{}, err
	}
	for i := range transcode.Renditions {
		transcode.Renditions[i].JobID = job.ID
		transcode.Renditions[i].VideoID = videoID
	}
	if len(transcode.Renditions) > 0 {
		if err := v.app.repo.CreateCommunityVideoRenditions(ctx, transcode.Renditions); err != nil {
			return model.CommunityVideoJobItem{}, mapStorageError(err)
		}
	}
	finished := v.app.now()
	job.Status = model.CommunityVideoJobStatusSucceeded
	job.Progress = 100
	job.VideoID = videoID
	job.OutputStorageKey = transcode.OutputStorageKey
	job.OutputPublicURL = transcode.MasterURL
	job.ErrorMessage = ""
	job.FinishedAt = &finished
	job.UpdatedAt = finished
	if err := v.app.repo.UpdateCommunityVideoJob(ctx, job); err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	return v.decorateJob(ctx, job)
}

func (v *configuredVideoService) publishTranscodedSubmission(ctx context.Context, submission model.CommunitySubmission, asset model.CommunityMediaAsset, transcode videoTranscodeResult, req model.CreateCommunityVideoJobRequest, reviewerID string) (string, error) {
	now := v.app.now()
	durationSeconds := transcode.DurationSeconds
	if req.DurationSeconds > 0 {
		durationSeconds = normalizeSubmissionReviewDuration(req.DurationSeconds)
	}
	if durationSeconds <= 0 {
		durationSeconds = 1
	}
	videoID := submissionVideoID(submission)
	slug := submissionVideoSlug(submission, req.Slug)
	thumbnailURL := transcode.ThumbnailURL
	if strings.TrimSpace(req.ThumbnailURL) != "" {
		thumbnailURL = normalizeSubmissionReviewThumbnailURL(req.ThumbnailURL, slug)
	}
	if strings.TrimSpace(thumbnailURL) == "" {
		thumbnailURL = normalizeSubmissionReviewThumbnailURL("", slug)
	}
	description := trimRunes(submission.Description, 720)
	var descriptionPtr *string
	if description != "" {
		descriptionPtr = &description
	}
	creator, err := submissionVideoCreator(submission, now)
	if err != nil {
		return "", err
	}
	video := model.Video{
		ID:              videoID,
		Slug:            slug,
		Title:           trimRunes(submission.Title, 240),
		Description:     descriptionPtr,
		ThumbnailURL:    thumbnailURL,
		DurationSeconds: durationSeconds,
		SourceURL:       transcode.MasterURL,
		PublishedAt:     now,
		UploaderID:      creator.ID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	hlsMime := "application/vnd.apple.mpegurl"
	nativeMime := optionalSourceMimeType(asset.MIMEType)
	sources := []model.VideoSourceOption{
		{
			ID:           submissionVideoSourceID(videoID),
			VideoID:      videoID,
			Src:          transcode.MasterURL,
			Kind:         model.VideoSourceKindHLS,
			Label:        "HLS",
			MimeType:     &hlsMime,
			QualityLabel: stringPtr("auto"),
			IsDefault:    true,
			Order:        10,
		},
		{
			ID:        submissionVideoNativeSourceID(videoID),
			VideoID:   videoID,
			Src:       asset.URL,
			Kind:      model.VideoSourceKindNative,
			Label:     "MP4",
			MimeType:  nativeMime,
			IsDefault: false,
			Order:     20,
		},
	}
	categorySlugs := []string{}
	if strings.TrimSpace(submission.CategorySlug) != "" {
		categorySlugs = append(categorySlugs, strings.TrimSpace(submission.CategorySlug))
	}
	if err := v.app.repo.CreateVideoFromSubmissionSources(ctx, creator, video, sources, categorySlugs, decodeSubmissionTags(submission.TagsJSON)); err != nil {
		return "", mapStorageError(err)
	}
	if err := applySubmissionReview(&submission, model.CommunitySubmissionStatusPublished, submission.ReviewNote, reviewerID, videoID, asset.ID, now); err != nil {
		return "", err
	}
	if err := v.app.repo.UpdateCommunitySubmissionReview(ctx, submission); err != nil {
		return "", mapStorageError(err)
	}
	if err := v.app.createNotification(ctx, submissionReviewNotification(submission)); err != nil {
		return "", err
	}
	return videoID, nil
}

func (v *configuredVideoService) ListJobs(ctx context.Context, filter model.CommunityVideoJobFilter) (model.CommunityVideoJobPayload, error) {
	if v.app == nil || v.app.repo == nil {
		return model.CommunityVideoJobPayload{}, ErrStorageUnavailable
	}
	filter.Status = strings.TrimSpace(filter.Status)
	filter.Limit = normalizeLimit(filter.Limit, 48)
	jobs, err := v.app.repo.ListCommunityVideoJobs(ctx, filter)
	if err != nil {
		return model.CommunityVideoJobPayload{}, mapStorageError(err)
	}
	items := make([]model.CommunityVideoJobItem, 0, len(jobs))
	for _, job := range jobs {
		item, err := v.decorateJob(ctx, job)
		if err != nil {
			return model.CommunityVideoJobPayload{}, err
		}
		items = append(items, item)
	}
	return model.CommunityVideoJobPayload{Items: model.PageResult[model.CommunityVideoJobItem]{Items: items}}, nil
}

func (v *configuredVideoService) GetJob(ctx context.Context, jobID string) (model.CommunityVideoJobItem, error) {
	if v.app == nil || v.app.repo == nil {
		return model.CommunityVideoJobItem{}, ErrStorageUnavailable
	}
	job, err := v.app.repo.FindCommunityVideoJob(ctx, strings.TrimSpace(jobID))
	if err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	return v.decorateJob(ctx, *job)
}

func (v *configuredVideoService) RetryJob(ctx context.Context, principal authtypes.Principal, jobID string) (model.CommunityVideoJobItem, error) {
	job, err := v.app.repo.FindCommunityVideoJob(ctx, strings.TrimSpace(jobID))
	if err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	if job.Status != model.CommunityVideoJobStatusFailed {
		return model.CommunityVideoJobItem{}, ErrInvalidInput
	}
	return v.CreateTranscodeJob(ctx, principal, job.SubmissionID, model.CreateCommunityVideoJobRequest{})
}

func (v *configuredVideoService) GetAsset(_ context.Context, assetPath string) (VideoAsset, error) {
	if v.app == nil || v.app.cfg.Storage == nil {
		return VideoAsset{}, ErrStorageUnavailable
	}
	assetPath = strings.TrimPrefix(path.Clean("/"+strings.TrimSpace(assetPath)), "/")
	if assetPath == "" || strings.Contains(assetPath, "..") {
		return VideoAsset{}, ErrInvalidInput
	}
	key := cleanStorageKey(v.cfg.Local.OutputRoot, assetPath)
	data, err := v.app.cfg.Storage.ReadFile(key)
	if err != nil {
		return VideoAsset{}, mapStorageError(err)
	}
	return VideoAsset{ContentType: communityVideoAssetMIME(key), Data: data}, nil
}

func (v *configuredVideoService) GetSourceAsset(ctx context.Context, assetID string) (VideoAsset, error) {
	if v.app == nil || v.app.repo == nil || v.app.cfg.Storage == nil {
		return VideoAsset{}, ErrStorageUnavailable
	}
	id, err := strconv.ParseInt(strings.TrimSpace(assetID), 10, 64)
	if err != nil || id <= 0 {
		return VideoAsset{}, ErrInvalidInput
	}
	asset, err := v.app.repo.FindMediaAssetByID(ctx, id)
	if err != nil {
		return VideoAsset{}, mapStorageError(err)
	}
	data, err := v.app.cfg.Storage.ReadFile(asset.StorageKey)
	if err != nil {
		return VideoAsset{}, mapStorageError(err)
	}
	return VideoAsset{ContentType: firstNonEmpty(asset.MIMEType, communityVideoAssetMIME(asset.StorageKey)), Data: data}, nil
}

func (v *configuredVideoService) decorateJob(ctx context.Context, job model.CommunityVideoJob) (model.CommunityVideoJobItem, error) {
	renditions, err := v.app.repo.ListCommunityVideoRenditions(ctx, job.ID)
	if err != nil {
		return model.CommunityVideoJobItem{}, mapStorageError(err)
	}
	return model.CommunityVideoJobItem{
		ID:               job.ID,
		SubmissionID:     job.SubmissionID,
		MediaAssetID:     job.MediaAssetID,
		VideoID:          job.VideoID,
		Provider:         job.Provider,
		Status:           job.Status,
		Progress:         job.Progress,
		InputStorageKey:  job.InputStorageKey,
		OutputStorageKey: job.OutputStorageKey,
		OutputPublicURL:  job.OutputPublicURL,
		ErrorMessage:     job.ErrorMessage,
		Renditions:       renditions,
		StartedAt:        job.StartedAt,
		FinishedAt:       job.FinishedAt,
		CreatedAt:        job.CreatedAt,
		UpdatedAt:        job.UpdatedAt,
	}, nil
}

type localVideoProvider struct {
	cfg     VideoConfig
	storage MediaStorage
}

func (p localVideoProvider) Name() string { return model.CommunityVideoProviderLocal }

func (p localVideoProvider) Transcode(ctx context.Context, job model.CommunityVideoJob, submission model.CommunitySubmission, asset model.CommunityMediaAsset, req model.CreateCommunityVideoJobRequest) (videoTranscodeResult, error) {
	if p.storage == nil {
		return videoTranscodeResult{}, ErrStorageUnavailable
	}
	ffmpeg, err := resolveExecutable(p.cfg.Local.FFmpegPath)
	if err != nil {
		return videoTranscodeResult{}, fmt.Errorf("ffmpeg executable unavailable: %w", err)
	}
	ffprobe, err := resolveExecutable(p.cfg.Local.FFprobePath)
	if err != nil {
		return videoTranscodeResult{}, fmt.Errorf("ffprobe executable unavailable: %w", err)
	}
	inputPath, cleanup, err := p.localReadablePath(asset.StorageKey)
	if err != nil {
		return videoTranscodeResult{}, err
	}
	if cleanup != nil {
		defer cleanup()
	}
	outputName := submissionVideoOutputDir(submission)
	outputKey := cleanStorageKey(p.cfg.Local.OutputRoot, outputName)
	outputDir := p.localPath(outputKey)
	if err := os.RemoveAll(outputDir); err != nil {
		return videoTranscodeResult{}, err
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return videoTranscodeResult{}, err
	}
	duration := req.DurationSeconds
	if duration <= 0 {
		duration = probeDurationSeconds(ctx, ffprobe, inputPath)
	}
	coverPath := filepath.Join(outputDir, "cover.jpg")
	if err := runCommand(ctx, ffmpeg, "-y", "-ss", "00:00:01", "-i", inputPath, "-frames:v", "1", coverPath); err != nil {
		return videoTranscodeResult{}, fmt.Errorf("generate cover.jpg: %w", err)
	}
	renditions := make([]model.CommunityVideoRendition, 0, len(p.cfg.HLS.Renditions))
	for index, rendition := range p.cfg.HLS.Renditions {
		label := safeRenditionLabel(rendition.Label, index)
		renditionDir := filepath.Join(outputDir, label)
		if err := os.MkdirAll(renditionDir, 0755); err != nil {
			return videoTranscodeResult{}, err
		}
		segmentPattern := filepath.Join(renditionDir, "segment_%03d.ts")
		playlistPath := filepath.Join(renditionDir, "index.m3u8")
		filter := fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease", rendition.Width, rendition.Height)
		if err := runCommand(ctx, ffmpeg,
			"-y",
			"-i", inputPath,
			"-vf", filter,
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-b:v", strconv.Itoa(rendition.VideoKbps)+"k",
			"-c:a", "aac",
			"-b:a", strconv.Itoa(firstPositive(rendition.AudioKbps, 128))+"k",
			"-f", "hls",
			"-hls_time", strconv.Itoa(p.cfg.HLS.SegmentSeconds),
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", segmentPattern,
			playlistPath,
		); err != nil {
			return videoTranscodeResult{}, fmt.Errorf("generate %s rendition: %w", label, err)
		}
		playlistKey := cleanStorageKey(outputKey, label, "index.m3u8")
		renditions = append(renditions, model.CommunityVideoRendition{
			ID:           "rendition-" + shortHash(job.ID+":"+label),
			QualityLabel: label,
			Width:        rendition.Width,
			Height:       rendition.Height,
			BitrateKbps:  rendition.VideoKbps,
			PlaylistURL:  publicURL(p.cfg.Local.PublicBaseURL, path.Join(outputName, label, "index.m3u8")),
			StorageKey:   playlistKey,
			CreatedAt:    time.Now().UTC(),
		})
	}
	master := buildMasterPlaylist(p.cfg.HLS.Renditions)
	if err := os.WriteFile(filepath.Join(outputDir, "master.m3u8"), []byte(master), 0644); err != nil {
		return videoTranscodeResult{}, err
	}
	return videoTranscodeResult{
		DurationSeconds:  firstPositive(duration, 1),
		ThumbnailURL:     publicURL(p.cfg.Local.PublicBaseURL, path.Join(outputName, "cover.jpg")),
		MasterURL:        publicURL(p.cfg.Local.PublicBaseURL, path.Join(outputName, "master.m3u8")),
		OutputStorageKey: outputKey,
		Renditions:       renditions,
	}, nil
}

func (p localVideoProvider) localReadablePath(storageKey string) (string, func(), error) {
	localPath := p.localPath(storageKey)
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil, nil
	}
	data, err := p.storage.ReadFile(storageKey)
	if err != nil {
		return "", nil, mapStorageError(err)
	}
	tmp, err := os.CreateTemp("", "community-video-source-*"+filepath.Ext(storageKey))
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { _ = os.Remove(tmp.Name()) }
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return "", nil, err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return "", nil, err
	}
	return tmp.Name(), cleanup, nil
}

func (p localVideoProvider) localPath(storageKey string) string {
	storageKey = filepath.FromSlash(strings.TrimPrefix(path.Clean("/"+storageKey), "/"))
	base := strings.TrimSpace(p.cfg.LocalBasePath)
	if strings.TrimSpace(p.cfg.LocalFSType) == "basepath" && base != "" {
		return filepath.Join(base, storageKey)
	}
	if base != "" {
		return filepath.Join(base, storageKey)
	}
	return storageKey
}

type cloudVideoProvider struct {
	cfg VideoConfig
}

func (p cloudVideoProvider) Name() string { return model.CommunityVideoProviderCloud }

func (p cloudVideoProvider) Transcode(context.Context, model.CommunityVideoJob, model.CommunitySubmission, model.CommunityMediaAsset, model.CreateCommunityVideoJobRequest) (videoTranscodeResult, error) {
	provider := strings.TrimSpace(p.cfg.Cloud.Provider)
	if provider == "" {
		provider = "cloud"
	}
	return videoTranscodeResult{}, fmt.Errorf("cloud video provider %q is configured but no VOD adapter is installed", provider)
}

func cleanUploadFilename(value string) string {
	value = strings.TrimSpace(filepath.Base(value))
	value = strings.ReplaceAll(value, "\x00", "")
	if value == "." || value == string(filepath.Separator) {
		return ""
	}
	return trimRunes(value, 240)
}

func normalizeUploadedVideoMIME(contentType string, data []byte) string {
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	if strings.HasPrefix(strings.ToLower(contentType), "video/") || contentType == "application/vnd.apple.mpegurl" {
		return contentType
	}
	detected := http.DetectContentType(data)
	if strings.HasPrefix(strings.ToLower(detected), "video/") {
		return detected
	}
	return contentType
}

func videoMIMEFromExtension(ext string) string {
	switch strings.ToLower(strings.TrimSpace(ext)) {
	case ".mp4", ".m4v":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	default:
		return ""
	}
}

func cleanStorageKey(parts ...string) string {
	joined := path.Join(parts...)
	joined = strings.TrimPrefix(path.Clean("/"+joined), "/")
	if joined == "." || strings.HasPrefix(joined, "../") {
		return ""
	}
	return joined
}

func resolveExecutable(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidInput
	}
	if strings.ContainsAny(value, `/\`) {
		if _, err := os.Stat(value); err != nil {
			return "", err
		}
		return value, nil
	}
	return exec.LookPath(value)
}

func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message != "" {
			return fmt.Errorf("%w: %s", err, trimRunes(message, 800))
		}
		return err
	}
	return nil
}

func probeDurationSeconds(ctx context.Context, ffprobe string, inputPath string) int {
	cmd := exec.CommandContext(ctx, ffprobe, "-v", "error", "-show_entries", "format=duration", "-of", "default=nw=1:nk=1", inputPath)
	raw, err := cmd.Output()
	if err != nil {
		return 0
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
	if err != nil || value <= 0 {
		return 0
	}
	return int(value + 0.5)
}

func safeRenditionLabel(value string, index int) string {
	label := safeASCIIIdentifier(value)
	if label == "" {
		label = strconv.Itoa(index+1) + "p"
	}
	return label
}

func buildMasterPlaylist(renditions []VideoRenditionConfig) string {
	var builder strings.Builder
	builder.WriteString("#EXTM3U\n#EXT-X-VERSION:3\n")
	for index, rendition := range renditions {
		label := safeRenditionLabel(rendition.Label, index)
		bandwidth := (rendition.VideoKbps + firstPositive(rendition.AudioKbps, 128)) * 1000
		builder.WriteString("#EXT-X-STREAM-INF:BANDWIDTH=")
		builder.WriteString(strconv.Itoa(bandwidth))
		builder.WriteString(",RESOLUTION=")
		builder.WriteString(strconv.Itoa(rendition.Width))
		builder.WriteString("x")
		builder.WriteString(strconv.Itoa(rendition.Height))
		builder.WriteString("\n")
		builder.WriteString(label)
		builder.WriteString("/index.m3u8\n")
	}
	return builder.String()
}

func publicURL(base string, subpath string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	subpath = strings.TrimLeft(path.Clean("/"+subpath), "/")
	if base == "" {
		return "/" + subpath
	}
	return base + "/" + subpath
}

func communityVideoAssetMIME(key string) string {
	switch strings.ToLower(path.Ext(key)) {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/mp2t"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".mp4":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

func firstPositive(values ...int) int {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func stringPtr(value string) *string {
	return &value
}

func submissionVideoOutputDir(submission model.CommunitySubmission) string {
	raw := strings.TrimPrefix(submissionVideoID(submission), "video-")
	raw = safeASCIIIdentifier(raw)
	if raw == "" {
		raw = shortHash(submission.ID)
	}
	return trimRunes("video_"+raw, 96)
}

func submissionVideoNativeSourceID(videoID string) string {
	raw := strings.TrimPrefix(strings.TrimSpace(videoID), "video-")
	raw = safeASCIIIdentifier(raw)
	if raw == "" {
		raw = shortHash(videoID)
	}
	return trimRunes("source-"+raw+"-native", 96)
}

func (s *service) newMediaAssetID() int64 {
	if s.cfg.NewIntID == nil {
		return s.now().UnixNano()
	}
	if id := s.cfg.NewIntID(); id > 0 {
		return id
	}
	return s.now().UnixNano()
}

func (s *service) newVideoJobID() string {
	raw := strings.TrimSpace(s.cfg.NewID())
	if raw == "" {
		raw = strconv.FormatInt(s.now().UnixNano(), 10)
	}
	if strings.HasPrefix(raw, "video-job-") {
		return raw
	}
	return "video-job-" + raw
}

func (s *service) UploadCommunityAccountSubmissionSource(ctx context.Context, principal authtypes.Principal, input UploadSourceInput) (model.CommunitySubmissionUploadResult, error) {
	if s.video == nil {
		return model.CommunitySubmissionUploadResult{}, ErrStorageUnavailable
	}
	return s.video.UploadSource(ctx, principal, input)
}

func (s *service) CreateCommunitySubmissionTranscodeJob(ctx context.Context, principal authtypes.Principal, submissionID string, req model.CreateCommunityVideoJobRequest) (model.CommunityVideoJobItem, error) {
	if s.video == nil {
		return model.CommunityVideoJobItem{}, ErrStorageUnavailable
	}
	return s.video.CreateTranscodeJob(ctx, principal, submissionID, req)
}

func (s *service) ListCommunityVideoJobs(ctx context.Context, filter model.CommunityVideoJobFilter) (model.CommunityVideoJobPayload, error) {
	if s.video == nil {
		return model.CommunityVideoJobPayload{}, ErrStorageUnavailable
	}
	return s.video.ListJobs(ctx, filter)
}

func (s *service) GetCommunityVideoJob(ctx context.Context, jobID string) (model.CommunityVideoJobItem, error) {
	if s.video == nil {
		return model.CommunityVideoJobItem{}, ErrStorageUnavailable
	}
	return s.video.GetJob(ctx, jobID)
}

func (s *service) RetryCommunityVideoJob(ctx context.Context, principal authtypes.Principal, jobID string) (model.CommunityVideoJobItem, error) {
	if s.video == nil {
		return model.CommunityVideoJobItem{}, ErrStorageUnavailable
	}
	return s.video.RetryJob(ctx, principal, jobID)
}

func (s *service) GetCommunityVideoAsset(ctx context.Context, assetPath string) (VideoAsset, error) {
	if s.video == nil {
		return VideoAsset{}, ErrStorageUnavailable
	}
	return s.video.GetAsset(ctx, assetPath)
}

func (s *service) GetCommunitySourceAsset(ctx context.Context, assetID string) (VideoAsset, error) {
	if s.video == nil {
		return VideoAsset{}, ErrStorageUnavailable
	}
	return s.video.GetSourceAsset(ctx, assetID)
}
