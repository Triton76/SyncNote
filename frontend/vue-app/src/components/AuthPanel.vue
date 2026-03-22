<script setup>
import { ref } from "vue";

const props = defineProps({
  busy: { type: Boolean, default: false },
});

const emit = defineEmits(["login", "register"]);

const mode = ref("login");
const loginForm = ref({ email: "", password: "" });
const registerForm = ref({ email: "", password: "" });

function submitLogin() {
  emit("login", { ...loginForm.value });
}

function submitRegister() {
  emit("register", { ...registerForm.value });
}
</script>

<template>
  <section class="card auth-card">
    <h2>欢迎使用 SyncNote</h2>
    <p class="muted">先登录或注册，再进入 Dashboard 进行团队与笔记操作。</p>

    <div class="auth-switch">
      <button :class="{ active: mode === 'login' }" @click="mode = 'login'">
        登录
      </button>
      <button
        :class="{ active: mode === 'register' }"
        @click="mode = 'register'"
      >
        注册
      </button>
    </div>

    <div v-if="mode === 'login'" class="form-grid">
      <input v-model="loginForm.email" type="email" placeholder="邮箱" />
      <input v-model="loginForm.password" type="password" placeholder="密码" />
      <button :disabled="props.busy" @click="submitLogin">
        登录并进入 Dashboard
      </button>
    </div>

    <div v-else class="form-grid">
      <input v-model="registerForm.email" type="email" placeholder="邮箱" />
      <input
        v-model="registerForm.password"
        type="password"
        placeholder="密码"
      />
      <button :disabled="props.busy" @click="submitRegister">注册</button>
    </div>
  </section>
</template>
