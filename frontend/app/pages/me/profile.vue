<script setup lang="ts">
import type { AccountProfileResponse } from "~/types/api"

const api = useAoiApi()
const authSession = useAuthSessionStore()
const { t } = useI18n()

const { profile, loadProfile } = inject("meProfile") as {
  profile: Ref<AccountProfileResponse | null>
  loadProfile: () => Promise<void>
}

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
const uploadingAvatar = ref(false)

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

async function onAvatarCropped(result: any) {
  const file = new File([result.blob], "avatar.webp", { type: "image/webp" })
  uploadingAvatar.value = true
  creatorMessage.value = null
  try {
    const res = await api.uploadAccountAvatar(file)
    creatorAvatarInput.value = res.avatarUrl
    if (profile.value) {
      profile.value.avatarUrl = res.avatarUrl
    }
    if (authSession.session && authSession.session.account) {
      (authSession.session.account as any).avatarUrl = res.avatarUrl
    }
    creatorMessage.value = { type: "success", text: "头像上传并保存成功" }
  } catch (err) {
    creatorMessage.value = { type: "error", text: "头像上传失败，请重试" }
  } finally {
    uploadingAvatar.value = false
  }
}

async function deleteAvatar() {
  if (uploadingAvatar.value) return
  uploadingAvatar.value = true
  creatorMessage.value = null
  try {
    const res = await api.deleteAccountAvatar()
    creatorAvatarInput.value = ""
    if (profile.value) {
      profile.value.avatarUrl = ""
    }
    if (authSession.session && authSession.session.account) {
      (authSession.session.account as any).avatarUrl = ""
    }
    creatorMessage.value = { type: "success", text: "头像已成功删除" }
  } catch (err) {
    creatorMessage.value = { type: "error", text: "头像删除失败，请重试" }
  } finally {
    uploadingAvatar.value = false
  }
}
</script>

<template>
  <div v-if="profile" class="me-profile-subpage">
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
        <div class="me-info-field me-info-field--full">
          <span class="me-field-name">{{ t("me.displayName") }}</span>
          <div v-if="!editingDisplayName" class="me-trigger-row">
            <span class="me-field-value">{{ profile.displayName }}</span>
            <AoiButton variant="outlined" tone="accent" @click="startEditDisplayName">
              {{ t("me.editDisplayName") }}
            </AoiButton>
          </div>
          <div v-else class="me-edit-form-content">
            <AoiTextField
              v-model="displayNameInput"
              appearance="outlined"
              :label="t('me.displayName')"
              :max-length="96"
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
    <AoiSurface surface="panel" padding="lg">
      <h2 class="me-pane-title">
        <AoiIcon name="sparkles" :size="18" decorative />
        {{ t("me.avatarAndBio") }}
      </h2>
      <div v-if="!editingCreatorProfile" class="me-trigger-row">
        <div class="me-trigger-desc">
          <p><strong>{{ t("me.bio") }}:</strong> {{ profile.bio || "-" }}</p>
          <p class="me-avatar-url-desc">
            <strong>{{ t("me.avatarUrl") }}:</strong> <span class="me-url-text">{{ profile.avatarUrl || "-" }}</span>
          </p>
        </div>
        <div class="me-profile-action-group">
          <AoiButton variant="outlined" tone="accent" @click="startEditCreatorProfile">
            {{ t("me.editAvatarAndBio") }}
          </AoiButton>
          <AoiButton
            v-if="profile.avatarUrl"
            variant="outlined"
            tone="danger"
            :disabled="uploadingAvatar"
            @click="deleteAvatar"
          >
            {{ uploadingAvatar ? "删除中..." : "删除头像" }}
          </AoiButton>
        </div>
      </div>
      <div v-else class="me-edit-form-content">
        <div class="me-avatar-uploader">
          <AoiImageClipboard
            label="剪切并上传头像 (WebP)"
            aspect-ratio="1:1"
            :aspect-ratios="[{ value: '1:1', label: '1:1 正方形' }]"
            @result="onAvatarCropped"
          />
        </div>
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
</template>

<style scoped>
.me-profile-subpage {
  display: grid;
  gap: 16px;
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
  font-size: 0.8rem;
  color: var(--aoi-text-muted);
}
.me-field-value {
  font-size: 1rem;
  color: var(--aoi-text);
  word-break: break-all;
}
.me-trigger-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}
.me-profile-action-group {
  display: flex;
  gap: 8px;
}
.me-trigger-desc {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}
.me-avatar-url-desc {
  font-size: 0.85rem;
  color: var(--aoi-text-muted);
}
.me-url-text {
  font-family: monospace;
  word-break: break-all;
}
.me-edit-form-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 12px;
}
.me-form-actions {
  display: flex;
  gap: 8px;
}
.me-form-feedback {
  margin-top: 16px;
}
.me-avatar-uploader {
  display: grid;
  gap: 12px;
  border: 1px dashed var(--aoi-border);
  padding: 16px;
  border-radius: var(--aoi-radius-sm);
  background: var(--aoi-surface-muted);
}
@media (max-width: 760px) {
  .me-info-grid {
    grid-template-columns: 1fr;
  }
  .me-info-field--full {
    grid-column: span 1;
  }
}
</style>
