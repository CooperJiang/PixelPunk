package share

import (
	"strings"

	"github.com/google/uuid"
)

func generateID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
