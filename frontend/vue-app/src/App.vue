<script setup>
import { computed, onMounted, ref, watch } from "vue";
import ApiConfigCard from "./components/ApiConfigCard.vue";
import AuthPanel from "./components/AuthPanel.vue";
import DebugPanel from "./components/DebugPanel.vue";
import GrantPermissionPanel from "./components/GrantPermissionPanel.vue";
import NoteActionsPanel from "./components/NoteActionsPanel.vue";
import NoteDetailPanel from "./components/NoteDetailPanel.vue";
import NoteSidebar from "./components/NoteSidebar.vue";
import TeamPanel from "./components/TeamPanel.vue";

const authBase = ref(
  localStorage.getItem("syncnote.authBase") ||
    import.meta.env.VITE_AUTH_BASE ||
    "",
);
const syncBase = ref(
  localStorage.getItem("syncnote.syncBase") ||
    import.meta.env.VITE_SYNC_BASE ||
    "",
);

const token = ref(localStorage.getItem("syncnote.token") || "");
const userId = ref(localStorage.getItem("syncnote.userId") || "");

const notes = ref([]);
const teams = ref([]);
const selectedNoteId = ref("");
const page = ref(1);
const pageSize = ref(5);

const busy = ref(false);
const message = ref("");

const isAuthed = computed(() => !!token.value);
const selectedNote = computed(
  () =>
    notes.value.find((item) => item.note_id === selectedNoteId.value) || null,
);

const debugInfo = computed(() => ({
  now: new Date().toISOString(),
  authBase: authBase.value || "(proxy)",
  syncBase: syncBase.value || "(proxy)",
  tokenPreview: token.value ? `${token.value.slice(0, 12)}...` : "",
  userId: userId.value,
  notesCount: notes.value.length,
  teamsCount: teams.value.length,
  selectedNoteId: selectedNoteId.value,
  page: page.value,
  pageSize: pageSize.value,
  busy: busy.value,
}));

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

watch(pageSize, () => {
  page.value = 1;
});

function trimTrailingSlash(value) {
  if (!value || value === "/") return value;
  return value.endsWith("/") ? value.slice(0, -1) : value;
}

function normalizeEndpointBase(raw, service) {
  const value = trimTrailingSlash((raw || "").trim());
  if (!value || value === "/") return "";

  try {
    const url = new URL(value);
    const host = url.hostname;
    const port = url.port || (url.protocol === "https:" ? "443" : "80");
    const isLocalHost = host === "127.0.0.1" || host === "localhost";

    if (isLocalHost) {
      if (service === "auth" && port === "8000") return "";
      if (service === "sync" && port === "8001") return "";
    }
  } catch {
    // Allow custom paths or non-url input.
  }

  return value;
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

function authHeaders() {
  return token.value
    ? {
        Authorization: `Bearer ${token.value}`,
      }
    : {};
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

function upsertNote(note) {
  if (!note?.note_id) return;
  const idx = notes.value.findIndex((item) => item.note_id === note.note_id);
  if (idx >= 0) {
    notes.value[idx] = note;
  } else {
    notes.value.unshift(note);
  }
  selectedNoteId.value = note.note_id;
  page.value = 1;
}

function upsertTeam(team) {
  if (!team?.team_id) return;
  const idx = teams.value.findIndex((item) => item.team_id === team.team_id);
  if (idx >= 0) {
    teams.value[idx] = team;
  } else {
    teams.value.unshift(team);
  }
}

function saveEndpoints() {
  const normalizedAuth = normalizeEndpointBase(authBase.value, "auth");
  const normalizedSync = normalizeEndpointBase(syncBase.value, "sync");

  authBase.value = normalizedAuth;
  syncBase.value = normalizedSync;

  localStorage.setItem("syncnote.authBase", normalizedAuth);
  localStorage.setItem("syncnote.syncBase", normalizedSync);

  if (!normalizedAuth && !normalizedSync) {
    setMessage("API 地址已保存（同源代理模式）");
    return;
  }

  setMessage("API 地址已保存");
}

async function login(payload) {
  busy.value = true;
  try {
    const data = await request(`${authBase.value}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: payload.email,
        password: payload.password,
        captcha: "",
      }),
    });

    token.value = data.token || "";
    userId.value = decodeUserIdFromJwt(token.value);
    setMessage("登录成功，已进入 Dashboard");
  } catch (err) {
    setMessage(`登录失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function register(payload) {
  busy.value = true;
  try {
    const data = await request(`${authBase.value}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: payload.email,
        password: payload.password,
        captcha: "",
      }),
    });
    setMessage(`注册成功，userId=${data.userId || "-"}`);
  } catch (err) {
    setMessage(`注册失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

function logout() {
  token.value = "";
  userId.value = "";
  selectedNoteId.value = "";
  setMessage("已退出登录");
}

async function createNote(payload, done) {
  busy.value = true;
  try {
    const data = await request(`${syncBase.value}/api/v1/notes`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders(),
      },
      body: JSON.stringify(payload),
    });

    upsertNote(data.note);
    if (typeof done === "function") done();
    setMessage("笔记创建成功");
  } catch (err) {
    setMessage(`创建笔记失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function fetchNoteById(noteId) {
  if (!noteId) {
    setMessage("请输入 note_id");
    return;
  }
  busy.value = true;
  try {
    const data = await request(`${syncBase.value}/api/v1/notes/${noteId}`, {
      headers: {
        ...authHeaders(),
      },
    });
    upsertNote(data.note);
    setMessage("笔记加载成功");
  } catch (err) {
    setMessage(`加载笔记失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function updateNote(payload) {
  busy.value = true;
  try {
    const data = await request(
      `${syncBase.value}/api/v1/notes/${payload.noteId}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          ...authHeaders(),
        },
        body: JSON.stringify({
          title: payload.title,
          content: payload.content,
          version: payload.version,
        }),
      },
    );
    upsertNote(data.note);
    setMessage("笔记已保存");
  } catch (err) {
    setMessage(`保存笔记失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function deleteNote(noteId) {
  busy.value = true;
  try {
    await request(`${syncBase.value}/api/v1/notes/${noteId}`, {
      method: "DELETE",
      headers: {
        ...authHeaders(),
      },
    });
    notes.value = notes.value.filter((item) => item.note_id !== noteId);
    if (selectedNoteId.value === noteId) {
      selectedNoteId.value = notes.value[0]?.note_id || "";
    }
    setMessage("笔记已删除");
  } catch (err) {
    setMessage(`删除笔记失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function createTeam(payload, done) {
  busy.value = true;
  try {
    const data = await request(`${syncBase.value}/api/v1/teams`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders(),
      },
      body: JSON.stringify(payload),
    });

    if (data.team) {
      teams.value.unshift(data.team);
    }
    if (typeof done === "function") done();
    setMessage("团队创建成功");
  } catch (err) {
    setMessage(`创建团队失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function joinTeam(teamId) {
  if (!teamId) {
    setMessage("请输入 team_id");
    return;
  }
  busy.value = true;
  try {
    await request(`${syncBase.value}/api/v1/teams/${teamId}/join`, {
      method: "POST",
      headers: {
        ...authHeaders(),
      },
    });

    const data = await request(`${syncBase.value}/api/v1/teams/${teamId}`, {
      headers: {
        ...authHeaders(),
      },
    });
    if (data.team) {
      upsertTeam(data.team);
    }

    setMessage("加入团队成功");
  } catch (err) {
    setMessage(`加入团队失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

async function grantPermission(payload, done) {
  const noteId = (payload.noteId || "").trim();
  const targetUserId = (payload.targetUserId || "").trim();
  const level = Number(payload.level || 1);

  if (!noteId || !targetUserId) {
    setMessage("请填写 note_id 和 target_user_id");
    return;
  }

  busy.value = true;
  try {
    await request(`${syncBase.value}/api/v1/notes/${noteId}/permissions`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders(),
      },
      body: JSON.stringify({
        target_user_id: targetUserId,
        level,
      }),
    });
    if (typeof done === "function") done();
    setMessage("授权成功");
  } catch (err) {
    setMessage(`授权失败: ${err.message}`);
  } finally {
    busy.value = false;
  }
}

onMounted(() => {
  if (token.value && !userId.value) {
    userId.value = decodeUserIdFromJwt(token.value);
  }

  const migratedAuth = normalizeEndpointBase(authBase.value, "auth");
  const migratedSync = normalizeEndpointBase(syncBase.value, "sync");
  const changed =
    migratedAuth !== authBase.value || migratedSync !== syncBase.value;

  if (changed) {
    authBase.value = migratedAuth;
    syncBase.value = migratedSync;
    localStorage.setItem("syncnote.authBase", migratedAuth);
    localStorage.setItem("syncnote.syncBase", migratedSync);
    setMessage("检测到旧地址，已自动迁移到同源代理模式");
  }
});
</script>

<template>
  <main class="app-shell">
    <header class="hero">
      <h1>SyncNote 控制台</h1>
      <p>
        登录后进入
        Dashboard，使用点击操作团队与笔记，支持笔记详情编辑与调试信息展示。
      </p>
    </header>

    <section v-if="!isAuthed" class="auth-layout">
      <ApiConfigCard
        :auth-base="authBase"
        :sync-base="syncBase"
        :busy="busy"
        @update:auth-base="authBase = $event"
        @update:sync-base="syncBase = $event"
        @save="saveEndpoints"
      />
      <AuthPanel :busy="busy" @login="login" @register="register" />
    </section>

    <section v-else class="dashboard-layout">
      <div class="top-row">
        <ApiConfigCard
          :auth-base="authBase"
          :sync-base="syncBase"
          :busy="busy"
          @update:auth-base="authBase = $event"
          @update:sync-base="syncBase = $event"
          @save="saveEndpoints"
        />

        <section class="card">
          <h3>会话信息</h3>
          <p class="muted">userId: {{ userId || "-" }}</p>
          <p class="muted">
            token: {{ token ? `${token.slice(0, 16)}...` : "-" }}
          </p>
          <button class="ghost" :disabled="busy" @click="logout">
            退出登录
          </button>
        </section>
      </div>

      <div class="top-row">
        <TeamPanel
          :busy="busy"
          :authed="isAuthed"
          :teams="teams"
          @create-team="createTeam"
          @join-team="joinTeam"
        />
        <NoteActionsPanel
          :busy="busy"
          :authed="isAuthed"
          @create-note="createNote"
          @fetch-note="fetchNoteById"
        />
      </div>

      <div class="single-row">
        <GrantPermissionPanel
          :busy="busy"
          :authed="isAuthed"
          @grant-permission="grantPermission"
        />
      </div>

      <section class="notes-layout">
        <NoteSidebar
          :notes="notes"
          :page="page"
          :page-size="pageSize"
          :selected-note-id="selectedNoteId"
          @update:page="page = $event"
          @update:page-size="pageSize = $event"
          @select="fetchNoteById"
        />

        <NoteDetailPanel
          :note="selectedNote"
          :busy="busy"
          @save="updateNote"
          @remove="deleteNote"
          @refresh="fetchNoteById"
        />
      </section>

      <DebugPanel :debug-info="debugInfo" />
    </section>

    <footer class="status" :class="{ error: message.includes('失败') }">
      {{ message || "就绪" }}
    </footer>
  </main>
</template>
