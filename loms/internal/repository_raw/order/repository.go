package order

import "gitlab.ozon.dev/kanat_9999/homework/loms/internal/pgcluster"

type Repository struct {
	cluster *pgcluster.Cluster
}

func NewRepository(cluster *pgcluster.Cluster) *Repository {
	return &Repository{
		cluster: cluster,
	}
}
