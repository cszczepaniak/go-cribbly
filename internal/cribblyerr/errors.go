package cribblyerr

import (
	"fmt"
)

type notFound struct {
	resource string
}

func (e notFound) Error() string {
	return fmt.Sprintf(`object %v not found`, e.resource)
}

func ErrNotFound(resource string) error {
	return notFound{
		resource: resource,
	}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(notFound)
	return ok
}
