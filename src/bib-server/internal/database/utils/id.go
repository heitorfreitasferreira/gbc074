package utils

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidId = errors.New("identificador inválido")
)

func ConvertUserBookIds(userIdStr, bookIdStr string) (int64, int64, error) {
	userId, err := strconv.Atoi(strings.TrimSpace(userIdStr))
	if err != nil {
		return 0, 0, ErrInvalidId
	}
	bookId, err := strconv.Atoi(strings.TrimSpace(bookIdStr))
	if err != nil {
		return 0, 0, ErrInvalidId
	}

	return int64(userId), int64(bookId), nil
}
