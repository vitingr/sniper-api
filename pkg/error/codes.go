package apierror

const (
	CodeInternalServer = "internal_server_error"
	CodeBadRequest     = "bad_request"
	CodeNotFound       = "not_found"
	CodeUnauthorized   = "unauthorized"
	CodeForbidden      = "forbidden"
	CodeConflict       = "conflict"

	CodeInvalidFormat = "invalid_format"
	CodeRequiredField = "required_field"

	CodeUserAlreadyExists = "user_already_exists"
	CodePageNotFound      = "page_not_found"
)
