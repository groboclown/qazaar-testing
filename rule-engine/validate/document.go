// Under the Apache-2.0 License
package validate

import (
	"context"
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func ValidateDocuments(doc *sdoc.Documents, ont *sont.AllowedDescriptors, probs problem.Adder) {
	if doc != nil {
		for _, d := range doc.Objects {
			ValidateDocument(d, ont, probs)
		}
	}
}

func ValidateDocument(doc *sdoc.DocumentObject, ont *sont.AllowedDescriptors, probs problem.Adder) {
	if doc != nil {
		for _, d := range doc.Descriptors {
			ValidateDescriptor(d, ont, doc.Sources, probs)
		}
	}
}

// ValidateDocumentsAsync validates all the documents, and returns a channel that reads once (and then closes) when they complete.
func ValidateDocumentsAsync(
	doc *sdoc.Documents,
	ont *sont.AllowedDescriptors,
	probs problem.Adder,
	ctx context.Context,
) <-chan bool {
	ret := make(chan bool)

	go func() {
		defer func() {
			ret <- true
			close(ret)
		}()

		var wg sync.WaitGroup

		if doc != nil {
			for _, d := range doc.Objects {
				if ctx.Err() != nil {
					break
				}
				if d != nil {
					for _, desc := range d.Descriptors {
						if ctx.Err() != nil {
							break
						}
						wg.Add(1)
						go func() {
							ValidateDescriptor(desc, ont, d.Sources, probs)
						}()
					}
				}
			}
		}

		wg.Wait()
	}()

	return ret
}
