package dbox

import (
    "fmt"
)

var (
	ErrNotFound = fmt.Errorf("not_found")
    ErrInvalidData = fmt.Errorf("invalid_data")

    ErrEmptyName = fmt.Errorf("empty_name")
    ErrEmptyID = fmt.Errorf("empty_id")
)