DROP TABLE IF EXISTS public.user_email_verifications;

ALTER TABLE public.users
	DROP CONSTRAINT IF EXISTS users_email_unique;

ALTER TABLE public.users
	DROP COLUMN IF EXISTS verified_at,
	DROP COLUMN IF EXISTS email_verified,
	DROP COLUMN IF EXISTS email,
	DROP COLUMN IF EXISTS full_name;
