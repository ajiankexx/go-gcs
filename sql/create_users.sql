CREATE TABLE public.t_users (
    id BIGSERIAL PRIMARY KEY,   -- 使用 BIGSERIAL 自动生成自增的 ID
    username character varying(50) NOT NULL,
    email character varying(254) NOT NULL,
    user_password character(32) NOT NULL,
    avatar_url character varying(1024) NOT NULL,
    gmt_created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    gmt_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    gmt_deleted timestamp without time zone
);

COMMENT ON TABLE public.t_users IS 'table for storing users information';
COMMENT ON COLUMN public.t_users.id IS 'Primary key for users';
COMMENT ON COLUMN public.t_users.username IS 'username for users';
COMMENT ON COLUMN public.t_users.email IS 'email for users';
COMMENT ON COLUMN public.t_users.user_password IS 'user password for users';
