<script setup lang="ts">
import type { AccountProfileResponse } from "~/types/api"
import type { AoiSegmentedItem } from "~/components/aoi/AoiSegmentedControl.vue"

const api = useAoiApi()
const authSession = useAuthSessionStore()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const authenticated = computed(() => authSession.authenticated)

watchEffect(() => {
  if (authSession.hydrated && !authenticated.value) {
    router.replace("/login")
  }
})

useHead({ title: `${t("me.headTitle")} - Aoi` })

const profile = ref<AccountProfileResponse | null>(null)
const profileError = ref<string | null>(null)
const profilePending = ref(false)

async function loadProfile() {
  if (!authenticated.value) return
  profilePending.value = true
  profileError.value = null
  try {
    profile.value = await api.getAccountProfile()
    authSession.profileAvatarUrl = profile.value.avatarUrl || ""
  } catch {
    profileError.value = t("me.loadError")
  } finally {
    profilePending.value = false
  }
}

watch(authenticated, (val) => {
  if (val) void loadProfile()
}, { immediate: true })

// Expose profile state to sub-routes
provide("meProfile", {
  profile,
  profilePending,
  profileError,
  loadProfile
})

// Segmented control active tab matches the sub-route path
const activeTab = computed({
  get: () => {
    const segments = route.path.split("/me/")
    return (segments[1] || "profile") as any
  },
  set: (val) => {
    navigateTo(`/me/${val}`)
  }
})

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

const tabItems = computed<AoiSegmentedItem[]>(() => [
  { value: "profile", label: t("me.tabs.profile"), icon: "circle-user-round" },
  { value: "following", label: t("me.tabs.following"), icon: "users" },
  { value: "collections", label: t("me.tabs.collections"), icon: "star" },
  { value: "history", label: t("me.tabs.history"), icon: "clock-3" },
  { value: "submissions", label: t("me.tabs.submissions"), icon: "upload" },
  { value: "security", label: t("me.tabs.security"), icon: "shield-alert" },
  { value: "sessions", label: t("me.tabs.sessions"), icon: "smartphone" }
])

// Watch route path changes to automatically redirect base route /me to /me/profile
watch([() => route.path, () => authSession.hydrated], () => {
  if (!authSession.hydrated) return
  if (route.path === "/me" || route.path === "/me/") {
    navigateTo("/me/profile", { replace: true })
  }
}, { immediate: true })
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
        <NuxtPage />
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

@media (max-width: 960px) {
  .me-layout {
    grid-template-columns: 1fr;
  }
}
</style>