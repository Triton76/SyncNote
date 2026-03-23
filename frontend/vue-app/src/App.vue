<script setup>
import { computed, onMounted, ref, watch } from "vue";

const authBase = ref(
  localStorage.getItem("syncnote.authBase") ||
    import.meta.env.VITE_AUTH_BASE ||
    "http://127.0.0.1:8000",
);
const syncBase = ref(
  localStorage.getItem("syncnote.syncBase") ||
    import.meta.env.VITE_SYNC_BASE ||
    "http://127.0.0.1:8001",
);

const token = ref(localStorage.getItem("syncnote.token") || "");
const userId = ref(localStorage.getItem("syncnote.userId") || "");

const authMode = ref("login");
const registerForm = ref({ email: "", password: "" });
const loginForm = ref({ email: "", password: "" });
const createNoteForm = ref({ title: "", content: "" });
const loadNoteId = ref("");
const createTeamForm = ref({ name: "", description: "" });
const joinTeamId = ref("");
const grantForm = ref({ noteId: "", targetUserId: "", level: 1 });

const notes = ref([]);
const teams = ref([]);

const page = ref(1);
const pageSize = ref(5);

const busy = ref(false);
const message = ref("");

const totalPages = computed(() =>
  Math.max(1, Math.ceil(notes.value.length / pageSize.value)),
);
const pagedNotes = computed(() => {
  const start = (page.value - 1) * pageSize.value;
  return notes.value.slice(start, start + pageSize.value);
});

watch(pageSize, () => {
  if (page.value > totalPages.value) {
    page.value = totalPages.value;
  }
});

watch(token, (next) => {
  if (next) {
    localStorage.setItem("syncnote.token", next);
  } else {
    localStorage.removeItem("syncnote.token");
  }
});

watch(userId, (next) => {
  if (next) {
    localStorage.setItem("syncnote.userId", next);
  } else {
    localStorage.removeItem("syncnote.userId");
  }
});

function saveEndpoints() {
  localStorage.setItem("syncnote.authBase", authBase.value.trim());
  localStorage.setItem("syncnote.syncBase", syncBase.value.trim());
  message.value = "API 地址已保存";
}

function setMessage(text) {
  message.value = text;
}

function decodeUserIdFromJwt(jwt) {
  try {
    const parts = jwt.split(".");
    if (parts.length !== 3) return "";
    const payload = JSON.parse(atob(parts[1]));
    return payload.userId || payload.user_id || "";
  } catch {
    return "";
  }
}

async function request(url, options = {}) {
  const res = await fetch(url, options);
  const text = await res.text();
  let data = {};

  try {
    data = text ? JSON.parse(text) : {};
  } catch {
    data = { raw: text };
  }

  if (!res.ok) {
    throw new Error(
      data.message || data.error || text || `request failed: ${res.status}`,
    );
  }

  return data;
}

function authHeaders() {
  return token.value
    ? {
        Authorization: `Bearer ${token.value}`,
      }
    : {};
}

function upsertNote(note) {
  if (!note?.note_id) return;
  const idx = notes.value.findIndex((n) => n.note_id === note.note_id);
  if (idx >= 0) {
    notes.value[idx] = note;
  } else {
    notes.value.unshift(note);
  }
}

async function register() {
  busy.value = true;
  try {
    const data = await request(`${authBase.value}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: registerForm.value.email,
        password: registerForm.value.password,
        captcha: "",
      }),
    });
    setMessage(`注册成功，userId=${data.userId || "-"}`);
    authMode.value = "login";
  } catch (err) {
    setMessage(`注册失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function login() {
  busy.value = true;
  try {
    const data = await request(`${authBase.value}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: loginForm.value.email,
        password: loginForm.value.password,
        captcha: "",
      }),
    });

    token.value = data.token || "";
    userId.value = decodeUserIdFromJwt(token.value);
    setMessage("登录成功");
  } catch (err) {
    setMessage(`登录失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

function logout() {
  token.value = "";
  userId.value = "";
  setMessage("已退出登录");
}

async function createNote() {
  busy.value = true;
  try {
    const data = await request(`${syncBase.value}/api/v1/notes`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders(),
      },
      body: JSON.stringify(createNoteForm.value),
    });

    upsertNote(data.note);
    createNoteForm.value = { title: "", content: "" };
    setMessage("笔记创建成功");
  } catch (err) {
    setMessage(`创建笔记失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function fetchNoteById() {
  if (!loadNoteId.value.trim()) {
    setMessage("请输入 note_id");
    return;
  }
  busy.value = true;
  try {
    const data = await request(
      `${syncBase.value}/api/v1/notes/${loadNoteId.value.trim()}`,
      {
        headers: {
          ...authHeaders(),
        },
      },
    );
    upsertNote(data.note);
    setMessage("笔记加载成功");
  } catch (err) {
    setMessage(`加载笔记失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function createTeam() {
  busy.value = true;
  try {
    const data = await request(`${syncBase.value}/api/v1/teams`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders(),
      },
      body: JSON.stringify(createTeamForm.value),
    });

    if (data.team) {
      teams.value.unshift(data.team);
    }
    createTeamForm.value = { name: "", description: "" };
    setMessage("团队创建成功");
  } catch (err) {
    setMessage(`创建团队失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function joinTeam() {
  if (!joinTeamId.value.trim()) {
    setMessage("请输入 team_id");
    return;
  }
  busy.value = true;
  try {
    await request(
      `${syncBase.value}/api/v1/teams/${joinTeamId.value.trim()}/join`,
      {
        method: "POST",
        headers: {
          ...authHeaders(),
        },
      },
    );
    setMessage("加入团队成功");
  } catch (err) {
    setMessage(`加入团队失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function grantPermission() {
  if (!grantForm.value.noteId || !grantForm.value.targetUserId) {
    setMessage("请填写 note_id 与 target_user_id");
    return;
  }
  busy.value = true;
  try {
    await request(
      `${syncBase.value}/api/v1/notes/${grantForm.value.noteId}/permissions`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...authHeaders(),
        },
        body: JSON.stringify({
          target_user_id: grantForm.value.targetUserId,
          level: Number(grantForm.value.level),
        }),
      },
    );
    setMessage("授权成功");
  } catch (err) {
    setMessage(`授权失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

function prevPage() {
  page.value = Math.max(1, page.value - 1);
}

function nextPage() {
  page.value = Math.min(totalPages.value, page.value + 1);
}

onMounted(() => {
  if (token.value && !userId.value) {
    userId.value = decodeUserIdFromJwt(token.value);
  }
});
</script>

<template>
  <main class="page-shell">
    <header class="hero">
      <h1>SyncNote Vue Console</h1>
      <p>登录注册、创建笔记、笔记分页、创建团队与加入团队的一体化调试面板。</p>
    </header>

    <section class="panel grid-2">
      <div class="card">
        <h2>API 配置</h2>
        <label>
          Auth Base
          <input v-model="authBase" placeholder="http://127.0.0.1:8000" />
        </label>
        <label>
          SyncNote Base
          <input v-model="syncBase" placeholder="http://127.0.0.1:8001" />
        </label>
        <button @click="saveEndpoints">保存地址</button>
      </div>

      <div class="card">
        <h2>会话信息</h2>
        <p><strong>userId:</strong> {{ userId || "-" }}</p>
        <p class="token"><strong>token:</strong> {{ token || "-" }}</p>
        <button class="ghost" @click="logout">退出登录</button>
      </div>
    </section>

    <section class="panel">
      <div class="card">
        <h2>登录 / 注册</h2>
        <div class="auth-switch">
          <button
            :class="{ active: authMode === 'login' }"
            @click="authMode = 'login'"
          >
            登录
          </button>
          <button
            :class="{ active: authMode === 'register' }"
            @click="authMode = 'register'"
          >
            注册
          </button>
        </div>

        <div v-if="authMode === 'register'" class="form-grid">
          <input v-model="registerForm.email" type="email" placeholder="邮箱" />
          <input
            v-model="registerForm.password"
            type="password"
            placeholder="密码"
          />
          <button :disabled="busy" @click="register">注册</button>
        </div>

        <div v-else class="form-grid">
          <input v-model="loginForm.email" type="email" placeholder="邮箱" />
          <input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
          />
          <button :disabled="busy" @click="login">登录</button>
        </div>
      </div>
    </section>

    <section class="panel grid-2">
      <div class="card">
        <h2>笔记操作</h2>
        <div class="form-grid">
          <input v-model="createNoteForm.title" placeholder="笔记标题" />
          <textarea
            v-model="createNoteForm.content"
            rows="4"
            placeholder="笔记内容"
          />
          <button :disabled="busy || !token" @click="createNote">
            创建笔记
          </button>
        </div>
        <div class="inline-tools">
          <input v-model="loadNoteId" placeholder="输入 note_id 查询" />
          <button :disabled="busy || !token" @click="fetchNoteById">
            查询笔记
          </button>
        </div>
      </div>

      <div class="card">
        <h2>团队操作</h2>
        <div class="form-grid">
          <input v-model="createTeamForm.name" placeholder="团队名称" />
          <textarea
            v-model="createTeamForm.description"
            rows="3"
            placeholder="团队描述"
          />
          <button :disabled="busy || !token" @click="createTeam">
            创建团队
          </button>
        </div>
        <div class="inline-tools">
          <input v-model="joinTeamId" placeholder="输入 team_id 加入" />
          <button :disabled="busy || !token" @click="joinTeam">加入团队</button>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="card">
        <h2>授权</h2>
        <div class="grant-grid">
          <input v-model="grantForm.noteId" placeholder="note_id" />
          <input
            v-model="grantForm.targetUserId"
            placeholder="target_user_id"
          />
          <select v-model="grantForm.level">
            <option :value="1">read (1)</option>
            <option :value="2">write (2)</option>
            <option :value="3">admin (3)</option>
          </select>
          <button :disabled="busy || !token" @click="grantPermission">
            授予权限
          </button>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="card">
        <div class="list-head">
          <h2>笔记列表（前端分页）</h2>
          <div class="pager-tools">
            <label>
              每页
              <select v-model.number="pageSize">
                <option :value="5">5</option>
                <option :value="10">10</option>
                <option :value="20">20</option>
              </select>
            </label>
            <button @click="prevPage" :disabled="page <= 1">上一页</button>
            <span>{{ page }} / {{ totalPages }}</span>
            <button @click="nextPage" :disabled="page >= totalPages">
              下一页
            </button>
          </div>
        </div>

        <div class="note-grid" v-if="pagedNotes.length">
          <article
            v-for="item in pagedNotes"
            :key="item.note_id"
            class="note-card"
          >
            <h3>{{ item.title }}</h3>
            <p class="note-meta">note_id: {{ item.note_id }}</p>
            <p class="note-meta">
              owner: {{ item.owner_id }} | version: {{ item.version }}
            </p>
            <p>{{ item.content }}</p>
          </article>
        </div>
        <p v-else class="empty">暂无笔记，先创建或输入 note_id 查询。</p>
      </div>
    </section>

    <section class="panel" v-if="teams.length">
      <div class="card">
        <h2>最近创建的团队</h2>
        <div class="team-list">
          <article v-for="team in teams" :key="team.team_id" class="team-item">
            <h3>{{ team.name }}</h3>
            <p>{{ team.description }}</p>
            <p class="note-meta">team_id: {{ team.team_id }}</p>
          </article>
        </div>
      </div>
    </section>

    <footer class="status" :class="{ error: message.includes('失败') }">
      {{ message || "就绪" }}
    </footer>
  </main>
</template>
