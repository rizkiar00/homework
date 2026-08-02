INSERT INTO public.actions (action_id, action_desc, action_type, endpoint)
VALUES
	(1, 'Get current user', 'GET', '/auth/me'),
	(2, 'Create test DB row', 'POST', '/test_db'),
	(3, 'Get test DB list', 'GET', '/test_db'),
	(4, 'Get test DB detail', 'GET', '/test_db/{test_id}'),
	(5, 'Update test DB row', 'PUT', '/test_db/{test_id}'),
	(6, 'Delete test DB row', 'DELETE', '/test_db/{test_id}'),
	(7, 'Get action list', 'GET', '/actions'),
	(8, 'Create role', 'POST', '/roles'),
	(9, 'Update role', 'PUT', '/roles/{role_id}'),
	(10, 'Set role actions', 'PUT', '/roles/{role_id}/actions'),
	(11, 'Assign user role', 'PUT', '/users/{user_id}/role'),
	(12, 'Logout current user', 'POST', '/auth/logout')
ON CONFLICT (action_type, endpoint) DO UPDATE
SET action_desc = EXCLUDED.action_desc;

SELECT setval(pg_get_serial_sequence('public.actions', 'action_id'), (SELECT MAX(action_id) FROM public.actions));

INSERT INTO public.roles (role_id, role_desc, is_active, created_at)
VALUES
	(1, 'admin', true, now()),
	(2, 'user', true, now())
ON CONFLICT (role_desc) DO UPDATE
SET is_active = EXCLUDED.is_active,
	updated_at = now();

SELECT setval(pg_get_serial_sequence('public.roles', 'role_id'), (SELECT MAX(role_id) FROM public.roles));

INSERT INTO public.role_accesses (role_access_id, role_id, action_id, created_at)
VALUES
	(gen_random_uuid(), 1, 1, now()),
	(gen_random_uuid(), 1, 2, now()),
	(gen_random_uuid(), 1, 3, now()),
	(gen_random_uuid(), 1, 4, now()),
	(gen_random_uuid(), 1, 5, now()),
	(gen_random_uuid(), 1, 6, now()),
	(gen_random_uuid(), 1, 7, now()),
	(gen_random_uuid(), 1, 8, now()),
	(gen_random_uuid(), 1, 9, now()),
	(gen_random_uuid(), 1, 10, now()),
	(gen_random_uuid(), 1, 11, now()),
	(gen_random_uuid(), 1, 12, now()),
	(gen_random_uuid(), 2, 2, now()),
	(gen_random_uuid(), 2, 3, now()),
	(gen_random_uuid(), 2, 4, now()),
	(gen_random_uuid(), 2, 5, now()),
	(gen_random_uuid(), 2, 6, now()),
	(gen_random_uuid(), 2, 12, now()),
	(gen_random_uuid(), 2, 1, now())
ON CONFLICT (role_id, action_id) DO NOTHING;

UPDATE public.users
SET role_id = 2,
	updated_at = now()
WHERE role_id IS NULL
	AND role = 'user';
