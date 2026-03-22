<script setup>
import { reactive, watch } from "vue";

const props = defineProps({
  note: { type: Object, default: null },
  busy: { type: Boolean, default: false },
});

const emit = defineEmits(["save", "remove", "refresh"]);

const draft = reactive({ title: "", content: "", version: 1 });

watch(
  () => props.note,
  (next) => {
    draft.title = next?.title || "";
    draft.content = next?.content || "";
    draft.version = Number(next?.version || 1);
  },
  { immediate: true },
);

function saveNote() {
  if (!props.note?.note_id) return;
  emit("save", {
    noteId: props.note.note_id,
    title: draft.title,
    content: draft.content,
    version: Number(draft.version || 1),
  });
}

function removeNote() {
  if (!props.note?.note_id) return;
  if (!window.confirm("确认删除该笔记？")) return;
  emit("remove", props.note.note_id);
}

function refreshNote() {
  if (!props.note?.note_id) return;
  emit("refresh", props.note.note_id);
}
</script>

<template>
  <section class="card detail-card">
    <h3>笔记详情</h3>

    <div v-if="!props.note" class="empty-box">
      <p class="muted">请在左侧点击一条笔记，或通过 note_id 查询后查看详情。</p>
    </div>

    <div v-else class="form-grid">
      <p class="muted">note_id: {{ props.note.note_id }}</p>
      <p class="muted">owner: {{ props.note.owner_id }}</p>
      <label>
        标题
        <input v-model="draft.title" />
      </label>
      <label>
        内容
        <textarea v-model="draft.content" rows="8" />
      </label>
      <label>
        版本号
        <input v-model.number="draft.version" type="number" min="1" />
      </label>
      <div class="row-actions">
        <button :disabled="props.busy" @click="saveNote">保存编辑</button>
        <button class="ghost" :disabled="props.busy" @click="refreshNote">
          刷新详情
        </button>
        <button class="danger" :disabled="props.busy" @click="removeNote">
          删除笔记
        </button>
      </div>
    </div>
  </section>
</template>
