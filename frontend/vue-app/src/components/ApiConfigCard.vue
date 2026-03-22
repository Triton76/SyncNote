<script setup>
const props = defineProps({
  authBase: { type: String, default: "" },
  syncBase: { type: String, default: "" },
  busy: { type: Boolean, default: false },
});

const emit = defineEmits(["update:authBase", "update:syncBase", "save"]);

function onInputAuth(event) {
  emit("update:authBase", event.target.value);
}

function onInputSync(event) {
  emit("update:syncBase", event.target.value);
}
</script>

<template>
  <section class="card">
    <h3>API 配置</h3>
    <label>
      Auth Base
      <input
        :value="props.authBase"
        placeholder="(留空=走前端代理)"
        @input="onInputAuth"
      />
    </label>
    <label>
      SyncNote Base
      <input
        :value="props.syncBase"
        placeholder="(留空=走前端代理)"
        @input="onInputSync"
      />
    </label>
    <button :disabled="props.busy" @click="emit('save')">保存地址</button>
  </section>
</template>
