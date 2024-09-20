package cache

import (
	"fmt"
)

/* -------------------------------- cache key ------------------------------- */
func getCacheKey() string {
	return fmt.Sprintf("test")
}
