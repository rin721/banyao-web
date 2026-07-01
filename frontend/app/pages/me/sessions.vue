<script setup lang="ts">
const api = useAoiApi()
const authSession = useAuthSessionStore()
const { t } = useI18n()

const sessionsList = ref<any[] | null>(null)
const sessionsPending = ref(false)
const sessionsError = ref<string | null>(null)

async function loadSessions() {
  sessionsPending.value = true
  sessionsError.value = null
  try {
    const payload = await api.getAccountSessions()
    sessionsList.value = payload.items || []
  } catch {
    sessionsError.value = t("me.loadError")
  } finally {
    sessionsPending.value = false
  }
}

function formatDate(iso: string | null | undefined) {
  if (!iso) return "-"
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

onMounted(() => {
  void loadSessions()
})
</script>

<template>
  <div class="me-sessions-subpage">
    <div v-if="sessionsPending" class="me-loading-wrapper">
      <AoiProgress indeterminate />
      <span class="me-loading-text">正在加载活跃会话...</span>
    </div>
    <div v-else-if="sessionsError" class="me-error-wrapper">
      <AoiStatusMessage intent="danger" icon="alert-circle">
        {{ sessionsError }}
      </AoiStatusMessage>
      <AoiButton variant="filled" tone="accent" @click="loadSessions">
        重新加载
      </AoiButton>
    </div>
    <template v-else>
      <div v-if="sessionsList && sessionsList.length > 0" class="me-sessions-list">
        <AoiSurface
          v-for="s in sessionsList"
          :key="s.id"
          surface="card"
          padding="md"
          class="me-session-card"
        >
          <div class="me-session-card__header">
            <div class="me-session-card__title">
              <AoiIcon :name="s.clientType.includes('mobile') ? 'smartphone' : 'monitor'" :size="18" decorative />
              <strong>{{ s.clientType }}</strong>
              <span v-if="s.id === authSession.session?.sessionId" class="me-current-session-badge">当前设备</span>
            </div>
            <span class="me-session-card__date">创建于 {{ formatDate(s.createdAt) }}</span>
          </div>
          <div class="me-session-card__body">
            <div class="me-session-detail">
              <span>IP 地址:</span>
              <strong>{{ s.ipAddress }}</strong>
            </div>
            <div class="me-session-detail">
              <span>User Agent:</span>
              <span class="me-ua-text">{{ s.userAgent }}</span>
            </div>
          </div>
        </AoiSurface>
      </div>
      <div v-else class="me-empty-sessions">
        没有找到活跃会话记录。
      </div>
    </template>
  </div>
</template>

<style scoped>
.me-sessions-list {
  display: grid;
  gap: var(--aoi-grid-gap-compact);
}
.me-session-card {
  display: grid;
  gap: 8px;
}
.me-session-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.me-session-card__title {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--aoi-text);
}
.me-current-session-badge {
  font-size: 11px;
  background: var(--aoi-accent-10);
  color: var(--aoi-accent-60);
  padding: 2px 8px;
  border-radius: var(--aoi-radius-round);
  font-weight: 750;
}
.me-session-card__date {
  font-size: 12px;
  color: var(--aoi-text-muted);
}
.me-session-card__body {
  font-size: 13px;
  display: grid;
  gap: 4px;
  border-top: 1px solid var(--aoi-border);
  padding-top: 8px;
}
.me-session-detail {
  display: flex;
  gap: 8px;
}
.me-session-detail span {
  color: var(--aoi-text-muted);
}
.me-session-detail strong {
  color: var(--aoi-text);
}
.me-ua-text {
  color: var(--aoi-text-muted);
  font-family: monospace;
  font-size: 11px;
  word-break: break-all;
}
.me-empty-sessions {
  text-align: center;
  padding: 48px 0;
  color: var(--aoi-text-muted);
}
.me-loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 48px 0;
}
.me-loading-text {
  font-size: 0.9rem;
  color: var(--aoi-text-muted);
}
.me-error-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 48px 0;
}
</style>
