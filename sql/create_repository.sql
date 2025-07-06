create table public.t_repository (
    id bigserial primary key,
    repository_name character varying(255) not null,
    repository_description character varying(255) not null,
    is_private boolean default false,
    user_id bigint not null,
    star integer default 0 not null,
    fork integer default 0 not null,
    watcher integer default 0 not null,
    https_url character varying(1024) not null,
    ssh_url character varying(1024) not null,
    gmt_created timestamp without time zone default current_timestamp not null,
    gmt_updated timestamp without time zone default current_timestamp not null,
    gmt_deleted timestamp without time zone
);

COMMENT ON TABLE public.t_repository IS 'Table for storing repository information.';

COMMENT ON COLUMN public.t_repository.id IS 'Primary key of the repository table.';
COMMENT ON COLUMN public.t_repository.repository_name IS 'Name of the repository.';
COMMENT ON COLUMN public.t_repository.repository_description IS 'Description of the repository.';
COMMENT ON COLUMN public.t_repository.is_private IS 'Indicates if the repository is private.';
COMMENT ON COLUMN public.t_repository.user_id IS 'ID of the user who owns the repository.';
COMMENT ON COLUMN public.t_repository.star IS 'Number of stars the repository has received.';
COMMENT ON COLUMN public.t_repository.fork IS 'Number of times the repository has been forked.';
COMMENT ON COLUMN public.t_repository.watcher IS 'Number of users watching the repository.';
COMMENT ON COLUMN public.t_repository.https_url IS 'Repository link under HTTPS protocol.';
COMMENT ON COLUMN public.t_repository.ssh_url IS 'Repository link under SSH protocol.';
COMMENT ON COLUMN public.t_repository.gmt_created IS 'Timestamp when the repository was created.';
COMMENT ON COLUMN public.t_repository.gmt_updated IS 'Timestamp when the repository was last updated.';
COMMENT ON COLUMN public.t_repository.gmt_deleted IS 'Timestamp when the repository was deleted.
If set to NULL, it indicates that the repository has not been deleted.';

