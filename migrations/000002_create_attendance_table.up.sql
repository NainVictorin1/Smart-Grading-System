CREATE TABLE IF NOT EXISTS attendance (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    student_name TEXT NOT NULL,
    date DATE NOT NULL,
    status TEXT CHECK (status IN ('Present', 'Absent')) NOT NULL
);
