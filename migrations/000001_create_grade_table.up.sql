CREATE TABLE IF NOT EXISTS grade(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    fullname text NOT NULL,
    subject text NOT NULL,
    grade NUMERIC NOT Null,
    email citext NOT NUll

);