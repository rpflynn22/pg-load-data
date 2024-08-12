package main

import (
	"context"
	"github.com/rpflynn22/pg-load-data/pkg/database"
	"github.com/rpflynn22/pg-load-data/pkg/insertdriver"
	"math/rand"

	"github.com/gofrs/uuid"
)

const query = `insert into skipscantest (
	account_id,
	sequence_id,
	data
) values (
	:account_id,
	:sequence_id,
	:data
)`

func main() {
	ctx := context.Background()
	db, err := database.GetDB(ctx)
	if err != nil {
		panic(err)
	}

	cfg := insertdriver.Config{
		Count:             10_000_000,
		BatchSize:         1_000,
		ConcurrencyFactor: 99,
	}

	s := sourcer{}
	s.init(10_000)
	if err := insertdriver.Insert(ctx, cfg, db, query, &s); err != nil {
		panic(err)
	}
}

type skipscandata struct {
	AccountID  string `db:"account_id"`
	SequenceID int    `db:"sequence_id"`
	Data       int    `db:"data"`
}

type sourcer struct {
	acctIDs []string
}

func (s *sourcer) init(numAccts int) {
	for i := 0; i < numAccts; i++ {
		s.acctIDs = append(s.acctIDs, uuid.Must(uuid.NewV4()).String())
	}
}

func (s *sourcer) SourceData(ctx context.Context) (skipscandata, error) {
	acctID := s.acctIDs[rand.Intn(len(s.acctIDs))]
	return skipscandata{
		AccountID:  acctID,
		SequenceID: rand.Intn(1000),
		Data:       rand.Intn(1000000),
	}, nil
}
