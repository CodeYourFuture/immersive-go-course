ALTER TABLE public.user ADD status_str VARCHAR(20) NOT NULL DEFAULT 'inactive';

UPDATE public.user SET status_str = CASE WHEN "status" = 1 THEN 'active' ELSE 'inactive' END;

ALTER TABLE public.user DROP COLUMN status;

ALTER TABLE public.user RENAME COLUMN status_str TO status;