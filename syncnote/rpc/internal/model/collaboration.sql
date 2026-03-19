CREATE DATABASE IF NOT EXISTS syncnote DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE syncnote;
CREATE TABLE IF NOT EXISTS note_permissions (
	permission_id CHAR(32) NOT NULL COMMENT '权限记录唯一ID',
	note_id CHAR(32) NOT NULL COMMENT '笔记ID',
	user_id CHAR(32) DEFAULT NULL COMMENT '被授权用户ID，user_id 和 team_id 二选一',
	team_id CHAR(32) DEFAULT NULL COMMENT '被授权团队ID（团队成员自动继承权限）',
	granted_by CHAR(32) NOT NULL COMMENT '授权者ID（note owner or admin）',
	role VARCHAR(20) NOT NULL DEFAULT 'viewer' COMMENT '权限角色：owner | admin | editor | viewer',
	status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'active | revoked | pending',
	granted_at BIGINT NOT NULL DEFAULT 0 COMMENT '授权时间戳(ms)',
	revoked_at BIGINT DEFAULT NULL COMMENT '撤销时间戳(ms)',
	PRIMARY KEY (permission_id),
	KEY idx_note_id (note_id),
	KEY idx_user_id (user_id),
	KEY idx_team_id (team_id),
	KEY idx_note_user (note_id, user_id),
	KEY idx_note_team (note_id, team_id),
	UNIQUE KEY uk_note_user (note_id, user_id),
	UNIQUE KEY uk_note_team (note_id, team_id),
	CONSTRAINT chk_permission_target CHECK (
		(
			user_id IS NOT NULL
			AND team_id IS NULL
		)
		OR (
			user_id IS NULL
			AND team_id IS NOT NULL
		)
	),
	CONSTRAINT chk_permission_role CHECK (role IN ('owner', 'admin', 'editor', 'viewer')),
	CONSTRAINT chk_permission_status CHECK (status IN ('active', 'revoked', 'pending'))
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '笔记权限表';
CREATE TABLE IF NOT EXISTS collaboration_events (
	event_id CHAR(32) NOT NULL COMMENT '事件唯一ID',
	note_id CHAR(32) NOT NULL COMMENT '操作的笔记ID',
	event_seq BIGINT NOT NULL COMMENT '该笔记内事件序列号（单调递增）',
	event_type VARCHAR(30) NOT NULL COMMENT '事件类型',
	operator_id CHAR(32) NOT NULL COMMENT '操作者用户ID',
	operator_name VARCHAR(255) DEFAULT NULL COMMENT '操作者用户名（冗余字段）',
	payload JSON DEFAULT NULL COMMENT '事件具体内容',
	note_version BIGINT DEFAULT NULL COMMENT '操作后笔记版本号',
	expected_version BIGINT DEFAULT NULL COMMENT '预期版本号（用于冲突检测）',
	is_conflict TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否冲突：0否 1是',
	related_event_id CHAR(32) DEFAULT NULL COMMENT '关联冲突事件ID',
	created_at BIGINT NOT NULL DEFAULT 0 COMMENT '事件发生时间戳(ms)',
	PRIMARY KEY (event_id),
	KEY idx_note_id (note_id),
	KEY idx_operator_id (operator_id),
	KEY idx_event_type (event_type),
	KEY idx_created_at (created_at),
	UNIQUE KEY uk_note_seq (note_id, event_seq),
	KEY idx_note_created (note_id, created_at),
	CONSTRAINT chk_event_type CHECK (
		event_type IN (
			'note_created',
			'note_updated',
			'note_deleted',
			'permission_granted',
			'permission_revoked',
			'conflict_detected',
			'view_started',
			'view_ended'
		)
	)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '协作事件流表';