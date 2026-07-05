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
	(11, 'Assign user role', 'PUT', '/users/{user_id}/role')
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
	('00000000-0000-0000-0000-000000000001', 1, 1, now()),
	('00000000-0000-0000-0000-000000000002', 1, 2, now()),
	('00000000-0000-0000-0000-000000000003', 1, 3, now()),
	('00000000-0000-0000-0000-000000000004', 1, 4, now()),
	('00000000-0000-0000-0000-000000000005', 1, 5, now()),
	('00000000-0000-0000-0000-000000000006', 1, 6, now()),
	('00000000-0000-0000-0000-000000000007', 1, 7, now()),
	('00000000-0000-0000-0000-000000000008', 1, 8, now()),
	('00000000-0000-0000-0000-000000000009', 1, 9, now()),
	('00000000-0000-0000-0000-000000000010', 1, 10, now()),
	('00000000-0000-0000-0000-000000000011', 1, 11, now()),
	('00000000-0000-0000-0000-000000000012', 2, 1, now())
ON CONFLICT (role_id, action_id) DO NOTHING;
