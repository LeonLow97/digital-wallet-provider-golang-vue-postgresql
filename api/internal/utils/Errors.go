package utils

type UnauthorizedError struct {
	Message string
}

type ServiceError struct {
	Message string
}

type RepositoryError struct {
	Message string
}

func (err UnauthorizedError) Error() string {
	return err.Message
}

func (err ServiceError) Error() string {
	return err.Message
}

func (err RepositoryError) Error() string {
	return err.Message
}
