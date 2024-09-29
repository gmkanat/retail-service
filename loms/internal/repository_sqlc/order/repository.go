package order

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/pgcluster"
	"time"
)

type Repository struct {
	cluster *pgcluster.Cluster
}

func NewRepository(cluster *pgcluster.Cluster) *Repository {
	return &Repository{
		cluster: cluster,
	}
}

func currentTimestamp() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}
}
