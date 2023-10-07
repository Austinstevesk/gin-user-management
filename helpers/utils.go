package helpers

import (
	"fmt"
	"strconv"
)


func ConvertStringToInt(str string) (int, error) {
	intVal, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("Could not convert string to integer %w", err)
	}
	return intVal, nil
}