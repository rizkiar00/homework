package constant

const (
	CodeSuccess            = "SUCCESS"
	CodeBadRequest         = "ERR_BAD_REQUEST"
	CodeUnauthorized       = "ERR_UNAUTHORIZED"
	CodeForbidden          = "ERR_FORBIDDEN"
	CodeNotFound           = "ERR_NOT_FOUND"
	CodeConflict           = "ERR_CONFLICT"
	CodePayloadTooLarge    = "ERR_PAYLOAD_TOO_LARGE"
	CodeTooManyRequests    = "ERR_TOO_MANY_REQUESTS"
	CodeInternalServer     = "ERR_INTERNAL_SERVER"
	CodeServiceUnavailable = "ERR_SERVICE_UNAVAILABLE"
)

const (
	MessageSuccess            = "success"
	MessageCreated            = "created"
	MessageDeleted            = "deleted"
	MessagePayloadTooLarge    = "payload too large"
	MessageTooManyRequests    = "too many requests"
	MessageInternalServer     = "internal server error"
	MessageServiceUnavailable = "service unavailable"
)
