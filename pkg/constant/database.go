package constant

const (
	MessageDatabaseNotConfigured     = "database is not configured"
	MessageDatabaseConfigIsAvailable = "database configuration is present"
	MessageRedisNotConfigured        = "redis is not configured"
	MessageRedisConfigIsAvailable    = "redis configuration is present"
)

const (
	TableUsers              = "public.users"
	TableTestTable          = "public.test_table"
	TableActions            = "public.actions"
	TableRoles              = "public.roles"
	TableRoleAccesses       = "public.role_accesses"
	TableEmailVerifications = "public.user_email_verifications"
	TablePasswordResets     = "public.user_password_resets"
)

const (
	ColumnUserID   = "user_id"
	ColumnUsername = "username"
	ColumnEmail    = "email"
	ColumnTestID   = "test_id"
	ColumnDescTest = "desc_test"
)
