<script setup>
const props = defineProps({
  debugInfo: { type: Object, required: true },
});
</script>

<template>
  <section class="card debug-card">
    <h3>调试信息</h3>
    <pre>{{ JSON.stringify(props.debugInfo, null, 2) }}</pre>
  </section>
</template>
