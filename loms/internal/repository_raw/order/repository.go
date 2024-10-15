package order

import (
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/pgcluster"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_raw/outbox"
)

type Repository struct {
	cluster  *pgcluster.Cluster
	notifier *outbox.Repository
}

func NewRepository(cluster *pgcluster.Cluster) *Repository {
	return &Repository{
		cluster: cluster,
	}
}
