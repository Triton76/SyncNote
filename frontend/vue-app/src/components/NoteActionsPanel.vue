<script setup>
import { ref } from "vue";

const props = defineProps({
  busy: { type: Boolean, default: false },
  authed: { type: Boolean, default: false },
});

const emit = defineEmits(["create-note", "fetch-note"]);

const createForm = ref({ title: "", content: "" });
const queryId = ref("");

function createNote() {
  emit("create-note", { ...createForm.value }, () => {
    createForm.value = { title: "", content: "" };
  });
}

function fetchNote() {
  emit("fetch-note", queryId.value.trim());
}
</script>

<template>
  <section class="card">
    <h3>笔记操作</h3>
    <div class="form-grid">
      <input v-model="createForm.title" placeholder="笔记标题" />
      <textarea v-model="createForm.content" rows="4" placeholder="笔记内容" />
      <button :disabled="props.busy || !props.authed" @click="createNote">
        创建笔记
      </button>
    </div>

    <div class="inline-tools">
      <input v-model="queryId" placeholder="输入 note_id 查询" />
      <button :disabled="props.busy || !props.authed" @click="fetchNote">
        查询笔记
      </button>
    </div>
  </section>
</template>
