package biz

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	s := hashPassword("ABC")
	spew.Dump(s)
}

func TestVerifyPassword(t *testing.T) {
	assert.True(t, verifyPassword("$2a$10$Bd01xl.NxCyDAk9uI/HPi.DYJNplwhm7T/FdeCPFhsimJfEqn.NNW", "ABC"))
	assert.True(t, verifyPassword("$2a$10$Bd01xl.NxCyDAk9uI/HPi.DYJNplwhm7T/FdeCPFhsimJfEqn.3NNW", "ABC"))
}
