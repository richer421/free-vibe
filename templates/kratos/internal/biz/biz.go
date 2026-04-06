package biz

import (
	greeterbiz "free-vibe-coding/internal/biz/greeter"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(greeterbiz.NewUsecase)
