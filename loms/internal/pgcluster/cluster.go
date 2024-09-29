package pgcluster

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cluster struct {
	write         *pgxpool.Pool
	read          []*pgxpool.Pool
	roundRobinIdx int
	mu            *sync.Mutex
}

func New() *Cluster {
	return &Cluster{
		roundRobinIdx: 0,
		mu:            &sync.Mutex{},
	}
}

func (c *Cluster) SetWriter(master *pgxpool.Pool) *Cluster {
	c.write = master
	return c
}

func (c *Cluster) AddReader(readPools ...*pgxpool.Pool) *Cluster {
	c.read = append(c.read, readPools...)
	return c
}

func (c *Cluster) GetWriter(ctx context.Context) (*pgxpool.Pool, error) {
	if c.write == nil {
		return nil, errors.New("writer pool is not set")
	}

	if err := c.write.Ping(ctx); err != nil {
		return nil, err
	}
	return c.write, nil
}

func (c *Cluster) GetReader(ctx context.Context) (*pgxpool.Pool, error) {
	if len(c.read) == 0 {
		return nil, errors.New("reader pools are not set")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	reader := c.read[c.roundRobinIdx]
	c.roundRobinIdx = (c.roundRobinIdx + 1) % len(c.read)

	if err := reader.Ping(ctx); err != nil {
		return nil, err
	}
	return reader, nil
}

func (c *Cluster) Close() {
	if c.write != nil {
		c.write.Close()
	}
	for _, pool := range c.read {
		if pool != nil {
			pool.Close()
		}
	}
}
