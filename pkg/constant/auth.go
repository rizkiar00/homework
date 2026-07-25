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
	MessageInvalidRegisterCredential = "full_name, username, email are required and password minimum length is 8"
	MessageDataNotFound              = "data not found"
	MessageInvalidAuthCredential     = "invalid username or password"
	MessageEmailNotVerified          = "email is not verified"
	MessageUsernameAlreadyExists     = "username already exists"
	MessageEmailAlreadyExists        = "email already exists"
	MessageInvalidVerificationCode   = "invalid or expired verification code"
	MessageInvalidResetPasswordCode  = "invalid or expired reset password code"
	MessageInvalidResetPassword      = "email, code, and new password minimum length is 8 are required"
	MessageVerificationSubmitted     = "registration submitted, please verify your email"
	MessageVerificationSent          = "verification code sent"
	MessageEmailVerified             = "email verified"
	MessageResetPasswordCodeSent     = "if the email is registered, reset password code has been sent"
	MessagePasswordReset             = "password has been reset"
	MessageEmailRequestTooFrequent   = "please wait before requesting another email code"
	MessageEmailRequestLimitExceeded = "too many email code requests, please try again later"
)
