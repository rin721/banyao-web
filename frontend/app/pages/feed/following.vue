<script setup lang="ts">
const api = useAoiApi()
const following = useFollowingStore()
const { t } = useI18n()
const dynamicAuthorName = ref(t("dynamics.composer.defaultAuthor"))
const dynamicError = ref("")
const dynamicSubmitRevision = ref(0)
const dynamicSubmitting = ref(false)

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
const dynamicItems = computed(() => feed.value?.dynamics.items || [])

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

async function publishDynamic(body: string) {
  dynamicError.value = ""
  dynamicSubmitting.value = true

  try {
    await api.createCommunityDynamic({
      authorName: dynamicAuthorName.value.trim(),
      body,
      clientId: following.ensureClientId()
    })
    dynamicSubmitRevision.value += 1
    await refresh()
  } catch (error) {
    dynamicError.value = error instanceof Error
      ? error.message
      : t("dynamics.composer.error")
  } finally {
    dynamicSubmitting.value = false
  }
}

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

      <AoiSection
        icon="sparkles"
        :title="t('following.dynamicComposerTitle')"
        :description="t('following.dynamicComposerDescription')"
        title-id="following-dynamic-composer-title"
      >
        <CommentComposer
          v-model:author-name="dynamicAuthorName"
          :author-label="t('dynamics.composer.authorLabel')"
          :body-label="t('dynamics.composer.bodyLabel')"
          :body-placeholder="t('dynamics.composer.bodyPlaceholder')"
          :hint="t('dynamics.composer.hint')"
          :submit-label="t('dynamics.composer.submit')"
          :error-text="dynamicError"
          :submitting="dynamicSubmitting"
          :submit-revision="dynamicSubmitRevision"
          :max-body-length="280"
          @submit="publishDynamic"
        />
      </AoiSection>

      <CommunityPulse
        :items="dynamicItems"
        :title="t('following.dynamicsTitle')"
        :description="feedMessage || t('following.dynamicsDescription')"
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
