ALTER TABLE public.users
	ADD COLUMN IF NOT EXISTS full_name varchar NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS email varchar NULL,
	ADD COLUMN IF NOT EXISTS email_verified boolean NOT NULL DEFAULT true,
	ADD COLUMN IF NOT EXISTS verified_at timestamp NULL;

UPDATE public.users
SET email = username || '@local.invalid'
WHERE email IS NULL;

ALTER TABLE public.users
	ALTER COLUMN email SET NOT NULL;

ALTER TABLE public.users
	ADD CONSTRAINT users_email_unique UNIQUE (email);

CREATE TABLE IF NOT EXISTS public.user_email_verifications (
	verification_id uuid NOT NULL,
	user_id uuid NOT NULL,
	email varchar NOT NULL,
	code_hash varchar NOT NULL,
	expires_at timestamp NOT NULL,
	attempt_count integer NOT NULL DEFAULT 0,
	used_at timestamp NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT user_email_verifications_pk PRIMARY KEY (verification_id),
	CONSTRAINT user_email_verifications_user_fk FOREIGN KEY (user_id) REFERENCES public.users(user_id)
);

CREATE INDEX IF NOT EXISTS user_email_verifications_email_idx
	ON public.user_email_verifications (email);

CREATE INDEX IF NOT EXISTS user_email_verifications_active_idx
	ON public.user_email_verifications (email, expires_at)
	WHERE used_at IS NULL;
