<script setup lang="ts">
const library = useLibraryStore()
const activeTab = ref("favorites")
const syncPending = ref(false)

const tabItems = [
  { icon: "star", label: "收藏", value: "favorites" },
  { icon: "clock-3", label: "稍后看", value: "watchLater" }
]

const activeVideos = computed(() => activeTab.value === "favorites"
  ? library.favoriteList
  : library.watchLaterList)

const hasVideos = computed(() => library.hydrated && activeVideos.value.length > 0)
const clearLabel = computed(() => activeTab.value === "favorites" ? "清空收藏" : "清空稍后看")
const emptyTitle = computed(() => activeTab.value === "favorites" ? "还没有收藏" : "稍后看为空")
const emptyDescription = computed(() => activeTab.value === "favorites"
  ? "在视频详情页点击收藏后会写入后端社区资料库，并在这里同步。"
  : "把暂时没时间看的视频加入稍后看，就能从后端社区资料库继续找回。")

async function syncLibrary() {
  if (!library.hydrated || syncPending.value) {
    return
  }

  syncPending.value = true
  try {
    await library.syncWithBackend()
  } finally {
    syncPending.value = false
  }
}

async function clearActiveList() {
  if (activeTab.value === "favorites") {
    await library.clearFavorites()
    return
  }

  await library.clearWatchLater()
}

watch(() => library.hydrated, (hydrated) => {
  if (hydrated) {
    syncLibrary()
  }
}, { immediate: true })

useHead({
  title: "Collections - Aoi"
})
</script>

<template>
  <div class="aoi-page">
    <PageHeader
      icon="star"
      title="收藏"
      description="收藏和稍后看会写入 Go 后端社区匿名资料库；本地只保留离线缓存和加载降级。"
    >
      <template #actions>
        <AoiButton tone="accent"
          variant="outlined"
          icon="trash-2"
          :disabled="!hasVideos || syncPending"
          :loading="syncPending"
          @click="clearActiveList"
        >
          {{ clearLabel }}
        </AoiButton>
      </template>
    </PageHeader>

    <AoiReveal variant="fade">
      <AoiTabs
        v-model="activeTab"
        :items="tabItems"
        aria-label="收藏分类"
      />
    </AoiReveal>

    <p v-if="library.syncError" class="collections-sync-note">
      {{ library.syncError }}
    </p>

    <AoiSection v-if="hasVideos" :reveal="false">
      <VideoGrid :videos="activeVideos" />
    </AoiSection>

    <PageState
      v-else-if="library.hydrated && !syncPending"
      icon="star"
      :title="emptyTitle"
      :description="emptyDescription"
      action-icon="search"
      action-label="去搜索"
      @action="navigateTo('/search')"
    />

    <PageState
      v-else
      icon="refresh-cw"
      title="正在同步资料库"
      description="正在从后端社区模块读取收藏和稍后看。"
    />
  </div>
</template>

<style scoped>
.collections-sync-note {
  margin: 0;
  border: 1px solid var(--aoi-border);
  border-radius: var(--aoi-radius-sm);
  background: var(--aoi-card-bg);
  color: var(--aoi-text-muted);
  font-size: 13px;
  line-height: 1.5;
  padding: 10px 12px;
}
</style>
