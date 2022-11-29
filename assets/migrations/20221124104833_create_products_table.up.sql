create table products
(
    id          int primary key,
    name        text                           not null default '',
    description text                           not null default '',
    price       decimal(12, 2)                 not null default 0.00,
    attributes  jsonb                          not null default '{}'::jsonb,
    tsv         tsvector                       not null default array_to_tsvector('{}'),
    created_at  timestamp(0) without time zone not null default NOW()::timestamp without time zone,
    updated_at  timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);

CREATE FUNCTION product_text_search(name text, description text)
    RETURNS tsvector
    LANGUAGE sql
    IMMUTABLE AS
$function$
SELECT setweight(to_tsvector(name), 'A') ||
       setweight(to_tsvector(description), 'B');
$function$;

CREATE INDEX attributes_gin_idx ON products USING GIN (attributes jsonb_path_ops);
CREATE INDEX attributes_gin_tsv_idx ON products USING GIN (tsv);
CREATE INDEX text_idx ON products USING GIN (product_text_search(name, description));
