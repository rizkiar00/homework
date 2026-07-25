CREATE TABLE IF NOT EXISTS public.user_password_resets (
	reset_id uuid NOT NULL,
	user_id uuid NOT NULL,
	email varchar NOT NULL,
	code_hash varchar NOT NULL,
	expires_at timestamp NOT NULL,
	attempt_count integer NOT NULL DEFAULT 0,
	used_at timestamp NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT user_password_resets_pk PRIMARY KEY (reset_id),
	CONSTRAINT user_password_resets_user_fk FOREIGN KEY (user_id) REFERENCES public.users(user_id)
);

CREATE INDEX IF NOT EXISTS user_password_resets_email_idx
	ON public.user_password_resets (email);

CREATE INDEX IF NOT EXISTS user_password_resets_active_idx
	ON public.user_password_resets (email, expires_at)
	WHERE used_at IS NULL;
