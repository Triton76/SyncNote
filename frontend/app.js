const els = {
  authHost: document.getElementById("authHost"),
  apiHost: document.getElementById("apiHost"),
  username: document.getElementById("username"),
  email: document.getElementById("email"),
  password: document.getElementById("password"),
  token: document.getElementById("token"),
  noteTitle: document.getElementById("noteTitle"),
  noteId: document.getElementById("noteId"),
  noteContent: document.getElementById("noteContent"),
  expectedVersion: document.getElementById("expectedVersion"),
  currentVersion: document.getElementById("currentVersion"),
  notesList: document.getElementById("notesList"),
  logOutput: document.getElementById("logOutput"),
  saveHostBtn: document.getElementById("saveHostBtn"),
  registerBtn: document.getElementById("registerBtn"),
  loginBtn: document.getElementById("loginBtn"),
  createBtn: document.getElementById("createBtn"),
  getBtn: document.getElementById("getBtn"),
  saveBtn: document.getElementById("saveBtn"),
  listBtn: document.getElementById("listBtn"),
};

const storageKeys = {
  authHost: "syncnote.authHost",
  apiHost: "syncnote.apiHost",
  token: "syncnote.token",
};

function boot() {
  els.authHost.value = localStorage.getItem(storageKeys.authHost) || els.authHost.value;
  els.apiHost.value = localStorage.getItem(storageKeys.apiHost) || els.apiHost.value;
  els.token.value = localStorage.getItem(storageKeys.token) || "";
  bindEvents();
  log(`Page origin: ${window.location.origin}`);
  log(`Auth Host: ${els.authHost.value.trim()} | API Host: ${els.apiHost.value.trim()}`);
  log("SyncNote Demo ready.");
}

function bindEvents() {
  els.saveHostBtn.addEventListener("click", () => {
    const authHost = normalizeHost(els.authHost.value);
    const apiHost = normalizeHost(els.apiHost.value);
    els.authHost.value = authHost;
    els.apiHost.value = apiHost;
    localStorage.setItem(storageKeys.authHost, authHost);
    localStorage.setItem(storageKeys.apiHost, apiHost);
    log(`Host 配置已保存。Auth=${authHost}, API=${apiHost}`);
  });

  els.registerBtn.addEventListener("click", () => runAction("注册", register));
  els.loginBtn.addEventListener("click", () => runAction("登录", login));
  els.createBtn.addEventListener("click", () => runAction("创建笔记", createNote));
  els.getBtn.addEventListener("click", () => runAction("读取笔记", getNote));
  els.saveBtn.addEventListener("click", () => runAction("保存笔记", saveNote));
  els.listBtn.addEventListener("click", () => runAction("拉取笔记列表", listNotes));
}

function jsonHeaders(withAuth = false) {
  const headers = {
    "Content-Type": "application/json",
  };
  if (withAuth && els.token.value.trim()) {
    headers.Authorization = `Bearer ${els.token.value.trim()}`;
  }
  return headers;
}

async function req(url, options = {}) {
  try {
    const res = await fetch(url, options);
    const text = await res.text();
    let data = text;
    try {
      data = JSON.parse(text);
    } catch (_e) {
      // keep plain text when response is not json
    }
    log(`${options.method || "GET"} ${url}\nstatus=${res.status}\n${pretty(data)}\n`);
    if (!res.ok) {
      throw new Error(`请求失败: ${res.status}`);
    }
    return data;
  } catch (err) {
    const msg = err && err.message ? err.message : String(err);
    const hint =
      msg.includes("Failed to fetch")
        ? "网络失败或被 CORS 拦截。请检查后端是否启动、端口是否正确、以及 OPTIONS 预检是否返回 204。"
        : "请求异常，请检查日志。";
    log(`请求异常: ${msg}\n诊断建议: ${hint}\n`);
    throw err;
  }
}

function normalizeHost(val) {
  return String(val || "").trim().replace(/\/+$/, "");
}

async function runAction(title, fn) {
  try {
    await fn();
  } catch (err) {
    const msg = err && err.message ? err.message : String(err);
    alert(`${title}失败: ${msg}`);
  }
}

async function register() {
  const authHost = normalizeHost(els.authHost.value);
  els.authHost.value = authHost;
  const payload = {
    username: els.username.value.trim(),
    email: els.email.value.trim(),
    password: els.password.value,
    captcha: "233",
  };
  const data = await req(`${authHost}/auth/register`, {
    method: "POST",
    headers: jsonHeaders(),
    body: JSON.stringify(payload),
  });
  const token = data.token || data.Token || "";
  if (token) {
    els.token.value = token;
    localStorage.setItem(storageKeys.token, token);
  }
}

async function login() {
  const authHost = normalizeHost(els.authHost.value);
  els.authHost.value = authHost;

  const loginId = (els.email.value || "").trim() || (els.username.value || "").trim();
  if (!loginId) {
    alert("登录需要 loginId。请优先填写邮箱（当前后端按邮箱登录）。");
    return;
  }

  const payload = {
    loginId,
    password: els.password.value,
  };

  log(`登录使用 loginId=${loginId}`);
  const data = await req(`${authHost}/auth/login`, {
    method: "POST",
    headers: jsonHeaders(),
    body: JSON.stringify(payload),
  });
  const token = data.token || data.Token || "";
  if (token) {
    els.token.value = token;
    localStorage.setItem(storageKeys.token, token);
  }
}

async function createNote() {
  const apiHost = normalizeHost(els.apiHost.value);
  els.apiHost.value = apiHost;
  const payload = {
    title: els.noteTitle.value.trim(),
    content: els.noteContent.value,
  };
  const data = await req(`${apiHost}/api/note/create`, {
    method: "POST",
    headers: jsonHeaders(true),
    body: JSON.stringify(payload),
  });

  const noteId = data.noteId || data.NoteId || "";
  const version = data.version || data.Version || 0;

  if (noteId) {
    els.noteId.value = noteId;
  }
  if (version) {
    els.currentVersion.value = version;
    els.expectedVersion.value = version;
  }
}

async function getNote() {
  const apiHost = normalizeHost(els.apiHost.value);
  els.apiHost.value = apiHost;
  const noteId = els.noteId.value.trim();
  if (!noteId) {
    alert("请先输入或选择 Note ID");
    return;
  }

  const data = await req(`${apiHost}/api/note/${encodeURIComponent(noteId)}`, {
    headers: {
      Authorization: `Bearer ${els.token.value.trim()}`,
    },
  });

  els.noteTitle.value = data.title || data.Title || "";
  els.noteContent.value = data.content || data.Content || "";
  const version = data.version || data.Version || 0;
  els.currentVersion.value = version;
  els.expectedVersion.value = version;
}

async function saveNote() {
  const apiHost = normalizeHost(els.apiHost.value);
  els.apiHost.value = apiHost;
  const noteId = els.noteId.value.trim();
  if (!noteId) {
    alert("请先输入或选择 Note ID");
    return;
  }

  const payload = {
    noteId,
    content: els.noteContent.value,
    expectedVersion: Number(els.expectedVersion.value || 0),
  };

  const data = await req(`${apiHost}/api/note/save`, {
    method: "POST",
    headers: jsonHeaders(true),
    body: JSON.stringify(payload),
  });

  const note = data.note || data.Note;
  if (note) {
    const version = note.version || note.Version || 0;
    if (version) {
      els.currentVersion.value = version;
      els.expectedVersion.value = version;
    }
  }
}

async function listNotes() {
  const apiHost = normalizeHost(els.apiHost.value);
  els.apiHost.value = apiHost;
  const data = await req(`${apiHost}/api/user/notes`, {
    headers: {
      Authorization: `Bearer ${els.token.value.trim()}`,
    },
  });

  const notes = data.notes || data.Notes || [];
  els.notesList.innerHTML = "";

  notes.forEach((item) => {
    const noteId = item.noteId || item.NoteId || "";
    const title = item.title || item.Title || "(untitled)";
    const version = item.version || item.Version || 0;
    const li = document.createElement("li");
    li.textContent = `${title} | id=${noteId} | v${version}`;
    li.addEventListener("click", async () => {
      els.noteId.value = noteId;
      await getNote();
    });
    els.notesList.appendChild(li);
  });
}

function pretty(val) {
  if (typeof val === "string") {
    return val;
  }
  try {
    return JSON.stringify(val, null, 2);
  } catch (_e) {
    return String(val);
  }
}

function log(msg) {
  const ts = new Date().toLocaleTimeString();
  els.logOutput.textContent = `[${ts}] ${msg}\n` + els.logOutput.textContent;
}

boot();
