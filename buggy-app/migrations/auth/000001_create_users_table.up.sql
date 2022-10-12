CREATE TABLE IF NOT EXISTS public.user(
   id serial PRIMARY KEY,
   status int NOT NULL,
   password VARCHAR (100) NOT NULL,
   created timestamp default current_timestamp,
   modified timestamp default current_timestamp
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_modified
BEFORE UPDATE ON public.user
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();