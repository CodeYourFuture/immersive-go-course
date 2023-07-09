CREATE TABLE IF NOT EXISTS public.images
(
    id         serial                              NOT NULL,
    title      text                                NOT NULL,
    url        text                                NOT NULL UNIQUE,
    alt_text   text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    resolution text,

    PRIMARY KEY (id)
);