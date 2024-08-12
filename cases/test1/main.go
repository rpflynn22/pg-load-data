package main

import (
	"context"
	"github.com/rpflynn22/pg-load-data/pkg/database"
	"github.com/rpflynn22/pg-load-data/pkg/insertdriver"
	"math/rand"
	"sync"
)

func main() {
	ctx := context.Background()

	db, err := database.GetDB(ctx)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cfg := insertdriver.Config{
		Count:             10_000_000,
		BatchSize:         10_000,
		ConcurrencyFactor: 99,
	}

	if err := insertdriver.Insert(ctx, cfg, db, query, &sourcer{}); err != nil {
		panic(err)
	}
}

const query = `insert into ryanftest1 (id, num, data) values (:id, :num, :data) on conflict do nothing`

type ryanftest1Data struct {
	ID   int    `db:"id"`
	Num  int    `db:"num"`
	Data string `db:"data"`
}

type sourcer struct {
	count int
	mux   sync.Mutex
}

func (s *sourcer) SourceData(ctx context.Context) (ryanftest1Data, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.count++
	return ryanftest1Data{
		ID:   s.count,
		Num:  rand.Intn(100),
		Data: "bean bag",
	}, nil
}
