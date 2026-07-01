<script setup lang="ts">
const api = useAoiApi()
const { t } = useI18n()

const followingList = ref<any[] | null>(null)
const followingPending = ref(false)
const followingError = ref<string | null>(null)

async function loadFollowing() {
  followingPending.value = true
  followingError.value = null
  try {
    const payload = await api.getAccountFollowingFeed()
    followingList.value = payload.creators || []
  } catch {
    followingError.value = t("me.loadError")
  } finally {
    followingPending.value = false
  }
}

onMounted(() => {
  void loadFollowing()
})
</script>

<template>
  <div class="me-following-subpage">
    <div v-if="followingPending" class="me-loading-wrapper">
      <AoiProgress indeterminate />
      <span class="me-loading-text">正在加载关注列表...</span>
    </div>
    <div v-else-if="followingError" class="me-error-wrapper">
      <AoiStatusMessage intent="danger" icon="alert-circle">
        {{ followingError }}
      </AoiStatusMessage>
      <AoiButton variant="filled" tone="accent" @click="loadFollowing">
        重新加载
      </AoiButton>
    </div>
    <template v-else>
      <div v-if="followingList && followingList.length > 0" class="me-creators-grid">
        <CreatorCard
          v-for="creator in followingList"
          :key="creator.id"
          :creator="creator"
          density="compact"
        />
      </div>
      <AoiSurface v-else surface="panel" padding="lg">
        <PageState
          icon="users"
          title="暂无关注"
          description="你还没关注任何社区创作者。去首页发现有趣的创作团队吧。"
          action-icon="search"
          action-label="发现创作者"
          @action="navigateTo('/')"
        />
      </AoiSurface>
    </template>
  </div>
</template>

<style scoped>
.me-creators-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--aoi-grid-gap);
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
