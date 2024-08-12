package insertdriver

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/constraints"
	"golang.org/x/sync/errgroup"
)

// Config contains configuration for the driver that loads data.
type Config struct {
	// Count specifies how many records we want to insert.
	Count int
	// BatchSize specifies how many records to insert in each query.
	// For best performance, push it right up against the 65k param
	// limit.
	BatchSize int
	// ConcurrencyFactor specifies how many goroutines to run concurrently
	// inserting data. By default, postgres has a max connection limit of 100,
	// so going above that isn't recommended.
	ConcurrencyFactor int
}

// DataSourcer is used to return data to insert. The SourceData method must be
// goroutine safe as it will be called concurrently if ConcurrencyFactor > 1.
type DataSourcer[T any] interface {
	SourceData(context.Context) (T, error)
}

// Insert spawns goroutines that insert data returned by the data sourcer.
func Insert[T any](
	ctx context.Context,
	cfg Config,
	db *sqlx.DB,
	query string,
	dataSourcer DataSourcer[T],
) error {
	grp, grpCtx := errgroup.WithContext(ctx)
	var countInserted int
	countIncrMux := sync.Mutex{}

	for i := 0; i < cfg.ConcurrencyFactor; i++ {

		grp.Go(func() error {

			ctx := grpCtx

			for {
				countIncrMux.Lock()
				numToInsert := min(cfg.BatchSize, max(0, cfg.Count-countInserted))
				countInserted += numToInsert
				countIncrMux.Unlock()

				if numToInsert <= 0 {
					return nil
				}

				dataSlc := make([]T, numToInsert)
				for i := 0; i < numToInsert; i++ {
					data, err := dataSourcer.SourceData(ctx)
					if err != nil {
						return fmt.Errorf("source failed: %w", err)
					}
					dataSlc[i] = data
				}

				_, err := db.NamedExecContext(ctx, query, dataSlc)
				if err != nil {
					return fmt.Errorf("write failed: %w", err)
				}
			}
		})
	}

	return grp.Wait()
}

func min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}
