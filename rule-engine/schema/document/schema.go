// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package document

import "encoding/json"
import "fmt"

// Document author comment text.
type Comment string

// List of document author comments.
type CommentList []Comment

// A shared primary document reference.  Source locations can refer to this
// document through the identifier, but should also include an anchor.
type CommonDocumentSource struct {
	// Comment corresponds to the JSON schema field "$comment".
	Comment *Comment `json:"$comment,omitempty" yaml:"$comment,omitempty" mapstructure:"$comment,omitempty"`

	// Comments corresponds to the JSON schema field "$comments".
	Comments CommentList `json:"$comments,omitempty" yaml:"$comments,omitempty" mapstructure:"$comments,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id Id `json:"id" yaml:"id" mapstructure:"id"`

	// The resource identifier within the 'repo'.  Depending on the repo type, this
	// most likely has a required format for that repository.
	Loc string `json:"loc" yaml:"loc" mapstructure:"loc"`

	// General repository category containing the source.  This might be 'git' if
	// stored in a Git repository, or 'aws-s3', if stored in an Amazon S3 key store,
	// or 'intranet' if stored in an Intranet source.  The different programs may have
	// their own requirements for this value.  It does not define a location within
	// the repository, though.
	Rep string `json:"rep" yaml:"rep" mapstructure:"rep"`

	// An identifier to reference the unique version of the source, as dictated by the
	// repository type.  This might a commit id, or document revision, or a date-time
	// stamp.  The repository type may not have the ability to retrieve this version
	// (someone may have deleted it, or the repository does not support versioning).
	Ver *string `json:"ver,omitempty" yaml:"ver,omitempty" mapstructure:"ver,omitempty"`
}

// Pool of document source references, which may be referenced from the source
// locations.
type CommonDocumentSourceList []CommonDocumentSource

// UnmarshalJSON implements json.Unmarshaler.
func (j *CommonDocumentSource) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in CommonDocumentSource: required")
	}
	if _, ok := raw["loc"]; raw != nil && !ok {
		return fmt.Errorf("field loc in CommonDocumentSource: required")
	}
	if _, ok := raw["rep"]; raw != nil && !ok {
		return fmt.Errorf("field rep in CommonDocumentSource: required")
	}
	type Plain CommonDocumentSource
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if len(plain.Loc) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "loc", 1)
	}
	if len(plain.Loc) > 8000 {
		return fmt.Errorf("field %s length: must be <= %d", "loc", 8000)
	}
	if len(plain.Rep) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "rep", 1)
	}
	if len(plain.Rep) > 200 {
		return fmt.Errorf("field %s length: must be <= %d", "rep", 200)
	}
	if plain.Ver != nil && len(*plain.Ver) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "ver", 1)
	}
	if plain.Ver != nil && len(*plain.Ver) > 200 {
		return fmt.Errorf("field %s length: must be <= %d", "ver", 200)
	}
	*j = CommonDocumentSource(plain)
	return nil
}

// Unique identifier for the descriptor. (Taken from the ontology schema)
type DescriptorKey string

// A numeric value.
type DescriptorNumericValue float64

// A textual value, either an enumerated or free value.
type DescriptorTextValue string

type DescriptorValue interface{}

// Details about source objects in terms of the ontology.  STATUS: ready for review
type DocumentDescriptionV1SchemaJson struct {
	// Comment corresponds to the JSON schema field "$comment".
	Comment *Comment `json:"$comment,omitempty" yaml:"$comment,omitempty" mapstructure:"$comment,omitempty"`

	// Comments corresponds to the JSON schema field "$comments".
	Comments CommentList `json:"$comments,omitempty" yaml:"$comments,omitempty" mapstructure:"$comments,omitempty"`

	// Schema corresponds to the JSON schema field "$schema".
	Schema Schema `json:"$schema" yaml:"$schema" mapstructure:"$schema"`

	// CommonSourceRefs corresponds to the JSON schema field "commonSourceRefs".
	CommonSourceRefs CommonDocumentSourceList `json:"commonSourceRefs,omitempty" yaml:"commonSourceRefs,omitempty" mapstructure:"commonSourceRefs,omitempty"`

	// List of document objects.
	Objects []DocumentObject `json:"objects" yaml:"objects" mapstructure:"objects"`

	// Sources corresponds to the JSON schema field "sources".
	Sources DocumentSources `json:"sources,omitempty" yaml:"sources,omitempty" mapstructure:"sources,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DocumentDescriptionV1SchemaJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["$schema"]; raw != nil && !ok {
		return fmt.Errorf("field $schema in DocumentDescriptionV1SchemaJson: required")
	}
	if _, ok := raw["objects"]; raw != nil && !ok {
		return fmt.Errorf("field objects in DocumentDescriptionV1SchemaJson: required")
	}
	type Plain DocumentDescriptionV1SchemaJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DocumentDescriptionV1SchemaJson(plain)
	return nil
}

// The ontological descriptor key and associated value.
type DocumentDescriptor struct {
	// Key corresponds to the JSON schema field "key".
	Key DescriptorKey `json:"key" yaml:"key" mapstructure:"key"`

	// The descriptor values.  Each entry must conform to the descriptor key's value
	// types.
	Values []DocumentDescriptorValuesElem `json:"values" yaml:"values" mapstructure:"values"`
}

type DocumentDescriptorValuesElem interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DocumentDescriptor) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["key"]; raw != nil && !ok {
		return fmt.Errorf("field key in DocumentDescriptor: required")
	}
	if _, ok := raw["values"]; raw != nil && !ok {
		return fmt.Errorf("field values in DocumentDescriptor: required")
	}
	type Plain DocumentDescriptor
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DocumentDescriptor(plain)
	return nil
}

// Description for a document object; a unique identifier, location information,
// and ontology descriptors.
type DocumentObject struct {
	// Comment corresponds to the JSON schema field "$comment".
	Comment *Comment `json:"$comment,omitempty" yaml:"$comment,omitempty" mapstructure:"$comment,omitempty"`

	// Comments corresponds to the JSON schema field "$comments".
	Comments CommentList `json:"$comments,omitempty" yaml:"$comments,omitempty" mapstructure:"$comments,omitempty"`

	// List of descriptors assigned to this object.  The list of descriptors must not
	// contain multiple items with the same key, and the values for the key must
	// conform to the document's ontology.
	Descriptors []DocumentDescriptor `json:"descriptors,omitempty" yaml:"descriptors,omitempty" mapstructure:"descriptors,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id Id `json:"id" yaml:"id" mapstructure:"id"`

	// Sources corresponds to the JSON schema field "sources".
	Sources DocumentSources `json:"sources" yaml:"sources" mapstructure:"sources"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DocumentObject) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in DocumentObject: required")
	}
	if _, ok := raw["sources"]; raw != nil && !ok {
		return fmt.Errorf("field sources in DocumentObject: required")
	}
	type Plain DocumentObject
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DocumentObject(plain)
	return nil
}

// Sources that contained the original definitions.  A tool collected those
// descriptions into this document.
type DocumentSources []SourceLocation

// Unique identifying string for the item.  These should be ASCII alpha-numeric +
// simple separators.
type Id string

// Data exchange schema format.
type Schema string

// Pointer to the location of the source.  Due to the prevalence of this object,
// property names use a truncated form to shrink file sizes.  The 'ref' points to a
// common document source identifier in the commonSourceRefs list.
type SourceLocation struct {
	// Comment corresponds to the JSON schema field "$comment".
	Comment *Comment `json:"$comment,omitempty" yaml:"$comment,omitempty" mapstructure:"$comment,omitempty"`

	// Comments corresponds to the JSON schema field "$comments".
	Comments CommentList `json:"$comments,omitempty" yaml:"$comments,omitempty" mapstructure:"$comments,omitempty"`

	// A location within the source.  This depends upon the source type; it might be
	// an HTML anchor tag, or a paragraph title, or a function name, or a line number,
	// or an opcode index.
	A *string `json:"a,omitempty" yaml:"a,omitempty" mapstructure:"a,omitempty"`

	// Ref corresponds to the JSON schema field "ref".
	Ref Id `json:"ref" yaml:"ref" mapstructure:"ref"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SourceLocation) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["ref"]; raw != nil && !ok {
		return fmt.Errorf("field ref in SourceLocation: required")
	}
	type Plain SourceLocation
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.A != nil && len(*plain.A) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "a", 1)
	}
	if plain.A != nil && len(*plain.A) > 4000 {
		return fmt.Errorf("field %s length: must be <= %d", "a", 4000)
	}
	*j = SourceLocation(plain)
	return nil
}
