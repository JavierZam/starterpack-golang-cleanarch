-- migrations/000001_create_users_table.up.sql (SIMPLIFIED FOR TESTING - DIRECT EXECUTION)
-- This version uses INTEGER ID and VARCHAR for tenant_id.
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_users_tenant_id ON users (tenant_id);
CREATE UNIQUE INDEX idx_users_email_tenant_id ON users (email, tenant_id);
