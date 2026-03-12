-- 笔记表
CREATE TABLE notes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    content_type VARCHAR(20) DEFAULT 'text',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id)
);

-- 笔记操作日志表（用于冲突解决）
CREATE TABLE note_operations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    note_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    op_type VARCHAR(20) NOT NULL,  -- insert, delete, replace
    start_pos INT NOT NULL,
    end_pos INT NOT NULL,
    content TEXT,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_note_timestamp (note_id, timestamp)
);