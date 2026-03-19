const storageKeys = {
  authHost: "syncnote.authHost",
  apiHost: "syncnote.apiHost",
  token: "syncnote.token",
};

function normalizeHost(val) {
  return String(val || "").trim().replace(/\/+$/, "");
}

function readTokenLikeField(payload) {
  return payload?.token || payload?.Token || "";
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

function hostAliasCandidates(url) {
  try {
    const u = new URL(url);
    if (u.hostname === "127.0.0.1") {
      u.hostname = "localhost";
      return [u.toString()];
    }
    if (u.hostname === "localhost") {
      u.hostname = "127.0.0.1";
      return [u.toString()];
    }
  } catch (_e) {
    return [];
  }
  return [];
}

const { createApp } = Vue;

createApp({
  data() {
    return {
      loading: false,
      authHost: localStorage.getItem(storageKeys.authHost) || "http://127.0.0.1:8889",
      apiHost: localStorage.getItem(storageKeys.apiHost) || "http://127.0.0.1:8888",
      username: "",
      email: "",
      password: "",
      token: localStorage.getItem(storageKeys.token) || "",
      noteTitle: "",
      noteId: "",
      noteContent: "",
      expectedVersion: 1,
      currentVersion: 0,
      notes: [],
      targetUserId: "",
      targetUserEmail: "",
      targetTeamId: "",
      permissionRole: "viewer",
      permissions: [],
      myTeams: [],
      eventStartSeq: 0,
      eventLimit: 50,
      events: [],
      logs: ["SyncNote Workbench ready."],
      connState: "unknown",
      flash: {
        type: "",
        text: "",
      },
    };
  },
  computed: {
    connText() {
      switch (this.connState) {
        case "ok":
          return "后端在线";
        case "fail":
          return "后端不可用";
        default:
          return "未检测";
      }
    },
  },
  methods: {
    setFlash(type, text) {
      this.flash.type = type;
      this.flash.text = text;
      setTimeout(() => {
        if (this.flash.text === text) {
          this.flash.type = "";
          this.flash.text = "";
        }
      }, 2800);
    },
    log(msg) {
      const ts = new Date().toLocaleTimeString();
      this.logs.unshift(`[${ts}] ${msg}`);
      if (this.logs.length > 80) {
        this.logs = this.logs.slice(0, 80);
      }
    },
    normalizeHosts() {
      this.authHost = normalizeHost(this.authHost);
      this.apiHost = normalizeHost(this.apiHost);
    },
    saveHosts() {
      this.normalizeHosts();
      localStorage.setItem(storageKeys.authHost, this.authHost);
      localStorage.setItem(storageKeys.apiHost, this.apiHost);
      this.log(`Host saved: auth=${this.authHost}, api=${this.apiHost}`);
      this.setFlash("success", "地址已保存");
    },
    persistToken() {
      if (this.token) {
        localStorage.setItem(storageKeys.token, this.token);
      } else {
        localStorage.removeItem(storageKeys.token);
      }
    },
    clearLogs() {
      this.logs = [];
      this.log("日志已清空");
    },
    logout() {
      this.token = "";
      this.persistToken();
      this.setFlash("success", "已退出，Token 已清空");
    },
    async copyToken() {
      if (!this.token) {
        return;
      }
      if (!navigator.clipboard) {
        this.setFlash("error", "当前环境不支持剪贴板复制");
        return;
      }
      await navigator.clipboard.writeText(this.token);
      this.setFlash("success", "Token 已复制");
    },
    jsonHeaders(withAuth = false) {
      const headers = { "Content-Type": "application/json" };
      if (withAuth && this.token) {
        headers.Authorization = `Bearer ${this.token}`;
      }
      return headers;
    },
    async request(url, options = {}) {
      try {
        const res = await fetch(url, options);
        const text = await res.text();
        let data = text;
        try {
          data = JSON.parse(text);
        } catch (_e) {
          // no-op
        }

        this.log(`${options.method || "GET"} ${url}\nstatus=${res.status}\n${pretty(data)}`);

        if (!res.ok) {
          const msg = typeof data === "string" ? data : (data?.message || data?.Message || `HTTP ${res.status}`);
          throw new Error(msg);
        }
        return data;
      } catch (err) {
        const message = err?.message || String(err);
        if (message.includes("Failed to fetch")) {
          const aliases = hostAliasCandidates(url);
          for (const nextUrl of aliases) {
            try {
              this.log(`Primary URL failed, retry alias URL: ${nextUrl}`);
              const retryRes = await fetch(nextUrl, options);
              const retryText = await retryRes.text();
              let retryData = retryText;
              try {
                retryData = JSON.parse(retryText);
              } catch (_e) {
                // no-op
              }
              this.log(`${options.method || "GET"} ${nextUrl}\nstatus=${retryRes.status}\n${pretty(retryData)}`);
              if (!retryRes.ok) {
                const retryMsg = typeof retryData === "string"
                  ? retryData
                  : (retryData?.message || retryData?.Message || `HTTP ${retryRes.status}`);
                throw new Error(retryMsg);
              }
              return retryData;
            } catch (_retryErr) {
              // continue trying candidates
            }
          }
          throw new Error(`网络失败（非业务报错）。请检查前端地址栏与 Host 配置是否可达。当前请求: ${url}`);
        }
        throw err;
      }
    },
    async wrapAction(title, fn) {
      this.loading = true;
      try {
        await fn();
        this.setFlash("success", `${title}成功`);
      } catch (err) {
        const msg = err?.message || String(err);
        this.log(`${title}失败: ${msg}`);
        this.setFlash("error", `${title}失败: ${msg}`);
        alert(`${title}失败: ${msg}`);
      } finally {
        this.loading = false;
      }
    },
    async checkConnection() {
      await this.wrapAction("连接检测", async () => {
        this.normalizeHosts();
        const authResp = await fetch(`${this.authHost}/auth/register`, {
          method: "OPTIONS",
          headers: {
            Origin: window.location.origin,
            "Access-Control-Request-Method": "POST",
            "Access-Control-Request-Headers": "content-type",
          },
        });

        const apiResp = await fetch(`${this.apiHost}/api/note/create`, {
          method: "OPTIONS",
          headers: {
            Origin: window.location.origin,
            "Access-Control-Request-Method": "POST",
            "Access-Control-Request-Headers": "authorization,content-type",
          },
        });

        if (authResp.status >= 400 || apiResp.status >= 400) {
          this.connState = "fail";
          throw new Error(`auth=${authResp.status}, api=${apiResp.status}`);
        }
        this.connState = "ok";
      });
    },
    async register() {
      await this.wrapAction("注册", async () => {
        this.normalizeHosts();
        const data = await this.request(`${this.authHost}/auth/register`, {
          method: "POST",
          headers: this.jsonHeaders(false),
          body: JSON.stringify({
            username: this.username.trim(),
            email: this.email.trim(),
            password: this.password,
            captcha: "233",
          }),
        });
        const token = readTokenLikeField(data);
        if (token) {
          this.token = token;
          this.persistToken();
        }
        this.connState = "ok";
        await this.listMyTeams(true);
      });
    },
    async login() {
      await this.wrapAction("登录", async () => {
        this.normalizeHosts();
        const loginId = this.email.trim() || this.username.trim();
        if (!loginId) {
          throw new Error("请填写邮箱或用户名作为 loginId");
        }
        const data = await this.request(`${this.authHost}/auth/login`, {
          method: "POST",
          headers: this.jsonHeaders(false),
          body: JSON.stringify({ loginId, password: this.password }),
        });
        const token = readTokenLikeField(data);
        if (token) {
          this.token = token;
          this.persistToken();
        }
        this.connState = "ok";
        await this.listMyTeams(true);
      });
    },
    ensureNoteId() {
      if (!this.noteId.trim()) {
        throw new Error("请先创建或选择 Note ID");
      }
      return this.noteId.trim();
    },
    async createNote() {
      await this.wrapAction("创建笔记", async () => {
        this.normalizeHosts();
        const data = await this.request(`${this.apiHost}/api/note/create`, {
          method: "POST",
          headers: this.jsonHeaders(true),
          body: JSON.stringify({ title: this.noteTitle.trim(), content: this.noteContent }),
        });
        this.noteId = data.noteId || data.NoteId || "";
        this.currentVersion = Number(data.version || data.Version || 1);
        this.expectedVersion = this.currentVersion;
      });
    },
    async getNote() {
      await this.wrapAction("读取笔记", async () => {
        this.normalizeHosts();
        const noteId = this.ensureNoteId();
        const data = await this.request(`${this.apiHost}/api/note/${encodeURIComponent(noteId)}`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        this.noteTitle = data.title || data.Title || "";
        this.noteContent = data.content || data.Content || "";
        this.currentVersion = Number(data.version || data.Version || 0);
        this.expectedVersion = this.currentVersion || this.expectedVersion;
      });
    },
    async saveNote() {
      await this.wrapAction("保存笔记", async () => {
        this.normalizeHosts();
        const noteId = this.ensureNoteId();
        const data = await this.request(`${this.apiHost}/api/note/save`, {
          method: "POST",
          headers: this.jsonHeaders(true),
          body: JSON.stringify({
            noteId,
            title: this.noteTitle.trim(),
            content: this.noteContent,
            expectedVersion: Number(this.expectedVersion || 0),
          }),
        });

        const note = data.note || data.Note;
        if (note) {
          this.currentVersion = Number(note.version || note.Version || this.currentVersion);
          this.expectedVersion = this.currentVersion;
        }
      });
    },
    async listNotes() {
      await this.wrapAction("拉取我的笔记", async () => {
        this.normalizeHosts();
        const data = await this.request(`${this.apiHost}/api/user/notes`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        const rows = data.notes || data.Notes || [];
        this.notes = rows.map((n) => ({
          noteId: n.noteId || n.NoteId || "",
          title: n.title || n.Title || "",
          version: Number(n.version || n.Version || 0),
        }));
        await this.listMyTeams(true);
      });
    },
    async pickNote(item) {
      this.noteId = item.noteId;
      await this.getNote();
      await this.listPermissions(true);
      await this.listEvents(true);
    },
    buildPermissionPayload() {
      const noteId = this.ensureNoteId();
      const targetUserId = this.targetUserId.trim();
      const targetUserEmail = this.targetUserEmail.trim();
      const targetTeamId = this.targetTeamId.trim();
      const hasUser = Boolean(targetUserId || targetUserEmail);
      const hasTeam = Boolean(targetTeamId);
      if ((!hasUser && !hasTeam) || (hasUser && hasTeam)) {
        throw new Error("用户(用户ID/邮箱)和团队ID必须二选一");
      }
      if (targetUserId && targetUserEmail) {
        throw new Error("targetUserId 和 targetUserEmail 不能同时填写");
      }
      return { noteId, targetUserId, targetUserEmail, targetTeamId };
    },
    async grantPermission() {
      await this.wrapAction("授予权限", async () => {
        this.normalizeHosts();
        const payload = {
          ...this.buildPermissionPayload(),
          role: this.permissionRole,
        };
        await this.request(`${this.apiHost}/api/permission/grant`, {
          method: "POST",
          headers: this.jsonHeaders(true),
          body: JSON.stringify(payload),
        });
        await this.listPermissions(true);
        await this.listEvents(true);
      });
    },
    async revokePermission() {
      await this.wrapAction("撤销权限", async () => {
        this.normalizeHosts();
        await this.request(`${this.apiHost}/api/permission/revoke`, {
          method: "POST",
          headers: this.jsonHeaders(true),
          body: JSON.stringify(this.buildPermissionPayload()),
        });
        await this.listPermissions(true);
        await this.listEvents(true);
      });
    },
    async listPermissions(silent = false) {
      const work = async () => {
        this.normalizeHosts();
        const noteId = this.ensureNoteId();
        const data = await this.request(`${this.apiHost}/api/permission/list/${encodeURIComponent(noteId)}`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        const rows = data.permissions || data.Permissions || [];
        this.permissions = rows.map((p) => ({
          permissionId: p.permissionId || p.PermissionId || "",
          noteId: p.noteId || p.NoteId || "",
          userId: p.userId || p.UserId || "",
          teamId: p.teamId || p.TeamId || "",
          role: p.role || p.Role || "",
          status: p.status || p.Status || "",
        }));
      };

      if (silent) {
        await work();
      } else {
        await this.wrapAction("查看权限列表", work);
      }
    },
    async listMyTeams(silent = false) {
      const work = async () => {
        this.normalizeHosts();
        const data = await this.request(`${this.apiHost}/api/team/me`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        const rows = data.teams || data.Teams || [];
        this.myTeams = rows.map((t) => ({
          teamId: t.teamId || t.TeamId || "",
          teamName: t.teamName || t.TeamName || "",
          role: t.role || t.Role || "",
          status: t.status || t.Status || "",
        }));
      };

      if (silent) {
        await work();
      } else {
        await this.wrapAction("查看我的团队", work);
      }
    },
    async listEvents(silent = false) {
      const work = async () => {
        this.normalizeHosts();
        const noteId = this.ensureNoteId();
        const qs = new URLSearchParams({
          startSeq: String(Number(this.eventStartSeq || 0)),
          limit: String(Number(this.eventLimit || 50)),
        });
        const data = await this.request(`${this.apiHost}/api/note/${encodeURIComponent(noteId)}/events?${qs.toString()}`, {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        const rows = data.events || data.Events || [];
        this.events = rows.map((e) => ({
          eventId: e.eventId || e.EventId || "",
          eventSeq: Number(e.eventSeq || e.EventSeq || 0),
          eventType: e.eventType || e.EventType || "",
          operatorId: e.operatorId || e.OperatorId || "",
          createdAt: e.createdAt || e.CreatedAt || 0,
        }));
      };

      if (silent) {
        await work();
      } else {
        await this.wrapAction("查看事件流", work);
      }
    },
  },
  mounted() {
    this.log(`Page origin: ${window.location.origin}`);
    this.log(`Auth Host: ${this.authHost} | API Host: ${this.apiHost}`);
    this.checkConnection();
  },
}).mount("#app");
