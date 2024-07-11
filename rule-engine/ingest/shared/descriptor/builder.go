// Under the Apache-2.0 License
package descriptor

import "strings"

type DescriptorValueBuilder[T DescriptorValueTypes] interface {
	AddList(n []T)
	Add(n DescriptorValueBuilder[T])
	List() []T
	Count() int
	Has(n T) bool
	Copy() DescriptorValueBuilder[T]
}

func NewNumericBuilder(distinct bool) DescriptorValueBuilder[float64] {
	if distinct {
		r := make(map[float64]bool)
		return (*distinctNumberValue)(&r)
	}
	r := make([]float64, 0)
	return (*duplicateNumberValue)(&r)
}

func NewTextBuilder(distinct bool, caseSensitive bool) DescriptorValueBuilder[string] {
	if distinct {
		r := make(map[string]bool)
		if caseSensitive {
			return (*distinctCaseSensitiveValue)(&r)
		}
		return (*distinctCaseInsensitiveValue)(&r)
	}
	r := make([]string, 0)
	if caseSensitive {
		return (*duplicateCaseSensitiveValue)(&r)
	}
	return (*duplicateCaseInsensitiveValue)(&r)
}

// ---------------------

type duplicateNumberValue []float64

func (d *duplicateNumberValue) Add(n DescriptorValueBuilder[float64]) {
	switch a := n.(type) {
	case *duplicateNumberValue:
		d.AddList(*a)
	case *distinctNumberValue:
		*d = appendExactMap(*d, *a)
	default:
		d.AddList(a.List())
	}
}
func (d *duplicateNumberValue) AddList(n []float64) {
	(*d) = append((*d), n...)
}
func (d *duplicateNumberValue) List() []float64 {
	return (*d)
}
func (d *duplicateNumberValue) Count() int {
	return len(*d)
}
func (d *duplicateNumberValue) Has(n float64) bool {
	return hasExactList(*d, n)
}
func (d *duplicateNumberValue) Copy() DescriptorValueBuilder[float64] {
	v := make([]float64, len(*d))
	copy(v, *d)
	return (*duplicateNumberValue)(&v)
}

// ---------------------

type distinctNumberValue map[float64]bool

func (d *distinctNumberValue) Add(n DescriptorValueBuilder[float64]) {
	switch a := n.(type) {
	case *duplicateNumberValue:
		d.AddList(*a)
	case *distinctNumberValue:
		addExactMap(*d, *a)
	default:
		d.AddList(n.List())
	}
}
func (d *distinctNumberValue) AddList(n []float64) {
	addExactList(*d, n)
}
func (d *distinctNumberValue) List() []float64 {
	return DistinctMapArray(*d)
}
func (d *distinctNumberValue) Count() int {
	return len(*d)
}
func (d *distinctNumberValue) Has(n float64) bool {
	return hasExactMap(*d, n)
}
func (d *distinctNumberValue) Copy() DescriptorValueBuilder[float64] {
	v := make(map[float64]bool)
	addExactMap(v, *d)
	return (*distinctNumberValue)(&v)
}

// ---------------------

type duplicateCaseSensitiveValue []string

func (d *duplicateCaseSensitiveValue) Add(n DescriptorValueBuilder[string]) {
	switch a := n.(type) {
	case *duplicateCaseSensitiveValue:
		d.AddList(*a)
	case *duplicateCaseInsensitiveValue:
		d.AddList(*a)
	case *distinctCaseInsensitiveValue:
		*d = appendExactMap(*d, *a)
	case *distinctCaseSensitiveValue:
		*d = appendExactMap(*d, *a)
	default:
		(*d) = append((*d), a.List()...)
	}
}
func (d *duplicateCaseSensitiveValue) AddList(n []string) {
	*d = append(*d, n...)
}
func (d *duplicateCaseSensitiveValue) List() []string {
	return *d
}
func (d *duplicateCaseSensitiveValue) Count() int {
	return len(*d)
}
func (d *duplicateCaseSensitiveValue) Has(n string) bool {
	return hasExactList(*d, n)
}
func (d *duplicateCaseSensitiveValue) Copy() DescriptorValueBuilder[string] {
	v := make([]string, len(*d))
	copy(v, *d)
	return (*duplicateCaseSensitiveValue)(&v)
}

// ---------------------

type duplicateCaseInsensitiveValue []string

func (d *duplicateCaseInsensitiveValue) Add(n DescriptorValueBuilder[string]) {
	switch a := n.(type) {
	case *duplicateCaseSensitiveValue:
		d.AddList(*a)
	case *duplicateCaseInsensitiveValue:
		*d = append(*d, (*a)...)
	case *distinctCaseInsensitiveValue:
		*d = appendExactMap(*d, *a)
	case *distinctCaseSensitiveValue:
		*d = appendSensitiveMap(*d, *a)
	default:
		(*d) = append((*d), a.List()...)
	}
}
func (d *duplicateCaseInsensitiveValue) AddList(n []string) {
	*d = appendSensitiveList(*d, n)
}
func (d *duplicateCaseInsensitiveValue) List() []string {
	return *d
}
func (d *duplicateCaseInsensitiveValue) Count() int {
	return len(*d)
}
func (d *duplicateCaseInsensitiveValue) Has(n string) bool {
	return hasSensitiveList(*d, n)
}
func (d *duplicateCaseInsensitiveValue) Copy() DescriptorValueBuilder[string] {
	v := make([]string, len(*d))
	copy(v, *d)
	return (*duplicateCaseInsensitiveValue)(&v)
}

// ---------------------

type distinctCaseInsensitiveValue map[string]bool

func (d *distinctCaseInsensitiveValue) Add(n DescriptorValueBuilder[string]) {
	switch a := n.(type) {
	case *duplicateCaseSensitiveValue:
		d.AddList(*a)
	case *duplicateCaseInsensitiveValue:
		addSensitiveList(*d, *a)
	case *distinctCaseInsensitiveValue:
		addExactMap(*d, *a)
	case *distinctCaseSensitiveValue:
		addSensitiveMap(*d, *a)
	default:
		d.AddList(a.List())
	}
}
func (d *distinctCaseInsensitiveValue) AddList(n []string) {
	addSensitiveList(*d, n)
}
func (d *distinctCaseInsensitiveValue) List() []string {
	return DistinctMapArray(*d)
}
func (d *distinctCaseInsensitiveValue) Count() int {
	return len(*d)
}
func (d *distinctCaseInsensitiveValue) Has(n string) bool {
	return hasSensitiveMap(*d, n)
}
func (d *distinctCaseInsensitiveValue) Copy() DescriptorValueBuilder[string] {
	v := make(map[string]bool)
	addExactMap(v, *d)
	return (*distinctCaseInsensitiveValue)(&v)
}

// ---------------------

type distinctCaseSensitiveValue map[string]bool

func (d *distinctCaseSensitiveValue) Add(n DescriptorValueBuilder[string]) {
	switch a := n.(type) {
	case *duplicateCaseSensitiveValue:
		d.AddList(*a)
	case *duplicateCaseInsensitiveValue:
		d.AddList(*a)
	case *distinctCaseInsensitiveValue:
		addExactMap(*d, *a)
	case *distinctCaseSensitiveValue:
		addExactMap(*d, *a)
	default:
		d.AddList(a.List())
	}
}
func (d *distinctCaseSensitiveValue) AddList(n []string) {
	addExactList(*d, n)
}
func (d *distinctCaseSensitiveValue) List() []string {
	return DistinctMapArray(*d)
}
func (d *distinctCaseSensitiveValue) Count() int {
	return len(*d)
}
func (d *distinctCaseSensitiveValue) Has(n string) bool {
	return hasExactMap(*d, n)
}
func (d *distinctCaseSensitiveValue) Copy() DescriptorValueBuilder[string] {
	v := make(map[string]bool)
	addExactMap(v, *d)
	return (*distinctCaseSensitiveValue)(&v)
}

// ---------------------

func appendExactMap[T DescriptorValueTypes](d []T, a map[T]bool) []T {
	i := len(d)
	tl := make([]T, i+len(a))
	copy(tl, d)
	for k := range a {
		tl[i] = k
		i++
	}
	return tl
}

func appendSensitiveMap(d []string, a map[string]bool) []string {
	i := len(d)
	tl := make([]string, i+len(a))
	copy(tl, d)
	for k := range a {
		tl[i] = strings.ToLower(k)
		i++
	}
	return tl
}

func appendSensitiveList(d []string, a []string) []string {
	i := len(d)
	tl := make([]string, i+len(a))
	copy(tl, d)
	for _, v := range a {
		tl[i] = strings.ToLower(v)
		i++
	}
	return tl
}

func addExactList[T DescriptorValueTypes](d map[T]bool, a []T) {
	for _, v := range a {
		d[v] = true
	}
}

func addExactMap[T DescriptorValueTypes](d map[T]bool, a map[T]bool) {
	for k := range a {
		d[k] = true
	}
}

func addSensitiveMap(d map[string]bool, a map[string]bool) {
	for k := range a {
		d[strings.ToLower(k)] = true
	}
}

func addSensitiveList(d map[string]bool, a []string) {
	for _, v := range a {
		d[strings.ToLower(v)] = true
	}
}

func hasExactList[T DescriptorValueTypes](d []T, e T) bool {
	for _, v := range d {
		if v == e {
			return true
		}
	}
	return false
}

func hasSensitiveList(d []string, e string) bool {
	return hasExactList(d, strings.ToLower(e))
}

func hasExactMap[T DescriptorValueTypes](d map[T]bool, e T) bool {
	_, ok := d[e]
	return ok
}

func hasSensitiveMap(d map[string]bool, e string) bool {
	return hasExactMap(d, strings.ToLower(e))
}
