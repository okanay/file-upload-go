package httperrors

import (
	"errors"
	"fmt"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func Handle(err error) Response {
	fmt.Print(err)

	pqErr := isPQError(err)
	if pqErr == true {
		return handlePQError(err)
	}

	var httpErr *HttpError
	ok := errors.As(err, &httpErr)
	if ok {
		return Response{
			Message: httpErr.Message,
			Status:  httpErr.Code,
		}
	}

	return Response{
		Message: "Internal Server Error",
		Status:  500,
	}
}
