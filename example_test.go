package subheaders_test

import (
	"log"
	"os"

	subheaders "github.com/mdigger/goldmark-subheaders"
	"github.com/yuin/goldmark"
)

func Example() {
	var md = goldmark.New(subheaders.Option)
	var source = []byte(`
## Title
## Subtitle
text

## Subtitle 1

## Subtitle 2

text
`)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// <header>
	// <h2>Title</h2>
	// <p>Subtitle</p>
	// </header>
	// <p>text</p>
	// <h2>Subtitle 1</h2>
	// <h2>Subtitle 2</h2>
	// <p>text</p>
}
