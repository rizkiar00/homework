package constant

const (
	RoleUser                  = "user"
	TokenTypeBearer           = "Bearer"
	AuthorizationBearerPrefix = "Bearer "
	ContextUserClaimsKey      = "user_claims"
)

const (
	MessageUnauthorized              = "unauthorized"
	MessageForbidden                 = "forbidden"
	MessageInvalidRequestBody        = "invalid request body"
	MessageInvalidRegisterCredential = "username is required and password minimum length is 8"
	MessageDataNotFound              = "data not found"
	MessageInvalidAuthCredential     = "invalid username or password"
	MessageUsernameAlreadyExists     = "username already exists"
)
