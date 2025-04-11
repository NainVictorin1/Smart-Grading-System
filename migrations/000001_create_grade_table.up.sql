CREATE TABLE grade (
    id BIGINT PRIMARY KEY DEFAULT nextval('grade_id_seq'::regclass),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
    fullname TEXT NOT NULL,
    subject TEXT NOT NULL,
    grade NUMERIC(5, 2) NOT NULL,
    email CITEXT NOT NULL
);