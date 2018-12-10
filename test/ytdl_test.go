package test

import (
	"testing"

	"gitlab.com/ttpcodes/prismriver/internal/app/sources"
)

func TestYtdl(t *testing.T) {
	_, err := sources.GetVideo("https://www.youtube.com/watch?v=Ys2p_bXOaAc")
	if err != nil {
		t.Error(err)
	}
}
