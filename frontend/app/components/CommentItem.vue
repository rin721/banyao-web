<script setup lang="ts">
import type { CommentView } from "~/types/comments"

const props = defineProps<{
  comment: CommentView
}>()

const isEdited = computed(() => props.comment.updatedAt !== props.comment.createdAt)
const { locale, t } = useI18n()

function formatTime(value: string) {
  return new Date(value).toLocaleString(locale.value, {
    dateStyle: "medium",
    timeStyle: "short"
  })
}
</script>

<template>
  <AoiSurface as="article" class="comment-item" surface="card" padding="md">
    <div class="comment-item__avatar" aria-hidden="true">
      {{ comment.authorName.slice(0, 1).toUpperCase() }}
    </div>

    <div class="comment-item__content">
      <header class="comment-item__header">
        <div>
          <strong>{{ comment.authorName }}</strong>
          <span>{{ formatTime(comment.createdAt) }}</span>
          <small v-if="isEdited">{{ t("comments.item.edited") }}</small>
        </div>
      </header>

      <p class="comment-item__body">{{ comment.body }}</p>
    </div>
  </AoiSurface>
</template>

<style scoped>
.comment-item {
  display: grid;
  grid-template-columns: 38px minmax(0, 1fr);
  gap: 12px;
}

.comment-item__avatar {
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: var(--aoi-radius-sm);
  background: var(--aoi-accent-10);
  color: var(--aoi-accent-60);
  font-weight: 800;
}

.comment-item__content {
  display: grid;
  min-width: 0;
  gap: 10px;
}

.comment-item__header {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 12px;
}

.comment-item__header div:first-child {
  display: flex;
  min-width: 0;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px 8px;
}

.comment-item__header strong {
  color: var(--aoi-text);
}

.comment-item__header span,
.comment-item__header small {
  color: var(--aoi-text-muted);
  font-size: 12px;
}

.comment-item__body {
  margin: 0;
  color: var(--aoi-text);
  line-height: 1.75;
  white-space: pre-wrap;
  word-break: break-word;
}

@media (max-width: 620px) {
  .comment-item {
    grid-template-columns: 1fr;
  }

  .comment-item__header {
    flex-direction: column;
  }

}
</style>
