package greeter

import "github.com/google/wire"

// ProviderSet is greeter data providers.
var ProviderSet = wire.NewSet(NewRepo)
