DROP DATABASE IF EXISTS syncnote;
CREATE DATABASE syncnote;
USE syncnote;
-- ============================================================
-- 数据库设置
-- ============================================================
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- ============================================================
-- 1. User 表（用户基础表）
-- ============================================================
CREATE TABLE `user` (
	`user_id` CHAR(36) NOT NULL COMMENT '用户ID (UUID)',
	`email` VARCHAR(255) NOT NULL COMMENT '邮箱',
	`password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
	`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	`deleted_at` DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
	PRIMARY KEY (`user_id`),
	UNIQUE KEY `uk_email` (`email`),
	KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';
-- ============================================================
-- 2. UserInfo 表（用户详细信息）
-- ============================================================
CREATE TABLE `user_info` (
	`user_id` CHAR(36) NOT NULL COMMENT '用户ID (UUID)',
	`username` VARCHAR(50) NOT NULL COMMENT '用户名',
	`synopsis` TEXT NULL COMMENT '个人简介',
	`avatar_url` VARCHAR(255) NULL COMMENT '头像URL',
	`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	PRIMARY KEY (`user_id`),
	CONSTRAINT `fk_user_info_user` FOREIGN KEY (`user_id`) REFERENCES `user`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户信息表';
-- ============================================================
-- 3. Team 表（团队表）
-- ============================================================
CREATE TABLE `team` (
	`team_id` CHAR(36) NOT NULL COMMENT '团队ID (UUID)',
	`name` VARCHAR(100) NOT NULL COMMENT '团队名称',
	`description` TEXT NULL COMMENT '团队描述',
	`owner_id` CHAR(36) NOT NULL COMMENT '团队所有者ID',
	`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	`deleted_at` DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
	PRIMARY KEY (`team_id`),
	KEY `idx_owner_id` (`owner_id`),
	KEY `idx_deleted_at` (`deleted_at`),
	CONSTRAINT `fk_team_owner` FOREIGN KEY (`owner_id`) REFERENCES `user`(`user_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '团队表';
-- ============================================================
-- 4. Note 表（笔记表）
-- ============================================================
CREATE TABLE `note` (
	`note_id` CHAR(36) NOT NULL COMMENT '笔记ID (UUID)',
	`owner_id` CHAR(36) NOT NULL COMMENT '笔记所有者ID',
	`title` VARCHAR(200) NOT NULL COMMENT '笔记标题',
	`content` TEXT NULL COMMENT '笔记内容',
	`version` INT NOT NULL DEFAULT 1 COMMENT '版本号 (乐观锁)',
	`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	`deleted_at` DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
	PRIMARY KEY (`note_id`),
	KEY `idx_owner_id` (`owner_id`),
	KEY `idx_deleted_at` (`deleted_at`),
	CONSTRAINT `fk_note_owner` FOREIGN KEY (`owner_id`) REFERENCES `user`(`user_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '笔记表';
-- ============================================================
-- 5. Note_User_Permission 表（笔记 - 用户权限表）
-- ============================================================
CREATE TABLE `note_user_permission` (
	`permission_id` char(36) PRIMARY KEY,
	`note_id` CHAR(36) NOT NULL COMMENT '笔记ID',
	`user_id` CHAR(36) NOT NULL COMMENT '用户ID',
	`permission_level` VARCHAR(20) NOT NULL COMMENT '权限级别 (read/write/admin)',
	`granted_by` CHAR(36) NULL DEFAULT NULL COMMENT '授权人ID (允许NULL用于SET NULL)',
	`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	KEY `idx_user_id` (`user_id`),
	KEY `idx_granted_by` (`granted_by`),
	UNIQUE KEY `uk_note_user` (`note_id`, `user_id`),
	CONSTRAINT `fk_nup_note` FOREIGN KEY (`note_id`) REFERENCES `note`(`note_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `fk_nup_user` FOREIGN KEY (`user_id`) REFERENCES `user`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `fk_nup_granted_by` FOREIGN KEY (`granted_by`) REFERENCES `user`(`user_id`) ON DELETE
	SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '笔记用户权限表';
-- ============================================================
-- 6. Note_Team_Permission 表（笔记 - 团队权限表）
-- ============================================================
CREATE TABLE `note_team_permission` (
	`permission_id` char(36) PRIMARY KEY,
	`note_id` CHAR(36) NOT NULL COMMENT '笔记ID',
	`team_id` CHAR(36) NOT NULL COMMENT '团队ID',
	`permission_level` VARCHAR(20) NOT NULL COMMENT '权限级别 (read/write/admin)',
	`granted_by` CHAR(36) NULL DEFAULT NULL COMMENT '授权人ID (允许NULL用于SET NULL)',
	`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	UNIQUE KEY `uk_note_team` (`note_id`, `team_id`),
	KEY `idx_team_id` (`team_id`),
	KEY `idx_granted_by` (`granted_by`),
	CONSTRAINT `fk_ntp_note` FOREIGN KEY (`note_id`) REFERENCES `note`(`note_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `fk_ntp_team` FOREIGN KEY (`team_id`) REFERENCES `team`(`team_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `fk_ntp_granted_by` FOREIGN KEY (`granted_by`) REFERENCES `user`(`user_id`) ON DELETE
	SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '笔记团队权限表';
-- ============================================================
-- 7. Team_members 表（团队成员表）
-- ============================================================
CREATE TABLE `team_members` (
	`id` CHAR(36) PRIMARY KEY,
	`team_id` CHAR(36) NOT NULL COMMENT '团队ID',
	`user_id` CHAR(36) NOT NULL COMMENT '用户ID',
	`joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
	UNIQUE KEY `uk_team_user` (`team_id`, `user_id`),
	KEY `idx_user_id` (`user_id`),
	CONSTRAINT `fk_tm_team` FOREIGN KEY (`team_id`) REFERENCES `team`(`team_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `fk_tm_user` FOREIGN KEY (`user_id`) REFERENCES `user`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '团队成员表';
-- ============================================================
-- 恢复外键检查
-- ============================================================
SET FOREIGN_KEY_CHECKS = 1;