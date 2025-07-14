create table public.t_ssh_key (
    id bigserial primary key,
    user_id bigint not null,
    name character varying(255) not null,
    public_key character varying(4096) not null,
    gmt_created timestamp without time zone default current_timestamp not null,
    gmt_updated timestamp without time zone not null,
    gmt_deleted timestamp without time zone
);
