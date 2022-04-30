package random

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	u := uuid.New()
	return strings.ReplaceAll(u.String(), `-`, ``)
}
