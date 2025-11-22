create table public.auth (
  password_hash text null,
  google_id character varying(100) null,
  email character varying(100) null,
  created_at timestamp without time zone null default CURRENT_TIMESTAMP,
  updated_at timestamp without time zone null default CURRENT_TIMESTAMP,
  github_id character varying null,
  id uuid not null,
  constraint auth_pkey primary key (id),
  constraint auth_google_id_key unique (google_id),
  constraint auth_id_key unique (id),
  constraint auth_id_fkey foreign KEY (id) references users (id)
) TABLESPACE pg_default;

create trigger set_updated_at BEFORE
update on auth for EACH row
execute FUNCTION update_updated_at_column ();
