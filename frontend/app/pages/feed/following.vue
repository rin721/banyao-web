<script setup lang="ts">
const api = useAoiApi()
const following = useFollowingStore()
const { t } = useI18n()

const { data: feed, error, pending, refresh } = useAsyncData(
  "following-feed",
  () => api.getFollowingFeed(following.clientId || undefined),
  {
    immediate: false,
    server: false
  }
)
const recommendedCreators = computed(() => feed.value?.followingCount
  ? []
  : feed.value?.creators.filter((creator) => !following.isFollowing(creator.id)) || [])
const feedMessage = computed(() => following.syncError || feed.value?.message)

watch(feed, (value) => {
  if (value) {
    following.applyBackendFeed(value)
  }
})

onMounted(async () => {
  if (!following.hydrated) {
    following.restore()
  }
  await refresh()
})

useHead({
  title: "Following - Aoi"
})
</script>

<template>
  <div class="aoi-page">
    <PageHeader
      icon="radio-tower"
      :title="t('following.title')"
      :description="t('following.description')"
    />

    <PageState
      v-if="!pending && error"
      icon="cloud-alert"
      :title="t('following.errorTitle')"
      :description="t('following.errorDescription')"
      action-icon="refresh-cw"
      :action-label="t('following.retry')"
      @action="refresh()"
    />

    <template v-else-if="!pending && feed">
      <PageState
        v-if="!feed.authenticated && following.hydrated && following.followedCount === 0"
        icon="user-round-plus"
        :title="t('following.emptyTitle')"
        :description="feedMessage || t('following.emptyDescription')"
        action-icon="search"
        :action-label="t('following.searchAction')"
        @action="navigateTo('/search')"
      />

      <div
        v-if="following.hydrated && (following.followedList.length || following.latestVideos.length)"
        class="following-dashboard"
      >
        <AoiSection
          v-if="following.followedList.length"
          :title="t('following.followedTitle')"
          :description="feedMessage || t('following.followedDescription')"
          title-id="following-creators-title"
        >
          <template #actions>
            <AoiButton tone="accent" variant="outlined" size="sm" icon="refresh-cw" @click="refresh()">{{ t("following.refresh") }}</AoiButton>
          </template>
          <AoiContentGrid min-width="260px" gap="compact" :mobile-columns="1">
            <AoiReveal
              v-for="(creator, index) in following.followedList"
              :key="creator.id"
              class="following-card-reveal"
              :index="index"
            >
              <CreatorCard :creator="creator" />
            </AoiReveal>
          </AoiContentGrid>
        </AoiSection>

        <AoiSection
          v-if="following.latestVideos.length"
          :title="t('following.latestTitle')"
          title-id="following-latest-title"
        >
          <VideoGrid :videos="following.latestVideos" />
        </AoiSection>
      </div>

      <AoiSection
        v-if="recommendedCreators.length"
        :title="t('following.recommendedTitle')"
        :description="t('following.recommendedDescription')"
        title-id="following-recommended-title"
      >
        <template #actions>
          <AoiButton tone="accent" variant="outlined" size="sm" icon="search" to="/search">{{ t("following.exploreMore") }}</AoiButton>
        </template>
        <AoiContentGrid min-width="260px" gap="compact" :mobile-columns="1">
          <AoiReveal
            v-for="(creator, index) in recommendedCreators"
            :key="creator.id"
            class="following-card-reveal"
            :index="index"
          >
            <CreatorCard :creator="creator" />
          </AoiReveal>
        </AoiContentGrid>
      </AoiSection>

      <AoiSection
        v-if="!following.latestVideos.length && feed.latest.items.length"
        :title="t('following.recommendedLatestTitle')"
        title-id="following-recommended-latest-title"
      >
        <VideoGrid :videos="feed.latest.items" />
      </AoiSection>
    </template>

    <PageState
      v-else-if="!pending"
      icon="user-round-plus"
      :title="t('following.noContentTitle')"
      :description="t('following.noContentDescription')"
      action-icon="refresh-cw"
      :action-label="t('following.retry')"
      @action="refresh()"
    />
  </div>
</template>

<style scoped>
.following-card-reveal {
  min-width: 0;
}

.following-dashboard {
  display: grid;
  grid-template-columns: minmax(280px, 420px) minmax(0, 1fr);
  align-items: start;
  gap: 18px;
}

@media (max-width: 900px) {
  .following-dashboard {
    grid-template-columns: 1fr;
  }
}
</style>
