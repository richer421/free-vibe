package greeter

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"

	greeterbiz "free-vibe-coding/internal/biz/greeter"
	"free-vibe-coding/internal/data"
	sqlcgen "free-vibe-coding/internal/data/gen/sqlc"

	"github.com/go-kratos/kratos/v2/log"
)

type repo struct {
	data *data.Data
	q    *sqlcgen.Queries
	log  *log.Helper
}

// NewRepo creates greeter repo.
func NewRepo(d *data.Data, logger log.Logger) greeterbiz.Repo {
	return &repo{
		data: d,
		q:    d.Queries(),
		log:  log.NewHelper(logger),
	}
}

func (r *repo) Save(ctx context.Context, g *greeterbiz.Entity) (*greeterbiz.Entity, error) {
	if g == nil {
		return nil, fmt.Errorf("greeter is nil")
	}
	res, err := r.q.CreateGreeter(ctx, g.Hello)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	row, err := r.q.GetGreeterByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toBizGreeter(row), nil
}

func (r *repo) Update(ctx context.Context, g *greeterbiz.Entity) (*greeterbiz.Entity, error) {
	if g == nil {
		return nil, fmt.Errorf("greeter is nil")
	}
	if g.ID == 0 {
		return nil, fmt.Errorf("greeter id is required")
	}
	err := r.q.UpdateGreeter(ctx, sqlcgen.UpdateGreeterParams{
		Hello: g.Hello,
		ID:    g.ID,
	})
	if err != nil {
		return nil, err
	}
	row, err := r.q.GetGreeterByID(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	return toBizGreeter(row), nil
}

func (r *repo) FindByID(ctx context.Context, id int64) (*greeterbiz.Entity, error) {
	row, err := r.q.GetGreeterByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, sql.ErrNoRows) {
			return nil, greeterbiz.ErrUserNotFound
		}
		return nil, err
	}
	return toBizGreeter(row), nil
}

func (r *repo) ListByHello(ctx context.Context, hello string) ([]*greeterbiz.Entity, error) {
	rows, err := r.q.ListGreetersByHello(ctx, hello)
	if err != nil {
		return nil, err
	}
	return toBizGreeterList(rows), nil
}

func (r *repo) ListAll(ctx context.Context) ([]*greeterbiz.Entity, error) {
	rows, err := r.q.ListGreeters(ctx)
	if err != nil {
		return nil, err
	}
	return toBizGreeterList(rows), nil
}

func toBizGreeter(row sqlcgen.Greeter) *greeterbiz.Entity {
	return &greeterbiz.Entity{
		ID:        row.ID,
		Hello:     row.Hello,
		CreatedAt: row.CreatedAt,
	}
}

func toBizGreeterList(rows []sqlcgen.Greeter) []*greeterbiz.Entity {
	list := make([]*greeterbiz.Entity, 0, len(rows))
	for _, row := range rows {
		list = append(list, toBizGreeter(row))
	}
	return list
}
