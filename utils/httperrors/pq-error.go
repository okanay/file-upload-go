package httperrors

import (
	"errors"
	"github.com/lib/pq"
	"net/http"
)

const (
	CodeUniqueViolation     = "23505"
	CodeForeignKeyViolation = "23503"
)

func handlePQError(err error) Response {
	var pqErr *pq.Error
	errors.As(err, &pqErr)

	switch pqErr.Code {
	case CodeUniqueViolation:
		return Response{
			Message: "A record with this data already exists",
			Status:  http.StatusConflict,
		}
	case CodeForeignKeyViolation:
		return Response{
			Message: "The referenced record does not exist",
			Status:  http.StatusNotFound,
		}
	default:
		return Response{
			Message: "An unexpected error occurred: " + pqErr.Message,
			Status:  http.StatusInternalServerError,
		}
	}
}

func isPQError(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr)
}
