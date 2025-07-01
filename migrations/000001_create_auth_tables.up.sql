    -- migrations/000001_create_auth_tables.up.sql
    -- This migration creates the 'users' table, which is central for authentication.
    -- It uses UUIDs for primary keys and tenant_id for multi-tenancy.
    -- Requires: PostgreSQL with "uuid-ossp" extension enabled.

    CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Ensure UUID generation extension is available

    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Using UUID as primary key
        tenant_id UUID NOT NULL,                        -- Tenant ID is also a UUID
        email VARCHAR(255) UNIQUE NOT NULL,             -- Unique email for login
        password_hash VARCHAR(255) NOT NULL,            -- Stores hashed password (e.g., bcrypt)
        name VARCHAR(255) NOT NULL,
        phone_number VARCHAR(50),                       -- Optional: phone number can be nullable
        role VARCHAR(50) NOT NULL DEFAULT 'user',       -- User role (e.g., 'admin', 'user')
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    -- Indexes for performance
    CREATE UNIQUE INDEX idx_users_email ON users (email); -- Unique index on email globally
    CREATE INDEX idx_users_tenant_id ON users (tenant_id); -- Index on tenant_id for filtering