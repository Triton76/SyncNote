笔记：
对应实体：

- note_id string uuid primary
- user_id string uuid
- title string
- content string
- version int64
- created_at timestamp
- updated_at timestamp

同步策略：

- 客户端读取笔记时拿到 content + version
- 客户端保存时提交 note_id + content + expected_version
- 服务端比对 version
- 如果 expected_version == current_version，则保存成功，version + 1
- 如果不一致，则返回冲突，让前端提示“内容已更新，请刷新后重试”

可选的在线协作（先不做核心依赖）：

- redis:
  - note_id
  - online_users
  - last_heartbeat

对应方法接口：
syncNoteInterface {

- createNote(ctx, userId, title) (noteId, error)
- getNote(ctx, noteId) (\*noteModel, error)
- saveNote(ctx, noteId, content, expectedVersion) (\*noteModel, error)
- deleteNote(ctx, noteId) error
  }

teams:

- id uuid primary key
- name varchar(100) not null default "nothing233"
- slug varchar(100) not null unique
  team_members:
- team_id uuid not null
- user_id uuid not null,
- role VARCHAR(20),
- primary key (user_id, team_id)

note_permissions:

- permission_id uuid primary key
- team_id
- note_id
- user_id
- role
- granted_by
- created_at, updated_at

CHECK (user_id is not null and team_id is null) or (user_id is null and team_id is not null)
