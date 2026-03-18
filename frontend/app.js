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
  log("SyncNote Demo ready.");
}

function bindEvents() {
  els.saveHostBtn.addEventListener("click", () => {
    localStorage.setItem(storageKeys.authHost, els.authHost.value.trim());
    localStorage.setItem(storageKeys.apiHost, els.apiHost.value.trim());
    log("Host 配置已保存。");
  });

  els.registerBtn.addEventListener("click", register);
  els.loginBtn.addEventListener("click", login);
  els.createBtn.addEventListener("click", createNote);
  els.getBtn.addEventListener("click", getNote);
  els.saveBtn.addEventListener("click", saveNote);
  els.listBtn.addEventListener("click", listNotes);
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
}

async function register() {
  const payload = {
    username: els.username.value.trim(),
    email: els.email.value.trim(),
    password: els.password.value,
    captcha: "233",
  };
  const data = await req(`${els.authHost.value.trim()}/auth/register`, {
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
  const payload = {
    username: els.username.value.trim(),
    password: els.password.value,
  };
  const data = await req(`${els.authHost.value.trim()}/auth/login`, {
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
  const payload = {
    title: els.noteTitle.value.trim(),
    content: els.noteContent.value,
  };
  const data = await req(`${els.apiHost.value.trim()}/api/note/create`, {
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
  const noteId = els.noteId.value.trim();
  if (!noteId) {
    alert("请先输入或选择 Note ID");
    return;
  }

  const data = await req(`${els.apiHost.value.trim()}/api/note/${encodeURIComponent(noteId)}`, {
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

  const data = await req(`${els.apiHost.value.trim()}/api/note/save`, {
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
  const data = await req(`${els.apiHost.value.trim()}/api/user/notes`, {
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
