package tests

import (
	"github.com/stretchr/testify/suite"
	testSuite "gitlab.ozon.dev/kanat_9999/homework/cart/e2e/suite"
	"testing"
)

func TestListItemSuite(t *testing.T) {
	suite.Run(t, new(testSuite.ListItemSuite))
}

func TestDeleteItemSuite(t *testing.T) {
	suite.Run(t, new(testSuite.DeleteItemSuite))
}
