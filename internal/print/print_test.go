package main

import (
	"testing"

	. "github.com/onsi/gomega"
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
			Input: GoFile{Installables: []Installable{
				{Name: "ginkgo", URI: "github.com/onsi/ginkgo"},
				{Name: "gomega", URI: "github.com/onsi/gomega"},
				{Name: "godepgraph", URI: "github.com/okisielk/godepgraph"},
			}},
		},
	}

	for name, test := range tcs {
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(PrettyPrint(test.Input)).Should(Equal(test.Expected))
		})
	}
}
