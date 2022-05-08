package random

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	u := uuid.New()
	return strings.ReplaceAll(u.String(), `-`, ``)
}

func Int() int {
	return rand.Int()
}
