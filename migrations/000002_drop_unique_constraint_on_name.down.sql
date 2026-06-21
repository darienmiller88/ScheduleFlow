ALTER TABLE specialists ADD CONSTRAINT specialists_first_name_key UNIQUE (first_name);
ALTER TABLE specialists ADD CONSTRAINT specialists_last_name_key UNIQUE (last_name);