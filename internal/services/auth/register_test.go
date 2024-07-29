package auth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

type RegisterSuite struct {
	suite.Suite
}
