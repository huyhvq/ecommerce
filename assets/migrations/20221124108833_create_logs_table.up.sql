create table logs
(
    id          uuid                                    default gen_random_uuid() PRIMARY KEY,
    search_term text                           not null default '',
    context     jsonb                          not null default '{}'::jsonb,
    created_at  timestamp(0) without time zone not null default NOW()::timestamp without time zone,
    updated_at  timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);
