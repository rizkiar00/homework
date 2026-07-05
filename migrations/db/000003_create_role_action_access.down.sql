ALTER TABLE public.users DROP CONSTRAINT IF EXISTS users_role_fk;

ALTER TABLE public.users
	DROP COLUMN IF EXISTS updated_at,
	DROP COLUMN IF EXISTS updated_by,
	DROP COLUMN IF EXISTS created_by,
	DROP COLUMN IF EXISTS is_active,
	DROP COLUMN IF EXISTS role_id;

ALTER TABLE public.test_table
	DROP COLUMN IF EXISTS updated_at,
	DROP COLUMN IF EXISTS updated_by,
	DROP COLUMN IF EXISTS created_at,
	DROP COLUMN IF EXISTS created_by,
	DROP COLUMN IF EXISTS is_active;

DROP TABLE IF EXISTS public.role_accesses;
DROP TABLE IF EXISTS public.roles;
DROP TABLE IF EXISTS public.actions;

ALTER TABLE public.users RENAME COLUMN user_id TO id_user;

ALTER TABLE public.test_table RENAME COLUMN test_id TO id_test;
