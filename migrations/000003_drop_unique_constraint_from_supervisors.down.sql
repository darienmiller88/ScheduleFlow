ALTER TABLE supervisors ADD CONSTRAINT supervisors_first_name_key UNIQUE (first_name);
ALTER TABLE supervisors ADD CONSTRAINT supervisors_last_name_key UNIQUE (last_name);