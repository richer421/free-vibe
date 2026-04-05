package greeter

import (
	"context"
	"time"

	v1 "free-vibe-coding/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Entity is a Greeter model.
type Entity struct {
	ID        int64
	Hello     string
	CreatedAt time.Time
}

// Repo is a Greater repo.
type Repo interface {
	Save(context.Context, *Entity) (*Entity, error)
	Update(context.Context, *Entity) (*Entity, error)
	FindByID(context.Context, int64) (*Entity, error)
	ListByHello(context.Context, string) ([]*Entity, error)
	ListAll(context.Context) ([]*Entity, error)
}

// Usecase is a Greeter usecase.
type Usecase struct {
	repo Repo
}

// NewUsecase new a Greeter usecase.
func NewUsecase(repo Repo) *Usecase {
	return &Usecase{repo: repo}
}

// Create creates a Greeter, and returns the new Greeter.
func (uc *Usecase) Create(ctx context.Context, g *Entity) (*Entity, error) {
	log.Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
