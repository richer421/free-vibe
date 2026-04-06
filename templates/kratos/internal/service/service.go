package service

import (
	greetersvc "free-vibe-coding/internal/service/greeter"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(greetersvc.NewService)
