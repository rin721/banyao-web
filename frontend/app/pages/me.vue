<script setup lang="ts">
import type { AccountProfileResponse, CommunitySubmissionItem } from "~/types/api"
import type { AoiSegmentedItem } from "~/components/aoi/AoiSegmentedControl.vue"

const api = useAoiApi()
const authSession = useAuthSessionStore()
const { t } = useI18n()
const router = useRouter()

const authenticated = computed(() => authSession.authenticated)

watchEffect(() => {
  if (authSession.hydrated && !authenticated.value) {
    router.replace("/login")
  }
})

useHead({ title: `${t("me.headTitle")} - Aoi` })

type MeTab = "profile" | "security" | "submissions"
const activeTab = ref<MeTab>("profile")

const profile = ref<AccountProfileResponse | null>(null)
const profileError = ref<string | null>(null)
const profilePending = ref(false)

const submissionsList = ref<CommunitySubmissionItem[] | null>(null)
const submissionsError = ref<string | null>(null)
const submissionsPending = ref(false)

async function loadProfile() {
  if (!authenticated.value) return
  profilePending.value = true
  profileError.value = null
  try {
    profile.value = await api.getAccountProfile()
  } catch {
    profileError.value = t("me.loadError")
  } finally {
    profilePending.value = false
  }
}

watch(authenticated, (val) => { if (val) void loadProfile() }, { immediate: true })

const editingDisplayName = ref(false)
const displayNameInput = ref("")
const displayNameSaving = ref(false)
const displayNameMessage = ref<{ type: "success" | "error"; text: string } | null>(null)

function startEditDisplayName() {
  displayNameInput.value = profile.value?.displayName ?? ""
  editingDisplayName.value = true
  displayNameMessage.value = null
}
function cancelEditDisplayName() {
  editingDisplayName.value = false
  displayNameMessage.value = null
}
async function saveDisplayName() {
  const name = displayNameInput.value.trim()
  if (!name || displayNameSaving.value) return
  displayNameSaving.value = true
  displayNameMessage.value = null
  try {
    profile.value = await api.updateAccountProfile({ displayName: name })
    editingDisplayName.value = false
    displayNameMessage.value = { type: "success", text: t("me.saveSuccess") }
  } catch (err) {
    displayNameMessage.value = { type: "error", text: t("me.saveError") }
  } finally {
    displayNameSaving.value = false
  }
}

const editingCreatorProfile = ref(false)
const creatorBioInput = ref("")
const creatorAvatarInput = ref("")
const creatorSaving = ref(false)
const creatorMessage = ref<{ type: "success" | "error"; text: string } | null>(null)
const isCreator = computed(() => profile.value?.role === "creator")

function startEditCreatorProfile() {
  creatorBioInput.value = profile.value?.bio ?? ""
  creatorAvatarInput.value = profile.value?.avatarUrl ?? ""
  editingCreatorProfile.value = true
  creatorMessage.value = null
}
function cancelEditCreatorProfile() {
  editingCreatorProfile.value = false
  creatorMessage.value = null
}
async function saveCreatorProfile() {
  if (creatorSaving.value) return
  creatorSaving.value = true
  creatorMessage.value = null
  try {
    profile.value = await api.updateAccountCreatorProfile({
      bio: creatorBioInput.value.trim() || null,
      avatarUrl: creatorAvatarInput.value.trim() || null
    })
    editingCreatorProfile.value = false
    creatorMessage.value = { type: "success", text: t("me.saveSuccess") }
  } catch (err) {
    creatorMessage.value = { type: "error", text: t("me.saveError") }
  } finally {
    creatorSaving.value = false
  }
}

const currentPasswordInput = ref("")
const newPasswordInput = ref("")
const passwordSaving = ref(false)
const passwordMessage = ref<{ type: "success" | "error"; text: string } | null>(null)

async function changePassword() {
  const currentPassword = currentPasswordInput.value.trim()
  const newPassword = newPasswordInput.value.trim()
  if (!currentPassword || newPassword.length < 8 || passwordSaving.value) return
  passwordSaving.value = true
  passwordMessage.value = null
  try {
    await api.changeAccountPassword({ currentPassword, newPassword })
    currentPasswordInput.value = ""
    newPasswordInput.value = ""
    passwordMessage.value = { type: "success", text: t("me.passwordSuccess") }
  } catch (err) {
    passwordMessage.value = { type: "error", text: t("me.passwordError") }
  } finally {
    passwordSaving.value = false
  }
}

const loggingOut = ref(false)
async function logout() {
  if (loggingOut.value) return
  loggingOut.value = true
  try {
    await authSession.logout()
    router.replace("/")
  } finally {
    loggingOut.value = false
  }
}

function roleLabel(role: string) {
  return role === "creator" ? t("me.roleCreator") : t("me.roleRegistered")
}
function formatDate(iso: string | null | undefined) {
  if (!iso) return "-"
  try { return new Date(iso).toLocaleDateString() } catch { return iso }
}

async function loadSubmissions() {
  if (!authenticated.value) return
  submissionsPending.value = true
  submissionsError.value = null
  try {
    const payload = await api.getCommunityAccountSubmissions()
    submissionsList.value = payload.items?.items || []
  } catch (err) {
    submissionsError.value = t("me.loadError")
  } finally {
    submissionsPending.value = false
  }
}

watch(activeTab, (tab) => {
  if (tab === "submissions" && !submissionsList.value) {
    void loadSubmissions()
  }
})

function formatBytes(bytes: number) {
  if (bytes <= 0) return "0 B"
  const k = 1024
  const sizes = ["B", "KB", "MB", "GB", "TB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

function statusIntent(status: string): "success" | "warning" | "danger" | "info" {
  if (status === "published") return "success"
  if (status === "approved") return "info"
  if (status === "rejected") return "danger"
  return "warning"
}

function statusText(status: string): string {
  if (status === "published") return "已发布"
  if (status === "approved") return "已通过 (转码中)"
  if (status === "rejected") return "被驳回"
  return "审核中"
}

const tabItems = computed<AoiSegmentedItem[]>(() => [
  { value: "profile", label: t("me.tabs.profile"), icon: "circle-user-round" },
  { value: "security", label: t("me.tabs.security"), icon: "shield-alert" },
  { value: "submissions", label: t("me.tabs.submissions"), icon: "upload" }
])
</script>

<template>
  <div class="aoi-page me-page">
    <PageHeader
      icon="circle-user-round"
      :title="t('me.headTitle')"
      :description="profile ? `@${profile.handle} • ${roleLabel(profile.role)}` : null"
    >
      <template #actions>
        <AoiButton
          variant="outlined"
          tone="neutral"
          icon="log-out"
          :disabled="loggingOut"
          @click="logout"
        >
          {{ loggingOut ? t("me.loggingOut") : t("me.logout") }}
        </AoiButton>
      </template>
    </PageHeader>

    <div v-if="!authSession.hydrated || profilePending" class="me-loading-wrapper">
      <AoiProgress indeterminate />
      <span class="me-loading-text">{{ t("me.loading") }}</span>
    </div>

    <div v-else-if="profileError && !profile" class="me-error-wrapper">
      <AoiStatusMessage intent="danger" icon="alert-circle">
        {{ profileError }}
      </AoiStatusMessage>
      <AoiButton variant="filled" tone="accent" @click="loadProfile">
        {{ t("me.saveChanges") }}
      </AoiButton>
    </div>

    <div v-else-if="profile" class="me-layout">
      <!-- Left sidebar -->
      <aside class="me-sidebar">
        <AoiSurface surface="panel" padding="md" class="me-profile-card">
          <div class="me-profile-avatar-container">
            <img
              v-if="profile.avatarUrl"
              :src="profile.avatarUrl"
              :alt="profile.displayName"
              class="me-profile-avatar"
            />
            <div v-else class="me-profile-avatar-fallback">
              {{ profile.displayName.charAt(0).toUpperCase() }}
            </div>
          </div>
          <div class="me-profile-meta">
            <h3>{{ profile.displayName }}</h3>
            <p>@{{ profile.handle }}</p>
            <span class="me-role-pill">{{ roleLabel(profile.role) }}</span>
          </div>
        </AoiSurface>

        <AoiSegmentedControl
          v-model="activeTab"
          :items="tabItems"
          selection-role="tab"
        />
      </aside>

      <!-- Main Panel content -->
      <main class="me-main-content">
        <!-- Tab: Profile -->
        <div v-if="activeTab === 'profile'" class="me-tab-pane">
          <!-- Info display -->
          <AoiSurface surface="panel" padding="lg">
            <h2 class="me-pane-title">
              <AoiIcon name="info" :size="18" decorative />
              {{ t("me.profileSection") }}
            </h2>
            <div class="me-info-grid">
              <div class="me-info-field">
                <span class="me-field-name">{{ t("me.handle") }}</span>
                <span class="me-field-value">@{{ profile.handle }}</span>
              </div>
              <div class="me-info-field">
                <span class="me-field-name">{{ t("me.email") }}</span>
                <span class="me-field-value">{{ profile.email }}</span>
              </div>
              <div class="me-info-field">
                <span class="me-field-name">{{ t("me.role") }}</span>
                <span class="me-field-value">{{ roleLabel(profile.role) }}</span>
              </div>
              <div class="me-info-field">
                <span class="me-field-name">{{ t("me.createdAt") }}</span>
                <span class="me-field-value">{{ formatDate(profile.createdAt) }}</span>
              </div>
              <div v-if="profile.lastLoginAt" class="me-info-field">
                <span class="me-field-name">{{ t("me.lastLoginAt") }}</span>
                <span class="me-field-value">{{ formatDate(profile.lastLoginAt) }}</span>
              </div>
              <div v-if="profile.bio" class="me-info-field me-info-field--full">
                <span class="me-field-name">{{ t("me.bio") }}</span>
                <span class="me-field-value me-bio-text">{{ profile.bio }}</span>
              </div>
            </div>
          </AoiSurface>

          <!-- Edit Display Name -->
          <AoiSurface surface="panel" padding="lg">
            <h2 class="me-pane-title">
              <AoiIcon name="user" :size="18" decorative />
              {{ t("me.editProfile") }}
            </h2>
            <div v-if="!editingDisplayName" class="me-trigger-row">
              <p class="me-trigger-desc">
                {{ t("me.displayName") }}: <strong>{{ profile.displayName }}</strong>
              </p>
              <AoiButton variant="outlined" tone="accent" @click="startEditDisplayName">
                {{ t("me.editProfile") }}
              </AoiButton>
            </div>
            <div v-else class="me-edit-form-content">
              <AoiTextField
                v-model="displayNameInput"
                appearance="outlined"
                :label="t('me.displayName')"
                :max-length="120"
                @enter="saveDisplayName"
              />
              <div class="me-form-actions">
                <AoiButton
                  variant="filled"
                  tone="accent"
                  :disabled="displayNameSaving || !displayNameInput.trim()"
                  @click="saveDisplayName"
                >
                  {{ displayNameSaving ? t("me.saving") : t("me.saveChanges") }}
                </AoiButton>
                <AoiButton variant="plain" tone="neutral" @click="cancelEditDisplayName">
                  {{ t("me.cancel") }}
                </AoiButton>
              </div>
            </div>
            <AoiStatusMessage
              v-if="displayNameMessage"
              :intent="displayNameMessage.type === 'success' ? 'success' : 'danger'"
              icon="info"
              class="me-form-feedback"
            >
              {{ displayNameMessage.text }}
            </AoiStatusMessage>
          </AoiSurface>

          <!-- Edit Creator Profile -->
          <AoiSurface v-if="isCreator" surface="panel" padding="lg">
            <h2 class="me-pane-title">
              <AoiIcon name="sparkles" :size="18" decorative />
              {{ t("me.editCreatorProfile") }}
            </h2>
            <div v-if="!editingCreatorProfile" class="me-trigger-row">
              <div class="me-trigger-desc">
                <p><strong>{{ t("me.bio") }}:</strong> {{ profile.bio || "-" }}</p>
                <p class="me-avatar-url-desc">
                  <strong>{{ t("me.avatarUrl") }}:</strong> <span class="me-url-text">{{ profile.avatarUrl || "-" }}</span>
                </p>
              </div>
              <AoiButton variant="outlined" tone="accent" @click="startEditCreatorProfile">
                {{ t("me.editCreatorProfile") }}
              </AoiButton>
            </div>
            <div v-else class="me-edit-form-content">
              <AoiTextField
                v-model="creatorBioInput"
                appearance="outlined"
                :label="t('me.bio')"
                multiline
                :rows="3"
                :max-length="640"
              />
              <AoiTextField
                v-model="creatorAvatarInput"
                appearance="outlined"
                :label="t('me.avatarUrl')"
                :max-length="512"
              />
              <div class="me-form-actions">
                <AoiButton
                  variant="filled"
                  tone="accent"
                  :disabled="creatorSaving"
                  @click="saveCreatorProfile"
                >
                  {{ creatorSaving ? t("me.saving") : t("me.saveChanges") }}
                </AoiButton>
                <AoiButton variant="plain" tone="neutral" @click="cancelEditCreatorProfile">
                  {{ t("me.cancel") }}
                </AoiButton>
              </div>
            </div>
            <AoiStatusMessage
              v-if="creatorMessage"
              :intent="creatorMessage.type === 'success' ? 'success' : 'danger'"
              icon="info"
              class="me-form-feedback"
            >
              {{ creatorMessage.text }}
            </AoiStatusMessage>
          </AoiSurface>
        </div>

        <!-- Tab: Security -->
        <div v-else-if="activeTab === 'security'" class="me-tab-pane">
          <AoiSurface surface="panel" padding="lg">
            <h2 class="me-pane-title">
              <AoiIcon name="key-round" :size="18" decorative />
              {{ t("me.changePassword") }}
            </h2>
            <div class="me-password-form">
              <AoiTextField
                v-model="currentPasswordInput"
                appearance="outlined"
                type="password"
                :label="t('me.currentPassword')"
              />
              <AoiTextField
                v-model="newPasswordInput"
                appearance="outlined"
                type="password"
                :label="t('me.newPassword')"
              />
              <div class="me-form-actions">
                <AoiButton
                  variant="filled"
                  tone="accent"
                  :disabled="passwordSaving || !currentPasswordInput.trim() || newPasswordInput.length < 8"
                  @click="changePassword"
                >
                  {{ passwordSaving ? t("me.saving") : t("me.changePassword") }}
                </AoiButton>
              </div>
            </div>
            <AoiStatusMessage
              v-if="passwordMessage"
              :intent="passwordMessage.type === 'success' ? 'success' : 'danger'"
              icon="info"
              class="me-form-feedback"
            >
              {{ passwordMessage.text }}
            </AoiStatusMessage>
          </AoiSurface>
        </div>

        <!-- Tab: Submissions -->
        <div v-else-if="activeTab === 'submissions'" class="me-tab-pane">
          <!-- Loading -->
          <div v-if="submissionsPending" class="me-loading-wrapper">
            <AoiProgress indeterminate />
            <span class="me-loading-text">正在加载投稿列表...</span>
          </div>

          <!-- Error -->
          <div v-else-if="submissionsError" class="me-error-wrapper">
            <AoiStatusMessage intent="danger" icon="alert-circle">
              {{ submissionsError }}
            </AoiStatusMessage>
            <AoiButton variant="filled" tone="accent" @click="loadSubmissions">
              重新加载
            </AoiButton>
          </div>

          <!-- Content -->
          <template v-else>
            <div v-if="submissionsList && submissionsList.length > 0" class="me-submissions-grid">
              <AoiSurface
                v-for="sub in submissionsList"
                :key="sub.id"
                surface="card"
                padding="md"
                class="me-sub-card"
              >
                <div class="me-sub-card__header">
                  <div class="me-sub-card__title-row">
                    <h3 class="me-sub-card__title">{{ sub.title }}</h3>
                    <span class="me-sub-card__date">{{ formatDate(sub.createdAt) }}</span>
                  </div>
                  <p class="me-sub-card__desc">{{ sub.description || "无简介" }}</p>
                </div>

                <div class="me-sub-card__meta">
                  <div class="me-meta-item">
                    <span class="me-meta-label">分区</span>
                    <span class="me-meta-val">{{ sub.categorySlug }}</span>
                  </div>
                  <div class="me-meta-item">
                    <span class="me-meta-label">源文件</span>
                    <span class="me-meta-val me-sub-card__file">{{ sub.sourceName }} ({{ formatBytes(sub.sourceSize) }})</span>
                  </div>
                </div>

                <div class="me-sub-card__status">
                  <AoiStatusMessage
                    :intent="statusIntent(sub.status)"
                    icon="info"
                    class="me-sub-status-msg"
                  >
                    <div class="me-status-body">
                      <strong>{{ statusText(sub.status) }}</strong>
                      <span v-if="sub.status === 'rejected' && sub.reviewNote" class="me-reject-note">
                        原因: {{ sub.reviewNote }}
                      </span>
                    </div>
                  </AoiStatusMessage>
                </div>

                <div class="me-sub-card__actions" v-if="sub.status === 'published' && sub.publishedVideoId">
                  <AoiButton
                    variant="outlined"
                    tone="accent"
                    icon="play"
                    @click="navigateTo(`/video/${sub.publishedVideoId}`)"
                  >
                    播放视频
                  </AoiButton>
                </div>
              </AoiSurface>
            </div>

            <!-- Empty State -->
            <AoiSurface v-else surface="panel" padding="lg">
              <PageState
                icon="video"
                title="暂无投稿"
                description="您还没有发布过任何视频投稿，现在就可以上传您的第一个视频作品。"
                action-icon="upload"
                action-label="前往投稿"
                @action="navigateTo('/upload')"
              />
            </AoiSurface>
          </template>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.me-page {
  max-width: var(--aoi-content-max-width);
}

.me-loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 60px 0;
  color: var(--aoi-text-muted);
}

.me-loading-text {
  font-size: 0.95rem;
}

.me-error-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 40px 0;
}

.me-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: var(--aoi-grid-gap);
  align-items: start;
}

.me-sidebar :deep(.aoi-segmented) {
  grid-template-columns: 1fr;
}

.me-sidebar {
  display: flex;
  flex-direction: column;
  gap: var(--aoi-grid-gap-compact);
}

.me-profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.me-profile-avatar-container {
  width: 88px;
  height: 88px;
  border-radius: var(--aoi-radius-round);
  overflow: hidden;
  margin-bottom: var(--aoi-grid-gap-compact);
  background: var(--aoi-accent-10);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  border: 2px solid var(--aoi-surface-solid);
}

.me-profile-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.me-profile-avatar-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: 700;
  color: var(--aoi-accent-60);
}

.me-profile-meta h3 {
  margin: 0 0 4px;
  font-size: 1.1rem;
  font-weight: 750;
  color: var(--aoi-text);
}

.me-profile-meta p {
  margin: 0 0 var(--aoi-grid-gap-compact);
  font-size: 0.85rem;
  color: var(--aoi-text-muted);
}

.me-role-pill {
  display: inline-block;
  padding: 3px 12px;
  border-radius: var(--aoi-radius-round);
  background: var(--aoi-accent-10);
  color: var(--aoi-accent-60);
  font-size: 0.72rem;
  font-weight: 750;
}

.me-main-content {
  min-width: 0;
}

.me-tab-pane {
  display: grid;
  gap: var(--aoi-grid-gap);
}

.me-pane-title {
  margin: 0 0 var(--aoi-grid-gap);
  font-size: 1.1rem;
  font-weight: 750;
  color: var(--aoi-text);
  display: flex;
  align-items: center;
  gap: 8px;
}

.me-info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.me-info-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.me-info-field--full {
  grid-column: span 2;
}

.me-field-name {
  font-size: 0.72rem;
  font-weight: 750;
  color: var(--aoi-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.me-field-value {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--aoi-text);
}

.me-bio-text {
  line-height: 1.6;
  white-space: pre-wrap;
}

.me-trigger-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.me-trigger-desc {
  font-size: 0.95rem;
  color: var(--aoi-text);
  margin: 0;
}

.me-trigger-desc p {
  margin: 4px 0;
}

.me-avatar-url-desc {
  color: var(--aoi-text-muted);
  font-size: 0.85rem;
}

.me-url-text {
  font-family: monospace;
  background: var(--aoi-surface-muted);
  padding: 2px 6px;
  border-radius: var(--aoi-radius-xs);
  word-break: break-all;
}

.me-edit-form-content {
  display: grid;
  gap: 14px;
}

.me-form-actions {
  display: flex;
  gap: 8px;
}

.me-form-feedback {
  margin-top: 12px;
}

.me-password-form {
  display: grid;
  gap: 16px;
  max-width: 480px;
}

.me-submissions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--aoi-grid-gap);
}

.me-sub-card {
  display: flex;
  flex-direction: column;
  gap: var(--aoi-grid-gap-compact);
  height: 100%;
}

.me-sub-card__header {
  display: grid;
  gap: 4px;
}

.me-sub-card__title-row {
  display: flex;
  justify-content: space-between;
  align-items: first baseline;
  gap: 12px;
}

.me-sub-card__title {
  font-size: 1rem;
  font-weight: 750;
  color: var(--aoi-text);
  margin: 0;
  word-break: break-word;
}

.me-sub-card__date {
  font-size: 0.75rem;
  color: var(--aoi-text-muted);
  white-space: nowrap;
}

.me-sub-card__desc {
  font-size: 0.85rem;
  color: var(--aoi-text-muted);
  margin: 4px 0 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}

.me-sub-card__meta {
  display: grid;
  gap: 6px;
  border-top: 1px solid var(--aoi-border);
  padding-top: 8px;
}

.me-meta-item {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
}

.me-meta-label {
  color: var(--aoi-text-muted);
}

.me-meta-val {
  font-weight: 600;
  color: var(--aoi-text);
}

.me-sub-card__file {
  font-family: monospace;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.me-sub-card__status {
  margin-top: auto;
}

.me-sub-status-msg {
  padding: var(--aoi-row-padding) !important;
}

.me-status-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.me-reject-note {
  font-size: 0.75rem;
  opacity: 0.9;
}

.me-sub-card__actions {
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid var(--aoi-border);
  padding-top: 8px;
}

@media (max-width: 760px) {
  .me-layout {
    grid-template-columns: 1fr;
  }
  .me-info-grid {
    grid-template-columns: 1fr;
  }
  .me-info-field--full {
    grid-column: span 1;
  }
  .me-submissions-grid {
    grid-template-columns: 1fr;
  }
}
</style>