package utils

import (
	"context"
	"errors"

	"employee-management/internal/models"
)

func AsAppError(err error) *models.AppError {
	if err == nil {
		return nil
	}

	var appErr *models.AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return models.ErrTimeout(err)
	}
	if errors.Is(err, context.Canceled) {
		return models.ErrTimeout(err)
	}

	return models.ErrInternal(err)
}
