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
- getNote(ctx, noteId) (*noteModel, error)
- saveNote(ctx, noteId, content, expectedVersion) (*noteModel, error)
- deleteNote(ctx, noteId) error
}