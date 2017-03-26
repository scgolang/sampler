package sampler_test

import (
	"testing"

	"github.com/scgolang/sampler"
)

func TestAdd(t *testing.T) {
}

func newTestSampler(t *testing.T) *Sampler {
	s, err := sampler.New("")
	if err != nil {
		t.Fatal(err)
	}
	return s
}
