package tests

import (
	"github.com/stretchr/testify/suite"
	grpcSuite "gitlab.ozon.dev/kanat_9999/homework/loms/e2e/suite"
	"testing"
)

func TestGRPCSuite(t *testing.T) {
	suite.Run(t, new(grpcSuite.GRPCSuite))
}
