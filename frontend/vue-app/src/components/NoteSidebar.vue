<script setup>
import { computed } from "vue";

const props = defineProps({
  notes: { type: Array, default: () => [] },
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 5 },
  selectedNoteId: { type: String, default: "" },
});

const emit = defineEmits(["update:page", "update:pageSize", "select"]);

const totalPages = computed(() =>
  Math.max(1, Math.ceil(props.notes.length / props.pageSize)),
);
const pagedNotes = computed(() => {
  const start = (props.page - 1) * props.pageSize;
  return props.notes.slice(start, start + props.pageSize);
});

function setPage(next) {
  const safePage = Math.max(1, Math.min(totalPages.value, Number(next) || 1));
  emit("update:page", safePage);
}

function prevPage() {
  setPage(props.page - 1);
}

function nextPage() {
  setPage(props.page + 1);
}

function onPageSizeChange(event) {
  emit("update:pageSize", Number(event.target.value));
}
</script>

<template>
  <aside class="card sidebar-card">
    <div class="sidebar-head">
      <h3>笔记列表</h3>
      <span class="muted">{{ props.notes.length }} 条</span>
    </div>

    <label>
      每页
      <select :value="props.pageSize" @change="onPageSizeChange">
        <option :value="5">5</option>
        <option :value="10">10</option>
        <option :value="20">20</option>
      </select>
    </label>

    <div class="sidebar-nav">
      <button :disabled="props.page <= 1" @click="prevPage">上一页</button>
      <button :disabled="props.page >= totalPages" @click="nextPage">
        下一页
      </button>
    </div>

    <div class="page-stack">
      <button
        v-for="n in totalPages"
        :key="`page-${n}`"
        class="chip"
        :class="{ active: n === props.page }"
        type="button"
        @click="setPage(n)"
      >
        第 {{ n }} 页
      </button>
    </div>

    <div class="note-mini-list" v-if="pagedNotes.length">
      <button
        v-for="item in pagedNotes"
        :key="item.note_id"
        class="note-mini"
        :class="{ active: item.note_id === props.selectedNoteId }"
        type="button"
        @click="emit('select', item.note_id)"
      >
        <strong>{{ item.title || "未命名笔记" }}</strong>
        <small>{{ item.note_id }}</small>
      </button>
    </div>

    <p v-else class="muted">暂无笔记</p>
  </aside>
</template>
