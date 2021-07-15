package entity

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "Not found"
}

func (e *NotFoundError) HTTPCode() int {
	return 404
}
