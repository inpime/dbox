package dbox

import (
    "strings"
    "github.com/satori/go.uuid"
)

func NewUUID() string {
    return strings.Replace(uuid.NewV4().String(), "-", "", -1) 
}