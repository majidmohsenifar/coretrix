package handlers

import "net/http"

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func getValidationError(err error) (ErrorResponse, int) {
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: err.Error(),
	}, http.StatusBadRequest
}
