CREATE TABLE IF NOT EXISTS public.users (
	id_user uuid NOT NULL,
	username varchar(100) NOT NULL,
	password_hash varchar NOT NULL,
	role varchar(50) NOT NULL DEFAULT 'user',
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT users_pk PRIMARY KEY (id_user),
	CONSTRAINT users_username_unique UNIQUE (username)
);
