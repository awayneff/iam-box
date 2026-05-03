-- Create extension for better performance
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Permissions table (who can do what)
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    granted_at TIMESTAMP DEFAULT NOW(),
    created_by TEXT
);

-- Audit log (every decision)
CREATE TABLE decisions (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    allowed BOOLEAN NOT NULL,
    reason TEXT,
    timestamp TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_permissions_user ON permissions(user_id);
CREATE INDEX idx_permissions_resource ON permissions(resource_type, resource_id);
CREATE INDEX idx_decisions_user ON decisions(user_id);
CREATE INDEX idx_decisions_timestamp ON decisions(timestamp);