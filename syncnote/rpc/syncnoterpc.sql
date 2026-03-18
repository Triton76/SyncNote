-- 创建数据库 (如果尚未创建)
CREATE DATABASE IF NOT EXISTS syncnote DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE syncnote;
-- 创建笔记表
CREATE TABLE IF NOT EXISTS notes (
	-- 对应 NoteReq.note_id / NoteResp.note_id
	-- 建议使用 VARCHAR 存储 UUID 或 Snowflake ID 字符串
	note_id VARCHAR(32) NOT NULL DEFAULT '' COMMENT '笔记唯一ID',
	-- 对应 user_id
	user_id VARCHAR(32) NOT NULL DEFAULT '' COMMENT '所属用户ID',
	-- 对应 title
	title VARCHAR(255) NOT NULL DEFAULT '' COMMENT '笔记标题',
	-- 对应 content (使用 LONGTEXT 存储大文本),这里content用NOT NULL原因是允许空字符串但是不允许NULL （在数据库语境中NULL表示不存在，而非空字符串）
	content LONGTEXT NOT NULL DEFAULT '' COMMENT '笔记内容',
	-- 对应 version (乐观锁核心字段)
	-- 初始版本通常为 1
	version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
	-- 对应 last_modified (存储 Unix 时间戳，单位毫秒或秒)
	last_modified BIGINT NOT NULL DEFAULT 0 COMMENT '最后修改时间戳',
	-- 软删除标记 (可选，但推荐)
	is_deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0:否 1:是',
	-- 创建时间 (额外补充，Proto中未体现但数据库通常需要)
	created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间戳',
	-- 主键约束
	PRIMARY KEY (note_id),
	-- 索引优化
	INDEX idx_user_id (user_id),
	-- 快速查询用户的所有笔记 (GetUserNotes)
	INDEX idx_user_modified (user_id, last_modified DESC) -- 优化按时间排序的列表查询
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户笔记表';