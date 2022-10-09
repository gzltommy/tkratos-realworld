package auth

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token := GenerateToken("abc", 1000)
	spew.Dump(token)
}
