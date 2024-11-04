
-- organizations 
CREATE TABLE IF NOT EXISTS organizations (
    organization_id INT GENERATED ALWAYS AS IDENTITY,
    organization_name TEXT NOT NULL,
    organization_description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(organization_id),
    UNIQUE(organization_name)
);

-- user that can access app
CREATE TABLE IF NOT EXISTS app_users (
    user_id INT GENERATED ALWAYS AS IDENTITY,
    organization_id INTEGER NOT NULL REFERENCES organizations(organization_id),
    username TEXT NOT NULL,
    display_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id),
    UNIQUE(username),
    UNIQUE(email)
);

