package pkg

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache is a global cache
var Cache = cache.New(5*time.Minute, 10*time.Minute)
