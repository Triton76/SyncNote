CREATE DATABASE IF NOT EXISTS syncnote DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE syncnote;
CREATE TABLE IF NOT EXISTS notes (
	note_id CHAR(36) NOT NULL COMMENT '笔记唯一ID',
	user_id CHAR(36) NOT NULL COMMENT '所属用户ID',
	title VARCHAR(255) NOT NULL DEFAULT '' COMMENT '笔记标题',
	content LONGTEXT NOT NULL COMMENT '笔记内容',
	version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号(乐观锁)',
	last_modified BIGINT NOT NULL DEFAULT 0 COMMENT '最后修改时间戳(ms)',
	is_deleted BIGINT NOT NULL DEFAULT 0 COMMENT '是否删除 0:否 1:是',
	created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间戳(ms)',
	PRIMARY KEY (note_id),
	KEY idx_user_id (user_id),
	KEY idx_user_deleted (user_id, is_deleted),
	KEY idx_last_modified (last_modified)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '笔记表';