<script setup>
import { ref } from "vue";

const props = defineProps({
  busy: { type: Boolean, default: false },
  authed: { type: Boolean, default: false },
});

const emit = defineEmits(["grant-permission"]);

const form = ref({ noteId: "", targetUserId: "", level: 1 });

function submitGrant() {
  emit("grant-permission", { ...form.value }, () => {
    form.value = { ...form.value, targetUserId: "" };
  });
}
</script>

<template>
  <section class="card">
    <h3>授权接口调试</h3>
    <div class="form-grid">
      <input v-model="form.noteId" placeholder="note_id" />
      <input v-model="form.targetUserId" placeholder="target_user_id" />
      <select v-model.number="form.level">
        <option :value="1">read (1)</option>
        <option :value="2">write (2)</option>
        <option :value="3">admin (3)</option>
      </select>
      <button :disabled="props.busy || !props.authed" @click="submitGrant">
        授予权限
      </button>
    </div>
  </section>
</template>
