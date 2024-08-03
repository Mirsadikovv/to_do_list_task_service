ALTER TABLE schedules ADD FOREIGN KEY (lesson_id) REFERENCES lessons (id);

ALTER TABLE journals ADD FOREIGN KEY (schedule_id) REFERENCES schedules (id);

ALTER TABLE student_performance ADD FOREIGN KEY (schedule_id) REFERENCES schedules (id);
