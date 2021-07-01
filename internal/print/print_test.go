package print

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/alex-held/gofile/internal/gofile"
)

func TestPrettyPrint(t *testing.T) {
	type tc struct {
		Expected string
		Input    PrettyPrintable
	}

	tcs := map[string]tc{
		"Gofile": {
			Expected: `NAME        URI                             
ginkgo      github.com/onsi/ginkgo          
gomega      github.com/onsi/gomega          
godepgraph  github.com/okisielk/godepgraph  
`,
			Input: gofile.GoFile{Installables: []gofile.Installable{
				{Name: "ginkgo", URI: "github.com/onsi/ginkgo"},
				{Name: "gomega", URI: "github.com/onsi/gomega"},
				{Name: "godepgraph", URI: "github.com/okisielk/godepgraph"},
			}},
		},
	}

	for name, test := range tcs {
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			_ = os.Chdir("/Users/dev/go/src/github.com/alex-held/gofile")
			actual, err := PrettyPrint(test.Input)
			println(actual)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(actual).Should(Equal(test.Expected))
		})
	}
}
