ALTER TABLE public.test_table RENAME COLUMN id_test TO test_id;

ALTER TABLE public.test_table
	ADD COLUMN is_active boolean NOT NULL DEFAULT true,
	ADD COLUMN created_by uuid NULL,
	ADD COLUMN created_at timestamp NOT NULL DEFAULT now(),
	ADD COLUMN updated_by uuid NULL,
	ADD COLUMN updated_at timestamp NULL;

ALTER TABLE public.users RENAME COLUMN id_user TO user_id;

CREATE TABLE IF NOT EXISTS public.actions (
	action_id bigserial NOT NULL,
	action_desc varchar NOT NULL,
	action_type varchar NOT NULL,
	endpoint varchar NOT NULL,
	CONSTRAINT actions_pk PRIMARY KEY (action_id),
	CONSTRAINT actions_unique_action UNIQUE (action_type, endpoint),
	CONSTRAINT actions_action_type_check CHECK (action_type IN ('GET', 'POST', 'PUT', 'PATCH', 'DELETE'))
);

CREATE TABLE IF NOT EXISTS public.roles (
	role_id bigserial NOT NULL,
	role_desc varchar NOT NULL,
	is_active boolean NOT NULL DEFAULT true,
	created_by uuid NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_by uuid NULL,
	updated_at timestamp NULL,
	CONSTRAINT roles_pk PRIMARY KEY (role_id),
	CONSTRAINT roles_unique_desc UNIQUE (role_desc)
);

CREATE TABLE IF NOT EXISTS public.role_accesses (
	role_access_id uuid NOT NULL,
	role_id bigint NOT NULL,
	action_id bigint NOT NULL,
	created_by uuid NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_by uuid NULL,
	updated_at timestamp NULL,
	CONSTRAINT role_accesses_pk PRIMARY KEY (role_access_id),
	CONSTRAINT role_accesses_unique_role_action UNIQUE (role_id, action_id),
	CONSTRAINT role_accesses_role_fk FOREIGN KEY (role_id) REFERENCES public.roles(role_id),
	CONSTRAINT role_accesses_action_fk FOREIGN KEY (action_id) REFERENCES public.actions(action_id)
);

ALTER TABLE public.users
	ADD COLUMN role_id bigint NULL,
	ADD COLUMN is_active boolean NOT NULL DEFAULT true,
	ADD COLUMN created_by uuid NULL,
	ADD COLUMN updated_by uuid NULL,
	ADD COLUMN updated_at timestamp NULL;

ALTER TABLE public.users
	ADD CONSTRAINT users_role_fk FOREIGN KEY (role_id) REFERENCES public.roles(role_id);
