
CREATE TYPE course_level AS ENUM ('beginner', 'elementary', 'intermediate', 'ielts');

CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    group_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    classroom VARCHAR(20),
    type_of_group course_level,
    task VARCHAR(35),
    deadline TIMESTAMP,
    score NUMERIC,
    start_time TIME,
    end_time TIME,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lessons (
    id UUID PRIMARY KEY,
    theme VARCHAR(35),
    links VARCHAR(255),
    type_of_group course_level,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS journals (
    id UUID PRIMARY KEY,
    schedule_id UUID NOT NULL,
    date_of_lesson DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);


CREATE TABLE IF NOT EXISTS student_performance (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    schedule_id UUID NOT NULL,
    attended BOOLEAN NOT NULL,
    task_score NUMERIC,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
