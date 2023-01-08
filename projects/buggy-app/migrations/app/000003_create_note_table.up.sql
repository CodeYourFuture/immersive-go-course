-- Create notes table
CREATE TABLE IF NOT EXISTS public.note(
   id VARCHAR (20) PRIMARY KEY,
   owner VARCHAR (20) NOT NULL REFERENCES public.user (id),
   content TEXT NOT NULL default '',
   created timestamp default current_timestamp,
   modified timestamp default current_timestamp
);

-- Add "modifier" trigger to note
CREATE TRIGGER note_update_modified
BEFORE UPDATE ON public.note
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- Add short ID trigger to note
CREATE TRIGGER note_gen_id
BEFORE INSERT ON public.note
FOR EACH ROW EXECUTE PROCEDURE gen_id();