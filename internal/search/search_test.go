package search

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSearch(t *testing.T) {
	g := NewWithT(t)

	expected := "github.com/onsi/ginkgo/..."
	result := SearchPackages( "ginkgo")

	g.Expect(result[0].Repo).Should(Equal(expected))
}
