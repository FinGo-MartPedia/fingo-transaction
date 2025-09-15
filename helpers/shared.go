package helpers

import (
	"fmt"
	"time"
)

func GenerateReference() string {
	now := time.Now()
	nowFormat := now.Format("20060102150405")
	reference := fmt.Sprintf("WT-%s", nowFormat)
	return reference
}
