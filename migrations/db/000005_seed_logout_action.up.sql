INSERT INTO public.actions (action_id, action_desc, action_type, endpoint)
VALUES
	(12, 'Logout current user', 'POST', '/auth/logout')
ON CONFLICT (action_type, endpoint) DO UPDATE
SET action_desc = EXCLUDED.action_desc;

SELECT setval(pg_get_serial_sequence('public.actions', 'action_id'), (SELECT MAX(action_id) FROM public.actions));

INSERT INTO public.role_accesses (role_access_id, role_id, action_id, created_at)
VALUES
	(gen_random_uuid(), 1, 12, now()),
	(gen_random_uuid(), 2, 12, now())
ON CONFLICT (role_id, action_id) DO NOTHING;
