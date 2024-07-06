// Under the Apache-2.0 License
package sources

// Source contains a general representation of the data-exchange source.
type Source struct {
	ref *innerRef
	a   *string
}

func (s Source) Loc() string {
	return s.ref.loc
}

func (s Source) Rep() string {
	return s.ref.rep
}

func (s Source) Ver() *string {
	return s.ref.ver
}

func (s Source) A() *string {
	return s.a
}

type SourceGen struct {
	refs map[string]*innerRef
}

func SourceGenerator() *SourceGen {
	return &SourceGen{refs: make(map[string]*innerRef)}
}

type innerRef struct {
	// The resource identifier within the 'repo'.  Depending on the repo type, this
	// most likely has a required format for that repository.
	loc string

	// General repository category containing the source.  This might be 'git' if
	// stored in a Git repository, or 'aws-s3', if stored in an Amazon S3 key store,
	// or 'intranet' if stored in an Intranet source.  The different programs may have
	// their own requirements for this value.  It does not define a location within
	// the repository, though.
	rep string

	// An identifier to reference the unique version of the source, as dictated by the
	// repository type.  This might a commit id, or document revision, or a date-time
	// stamp.  The repository type may not have the ability to retrieve this version
	// (someone may have deleted it, or the repository does not support versioning).
	ver *string
}
