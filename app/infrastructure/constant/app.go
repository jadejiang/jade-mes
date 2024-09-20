package constant

import (
	"time"
)

const (
	// NameSpace ...
	NameSpace = "jade-mes"

	// LogCategoryAccess ...
	LogCategoryAccess = "access_log"
	// LogCategoryError ...
	LogCategoryError = "error"
	// LogCategoryWarn ...
	LogCategoryWarn = "warn"

	// MaxDBRetry ...
	MaxDBRetry = 3

	MaxAttempts = 3

	// CacheTTL ...
	CacheTTL = 12 * time.Hour
)
