<script setup lang="ts">
import type { CommunityNotificationItem } from "~/types/api"

const api = useAoiApi()
const library = useLibraryStore()
const { t } = useI18n()
const markReadPending = ref(false)

const { data: payload, error, pending, refresh } = useAsyncData(
  "community-notifications",
  () => api.getCommunityNotifications(library.ensureClientId()),
  {
    immediate: false,
    server: false
  }
)

const notifications = computed(() => payload.value?.items.items || [])
const unreadCount = computed(() => payload.value?.unreadCount || 0)
const hasNotifications = computed(() => notifications.value.length > 0)
const statusLabel = computed(() => unreadCount.value > 0
  ? t("notifications.unreadStatus", { count: unreadCount.value })
  : t("notifications.readStatus"))
const sourceMessage = computed(() => payload.value?.message || t("notifications.defaultMessage"))

onMounted(async () => {
  if (!library.hydrated) {
    library.restore()
  }
  await refresh()
})

async function markAllRead() {
  if (markReadPending.value || unreadCount.value === 0) {
    return
  }
  markReadPending.value = true
  try {
    payload.value = await api.markCommunityNotificationsRead({
      clientId: library.ensureClientId()
    })
  } finally {
    markReadPending.value = false
  }
}

function iconFor(item: CommunityNotificationItem) {
  if (item.kind === "comment") {
    return "message-circle"
  }
  if (item.kind === "danmaku") {
    return "message-square-text"
  }
  if (item.kind === "follow") {
    return "user-round-check"
  }
  if (item.kind === "report") {
    return "shield-check"
  }
  return "sparkles"
}

function formatCreatedAt(value: string) {
  return new Intl.DateTimeFormat("zh-CN", {
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    month: "2-digit"
  }).format(new Date(value))
}

useHead({
  title: "Notifications - Aoi"
})
</script>

<template>
  <div class="aoi-page notifications-page">
    <PageHeader
      icon="bell"
      :title="t('notifications.title')"
      :description="t('notifications.description')"
    >
      <template #actions>
        <AoiButton
          tone="accent"
          variant="outlined"
          icon="check-check"
          :disabled="unreadCount === 0 || markReadPending"
          :loading="markReadPending"
          @click="markAllRead"
        >
          {{ t("notifications.markAllRead") }}
        </AoiButton>
      </template>
    </PageHeader>

    <PageState
      v-if="!pending && error"
      icon="cloud-alert"
      :title="t('notifications.errorTitle')"
      :description="t('notifications.errorDescription')"
      action-icon="refresh-cw"
      :action-label="t('notifications.retry')"
      @action="refresh()"
    />

    <template v-else-if="!pending && payload">
      <section class="notifications-summary" :aria-label="t('notifications.summaryLabel')">
        <div class="notifications-summary__metric">
          <span class="notifications-summary__value">{{ unreadCount }}</span>
          <span class="notifications-summary__label">{{ t("notifications.unreadLabel") }}</span>
        </div>
        <div class="notifications-summary__copy">
          <p class="notifications-summary__title">{{ statusLabel }}</p>
          <p class="notifications-summary__message">{{ sourceMessage }}</p>
        </div>
        <AoiButton
          class="notifications-summary__mobile-action"
          tone="accent"
          variant="outlined"
          icon="check-check"
          :disabled="unreadCount === 0 || markReadPending"
          :loading="markReadPending"
          @click="markAllRead"
        >
          {{ t("notifications.markAllRead") }}
        </AoiButton>
      </section>

      <AoiSection
        v-if="hasNotifications"
        :title="t('notifications.listTitle')"
        :description="t('notifications.listDescription')"
        :count="notifications.length"
        title-id="notifications-list-title"
      >
        <div class="notifications-list" role="list">
          <article
            v-for="(item, index) in notifications"
            :key="item.id"
            v-aoi-reveal="{ delay: index * 32 }"
            class="notification-card"
            :class="{ 'notification-card--unread': !item.readAt }"
            role="listitem"
          >
            <div class="notification-card__icon" aria-hidden="true">
              <AoiIcon :name="iconFor(item)" :size="18" decorative />
            </div>
            <div class="notification-card__body">
              <div class="notification-card__meta">
                <span class="notification-card__kind">{{ t(`notifications.kind.${item.kind}`) }}</span>
                <time :datetime="item.createdAt">{{ formatCreatedAt(item.createdAt) }}</time>
              </div>
              <h2 class="notification-card__title">{{ item.title }}</h2>
              <p class="notification-card__text">{{ item.body }}</p>
              <AoiButton
                v-if="item.link"
                class="notification-card__link"
                tone="accent"
                variant="plain"
                size="sm"
                trailing-icon="arrow-right"
                :to="item.link"
              >
                {{ t("notifications.openTarget") }}
              </AoiButton>
            </div>
            <span v-if="!item.readAt" class="notification-card__dot" :aria-label="t('notifications.unreadDot')" />
          </article>
        </div>
      </AoiSection>

      <PageState
        v-else
        icon="bell-off"
        :title="t('notifications.emptyTitle')"
        :description="t('notifications.emptyDescription')"
        action-icon="home"
        :action-label="t('notifications.exploreAction')"
        @action="navigateTo('/')"
      />
    </template>

    <PageState
      v-else-if="!pending"
      icon="bell"
      :title="t('notifications.noContentTitle')"
      :description="t('notifications.noContentDescription')"
      action-icon="refresh-cw"
      :action-label="t('notifications.retry')"
      @action="refresh()"
    />
  </div>
</template>

<style scoped>
.notifications-page {
  gap: 20px;
}

.notifications-summary {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 16px;
  align-items: center;
  border: 1px solid var(--aoi-state-border-active);
  border-radius: var(--aoi-radius-sm);
  background:
    linear-gradient(135deg, var(--aoi-accent-10), rgba(255, 255, 255, 0.92) 54%),
    var(--aoi-surface);
  box-shadow: var(--aoi-shadow-sm);
  padding: 16px;
}

.notifications-summary__metric {
  display: grid;
  min-width: 82px;
  place-items: center;
  border-right: 1px solid var(--aoi-border);
  padding-right: 16px;
}

.notifications-summary__value {
  color: var(--aoi-accent-60);
  font-size: 30px;
  font-weight: 850;
  line-height: 1;
}

.notifications-summary__label {
  margin-top: 4px;
  color: var(--aoi-text-muted);
  font-size: 12px;
  font-weight: 800;
}

.notifications-summary__copy {
  display: grid;
  min-width: 0;
  gap: 5px;
}

.notifications-summary__mobile-action {
  display: none;
  justify-self: start;
}

.notifications-summary__title,
.notifications-summary__message {
  margin: 0;
}

.notifications-summary__title {
  color: var(--aoi-text);
  font-size: 17px;
  font-weight: 800;
}

.notifications-summary__message {
  color: var(--aoi-text-muted);
  line-height: 1.7;
}

.notifications-list {
  display: grid;
  gap: 10px;
}

.notification-card {
  position: relative;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 12px;
  align-items: start;
  border: 1px solid var(--aoi-border);
  border-radius: var(--aoi-radius-sm);
  background: var(--aoi-surface);
  box-shadow: var(--aoi-shadow-sm);
  padding: 14px;
}

.notification-card--unread {
  border-color: var(--aoi-state-border-active);
  background: color-mix(in srgb, var(--aoi-accent-10) 48%, var(--aoi-surface-solid));
}

.notification-card__icon {
  display: grid;
  width: 36px;
  height: 36px;
  place-items: center;
  border-radius: var(--aoi-radius-sm);
  background: var(--aoi-accent-10);
  color: var(--aoi-accent-60);
}

.notification-card__body {
  display: grid;
  min-width: 0;
  gap: 7px;
}

.notification-card__meta {
  display: flex;
  min-width: 0;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  color: var(--aoi-text-muted);
  font-size: 12px;
}

.notification-card__kind {
  color: var(--aoi-accent-60);
  font-weight: 800;
}

.notification-card__title,
.notification-card__text {
  margin: 0;
}

.notification-card__title {
  color: var(--aoi-text);
  font-size: 16px;
  line-height: 1.35;
}

.notification-card__text {
  color: var(--aoi-text-muted);
  line-height: 1.7;
  overflow-wrap: anywhere;
}

.notification-card__link {
  justify-self: start;
}

.notification-card__dot {
  position: absolute;
  top: 14px;
  right: 14px;
  width: 8px;
  height: 8px;
  border-radius: var(--aoi-radius-round);
  background: var(--aoi-accent-60);
  box-shadow: 0 0 0 4px var(--aoi-accent-10);
}

@media (max-width: 639px) {
  .notifications-summary {
    grid-template-columns: 1fr;
  }

  .notifications-summary__metric {
    min-width: 0;
    justify-items: start;
    border-right: 0;
    border-bottom: 1px solid var(--aoi-border);
    padding-right: 0;
    padding-bottom: 12px;
  }

  .notifications-summary__mobile-action {
    display: inline-flex;
  }

  .notification-card {
    grid-template-columns: 1fr;
  }

  .notification-card__icon {
    width: 34px;
    height: 34px;
  }
}
</style>
