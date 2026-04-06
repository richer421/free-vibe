package data

import (
	"context"
	"database/sql"
	"free-vibe-coding/internal/conf"
	sqlcgen "free-vibe-coding/internal/data/gen/sqlc"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData)

// Data .
type Data struct {
	db *sql.DB
	q  *sqlcgen.Queries
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	driver := c.GetDatabase().GetDriver()
	if driver == "" {
		driver = "mysql"
	}

	db, err := sql.Open(driver, c.GetDatabase().GetSource())
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	d := &Data{
		db: db,
		q:  sqlcgen.New(db),
	}
	cleanup := func() {
		log.Info("closing the data resources")
		if err := db.Close(); err != nil {
			log.Errorf("failed to close db: %v", err)
		}
	}
	return d, cleanup, nil
}

// Queries returns sqlc query object.
func (d *Data) Queries() *sqlcgen.Queries {
	return d.q
}
