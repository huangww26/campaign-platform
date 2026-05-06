-- 建库: createdb campaign_platform
-- 执行: psql -d campaign_platform -f migrations/001_init.sql

CREATE TABLE IF NOT EXISTS campaign_templates (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(64) UNIQUE NOT NULL,
    schema_def  JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS campaigns (
    id          SERIAL PRIMARY KEY,
    template_id INT REFERENCES campaign_templates(id),
    name        VARCHAR(128) NOT NULL,
    slug        VARCHAR(128) UNIQUE NOT NULL,
    status      VARCHAR(16) NOT NULL DEFAULT 'draft',
    config      JSONB NOT NULL DEFAULT '{}',
    version     INT NOT NULL DEFAULT 1,
    created_by  VARCHAR(64),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS components (
    name         VARCHAR(64) PRIMARY KEY,
    version      VARCHAR(16) NOT NULL DEFAULT 'v1',
    category     VARCHAR(16) NOT NULL DEFAULT 'basic',
    props_schema JSONB NOT NULL DEFAULT '{}',
    status       VARCHAR(16) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS campaign_versions (
    id          SERIAL PRIMARY KEY,
    campaign_id INT NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    version     INT NOT NULL,
    config      JSONB NOT NULL,
    changelog   TEXT,
    deployed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);
CREATE INDEX IF NOT EXISTS idx_campaigns_slug ON campaigns(slug);
CREATE INDEX IF NOT EXISTS idx_campaign_versions_campaign ON campaign_versions(campaign_id);

-- 预置组件
INSERT INTO components (name, version, category, props_schema) VALUES
    ('hero_banner', 'v1', 'basic', '{
        "type":"object",
        "properties":{
            "bg_image":{"type":"string"},
            "title":{"type":"object"},
            "subtitle":{"type":"object"},
            "cta_text":{"type":"object"},
            "cta_link":{"type":"string"}
        }
    }'),
    ('cta_button', 'v1', 'basic', '{
        "type":"object",
        "properties":{
            "text":{"type":"object"},
            "action":{"type":"string"}
        }
    }')
ON CONFLICT (name) DO NOTHING;
