// Under the Apache-2.0 License
package sdoc

import (
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/document"
	"github.com/groboclown/qazaar-testing/rule-engine/sources"
)

// DocumentObject is a localized version of the document schema.
type DocumentObject struct {
	Comments    []string
	Descriptors []document.DocumentDescriptor
	Id          document.Id
	Sources     []sources.Source
}

// Documents simplifies and unifies the representation of documents.
type Documents struct {
	Objects  []DocumentObject
	Refs     map[document.Id]DocumentRef
	Problems *problem.ProblemSet
	sources  *sources.SourceGen
}

// DocumentRef simplifies the document.CommonDocumentSource structure.
type DocumentRef struct {
	// The resource identifier within the 'repo'.  Depending on the repo type, this
	// most likely has a required format for that repository.
	Loc string

	// General repository category containing the source.  This might be 'git' if
	// stored in a Git repository, or 'aws-s3', if stored in an Amazon S3 key store,
	// or 'intranet' if stored in an Intranet source.  The different programs may have
	// their own requirements for this value.  It does not define a location within
	// the repository, though.
	Rep string

	// An identifier to reference the unique version of the source, as dictated by the
	// repository type.  This might a commit id, or document revision, or a date-time
	// stamp.  The repository type may not have the ability to retrieve this version
	// (someone may have deleted it, or the repository does not support versioning).
	Ver *string
}

// New creates a new, shared Documents structure.
func New() *Documents {
	return &Documents{
		Objects:  make([]DocumentObject, 0),
		Refs:     make(map[document.Id]DocumentRef),
		Problems: problem.New(),
		sources:  sources.SourceGenerator(),
	}
}
