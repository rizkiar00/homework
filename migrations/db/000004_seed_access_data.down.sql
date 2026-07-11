DELETE FROM public.role_accesses
WHERE role_id IN (1, 2)
	AND action_id BETWEEN 1 AND 12;

DELETE FROM public.roles
WHERE role_id IN (1, 2);

DELETE FROM public.actions
WHERE action_id BETWEEN 1 AND 12;
