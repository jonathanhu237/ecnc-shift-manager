CREATE TABLE IF NOT EXISTS schedule_templates (
    id BIGSERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS schedule_template_shifts (
    id BIGSERIAL PRIMARY KEY,
    schedule_template_id BIGINT REFERENCES schedule_templates(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL,
    start_time TIME WITH TIME ZONE NOT NULL,
    end_time TIME WITH TIME ZONE NOT NULL,
    assistants_required INTEGER NOT NULL
);