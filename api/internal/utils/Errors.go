package utils

type UnauthorizedError struct {
	Message string
}

type InternalServerError struct {
	Message string
}

type ServiceError struct {
	Message string
}

type RepositoryError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

func (err UnauthorizedError) Error() string {
	return err.Message
}

func (err InternalServerError) Error() string {
	return err.Message
}

func (err ServiceError) Error() string {
	return err.Message
}

func (err RepositoryError) Error() string {
	return err.Message
}

func (err BadRequestError) Error() string {
	return err.Message
}
