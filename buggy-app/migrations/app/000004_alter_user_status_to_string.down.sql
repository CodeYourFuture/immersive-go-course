ALTER TABLE public.user ADD status_int int NOT NULL default 0;

UPDATE public.user SET status_int = CASE WHEN status = 'active' THEN 1 ELSE 0 END;

ALTER TABLE public.user DROP COLUMN status;

ALTER TABLE public.user RENAME COLUMN status_int TO status;