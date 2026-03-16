CREATE DATABASE IF NOT EXISTS syncnote 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE syncnote;

CREATE TABLE IF NOT EXISTS users (
    -- 主键（单独声明，goctl 必须这样识别）
    id CHAR(36) NOT NULL COMMENT '用户 ID(UUID)',
    
    -- 登录相关（添加 DEFAULT 消除警告）
    password_hash VARCHAR(255) NOT NULL DEFAULT '' COMMENT '加密后的密码',
    
    -- 用户信息（允许 NULL，保持原样）
    email VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    username VARCHAR(50) DEFAULT NULL COMMENT '用户昵称',
    
    -- 状态管理
    status TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态:1 正常 0 禁用',
    
    -- 时间戳（添加 DEFAULT 0 消除警告）
    created_at BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间戳',
    
    -- 主键约束（单独一行，goctl 关键！）
    PRIMARY KEY (id),
    
    -- 索引
    UNIQUE KEY uk_email (email) COMMENT '邮箱唯一',
    KEY idx_status (status) COMMENT '状态索引',
    KEY idx_created (created_at) COMMENT '创建时间索引'
    
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='用户表';