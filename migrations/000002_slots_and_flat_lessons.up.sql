CREATE TABLE IF NOT EXISTS slots (
    id SERIAL PRIMARY KEY,
    slot_number INT NOT NULL CHECK (slot_number >= 1 AND slot_number <= 20),
    start_minute INT NOT NULL CHECK (start_minute >= 0 AND start_minute < 1440),
    end_minute INT NOT NULL CHECK (end_minute > start_minute AND end_minute <= 1440),
    day_of_week INT CHECK (day_of_week IS NULL OR (day_of_week >= 0 AND day_of_week <= 6)),
    CONSTRAINT uq_slot_number_day UNIQUE (slot_number, day_of_week)
);
CREATE INDEX IF NOT EXISTS idx_slots_day_number ON slots(day_of_week, slot_number);

CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    teacher_id INT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id INT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    classroom VARCHAR(50) NOT NULL,
    day_of_week INT NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    slot_number INT NOT NULL CHECK (slot_number >= 1),
    slot_length INT NOT NULL DEFAULT 1 CHECK (slot_length >= 1),
    week_type VARCHAR(10) NOT NULL DEFAULT 'EVEN' CHECK (week_type IN ('EVEN', 'ODD', 'BOTH', 'OTHER')),
    is_available BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_lessons_group_day ON lessons(group_id, day_of_week, week_type);
CREATE INDEX IF NOT EXISTS idx_lessons_teacher_day ON lessons(teacher_id, day_of_week, week_type);

CREATE TABLE IF NOT EXISTS replaces (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    slot_number INT NOT NULL CHECK (slot_number >= 1),
    group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    teacher_id INT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id INT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    classroom VARCHAR(50) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_replaces_group_date ON replaces(group_id, date);
CREATE INDEX IF NOT EXISTS idx_replaces_teacher_date ON replaces(teacher_id, date);

-- Default bell timetable
INSERT INTO slots (slot_number, start_minute, end_minute, day_of_week) VALUES
(1, 510, 600, NULL),   -- 08:30 - 10:00
(2, 610, 700, NULL),   -- 10:10 - 11:40
(3, 720, 810, NULL),   -- 12:00 - 13:30
(4, 830, 920, NULL),   -- 13:50 - 15:20
(5, 930, 1020, NULL),  -- 15:30 - 17:00
(6, 1030, 1120, NULL), -- 17:10 - 18:40
(7, 1130, 1220, NULL)  -- 18:50 - 20:20
ON CONFLICT DO NOTHING;
