CREATE TABLE IF NOT EXISTS schedule_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS schedule_template_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_template_id UUID NOT NULL REFERENCES schedule_templates(id) ON DELETE CASCADE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    required_assistants INTEGER NOT NULL
);

CREATE TABLE schedule_template_shifts_availability (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_template_shift_id UUID NOT NULL REFERENCES schedule_template_shifts(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL
);