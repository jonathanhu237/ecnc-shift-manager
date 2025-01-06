CREATE TABLE IF NOT EXISTS schedule_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    submission_start_time TIMESTAMPTZ NOT NULL,
    submission_end_time TIMESTAMPTZ NOT NULL,
    active_start_time TIMESTAMPTZ NOT NULL,
    active_end_time TIMESTAMPTZ NOT NULL,
    schedule_template_name TEXT NOT NULL REFERENCES schedule_templates(name),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);