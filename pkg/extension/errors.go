package extension

type HttpError struct {
	Message string
}

func (e *HttpError) Error() string { return e.Message }

type AuthError struct {
	HttpError
}

type NotFoundError struct {
	HttpError
}

type BadRequestError struct {
	HttpError
}
