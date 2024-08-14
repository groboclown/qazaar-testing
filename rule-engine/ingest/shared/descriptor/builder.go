// Under the Apache-2.0 License
package descriptor

import "strings"

type ImmutableDescriptorValue[T DescriptorValueTypes] interface {
	List() []T
	Count() int
	Has(n T) bool
	IsDistinct() bool
	IsCaseSensitive() bool
	Copy() DescriptorValueBuilder[T]
}

type DescriptorValueBuilder[T DescriptorValueTypes] interface {
	AddList(n []T)
	Add(n DescriptorValueBuilder[T])
	RemoveOnce(n T)
	RemoveAll(n T)
	Clear()
	List() []T
	Count() int
	Has(n T) bool
	IsDistinct() bool
	IsCaseSensitive() bool
	Copy() DescriptorValueBuilder[T]
	Seal() ImmutableDescriptorValue[T]
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
func (d *duplicateNumberValue) RemoveOnce(n float64) {
	removeExactListOnce(*d, n)
}
func (d *duplicateNumberValue) RemoveAll(n float64) {
	removeExactListAll(*d, n)
}
func (d *duplicateNumberValue) Clear() {
	(*d) = make([]float64, 0)
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
func (d *duplicateNumberValue) IsDistinct() bool {
	return false
}
func (d *duplicateNumberValue) IsCaseSensitive() bool {
	return false
}
func (d *duplicateNumberValue) Seal() ImmutableDescriptorValue[float64] {
	return d.Copy()
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
func (d *distinctNumberValue) RemoveOnce(n float64) {
	removeExactMap(*d, n)
}
func (d *distinctNumberValue) RemoveAll(n float64) {
	removeExactMap(*d, n)
}
func (d *distinctNumberValue) Clear() {
	(*d) = make(map[float64]bool)
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
func (d *distinctNumberValue) IsDistinct() bool {
	return true
}
func (d *distinctNumberValue) IsCaseSensitive() bool {
	return false
}
func (d *distinctNumberValue) Seal() ImmutableDescriptorValue[float64] {
	return d.Copy()
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
func (d *duplicateCaseSensitiveValue) RemoveOnce(n string) {
	removeExactListOnce(*d, n)
}
func (d *duplicateCaseSensitiveValue) RemoveAll(n string) {
	removeExactListAll(*d, n)
}
func (d *duplicateCaseSensitiveValue) Clear() {
	(*d) = make([]string, 0)
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
func (d *duplicateCaseSensitiveValue) IsDistinct() bool {
	return false
}
func (d *duplicateCaseSensitiveValue) IsCaseSensitive() bool {
	return true
}
func (d *duplicateCaseSensitiveValue) Seal() ImmutableDescriptorValue[string] {
	return d.Copy()
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
		*d = appendInsensitiveMap(*d, *a)
	case *distinctCaseSensitiveValue:
		*d = appendExactMap(*d, *a)
	default:
		(*d) = append((*d), a.List()...)
	}
}
func (d *duplicateCaseInsensitiveValue) AddList(n []string) {
	*d = appendInsensitiveList(*d, n)
}
func (d *duplicateCaseInsensitiveValue) RemoveOnce(n string) {
	removeInsensitiveListOnce(*d, n)
}
func (d *duplicateCaseInsensitiveValue) RemoveAll(n string) {
	removeInsensitiveListAll(*d, n)
}
func (d *duplicateCaseInsensitiveValue) Clear() {
	(*d) = make([]string, 0)
}
func (d *duplicateCaseInsensitiveValue) List() []string {
	return *d
}
func (d *duplicateCaseInsensitiveValue) Count() int {
	return len(*d)
}
func (d *duplicateCaseInsensitiveValue) Has(n string) bool {
	return hasInsensitiveList(*d, n)
}
func (d *duplicateCaseInsensitiveValue) IsDistinct() bool {
	return false
}
func (d *duplicateCaseInsensitiveValue) IsCaseSensitive() bool {
	return false
}
func (d *duplicateCaseInsensitiveValue) Seal() ImmutableDescriptorValue[string] {
	return d.Copy()
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
		addInsensitiveList(*d, *a)
	case *distinctCaseInsensitiveValue:
		addExactMap(*d, *a)
	case *distinctCaseSensitiveValue:
		addInsensitiveMap(*d, *a)
	default:
		d.AddList(a.List())
	}
}
func (d *distinctCaseInsensitiveValue) AddList(n []string) {
	addInsensitiveList(*d, n)
}
func (d *distinctCaseInsensitiveValue) RemoveOnce(n string) {
	removeInsensitiveMap(*d, n)
}
func (d *distinctCaseInsensitiveValue) RemoveAll(n string) {
	removeInsensitiveMap(*d, n)
}
func (d *distinctCaseInsensitiveValue) Clear() {
	(*d) = make(map[string]bool)
}
func (d *distinctCaseInsensitiveValue) List() []string {
	return DistinctMapArray(*d)
}
func (d *distinctCaseInsensitiveValue) Count() int {
	return len(*d)
}
func (d *distinctCaseInsensitiveValue) Has(n string) bool {
	return hasInsensitiveMap(*d, n)
}
func (d *distinctCaseInsensitiveValue) IsDistinct() bool {
	return true
}
func (d *distinctCaseInsensitiveValue) IsCaseSensitive() bool {
	return false
}
func (d *distinctCaseInsensitiveValue) Seal() ImmutableDescriptorValue[string] {
	return d.Copy()
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
func (d *distinctCaseSensitiveValue) RemoveOnce(n string) {
	removeExactMap(*d, n)
}
func (d *distinctCaseSensitiveValue) RemoveAll(n string) {
	removeExactMap(*d, n)
}
func (d *distinctCaseSensitiveValue) Clear() {
	(*d) = make(map[string]bool)
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
func (d *distinctCaseSensitiveValue) IsDistinct() bool {
	return true
}
func (d *distinctCaseSensitiveValue) IsCaseSensitive() bool {
	return true
}
func (d *distinctCaseSensitiveValue) Seal() ImmutableDescriptorValue[string] {
	return d.Copy()
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

func appendInsensitiveMap(d []string, a map[string]bool) []string {
	i := len(d)
	tl := make([]string, i+len(a))
	copy(tl, d)
	for k := range a {
		tl[i] = strings.ToLower(k)
		i++
	}
	return tl
}

func appendInsensitiveList(d []string, a []string) []string {
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

func addInsensitiveMap(d map[string]bool, a map[string]bool) {
	for k := range a {
		d[strings.ToLower(k)] = true
	}
}

func addInsensitiveList(d map[string]bool, a []string) {
	for _, v := range a {
		d[strings.ToLower(v)] = true
	}
}

func removeExactListOnce[T DescriptorValueTypes](d []T, v T) []T {
	n := len(d) - 1
	if n < 0 {
		return d
	}
	if d[0] == v {
		return d[1:]
	}
	if d[n] == v {
		return d[:n]
	}
	for i := 1; i < n; i++ {
		if d[i] == v {
			return append(d[:i], d[i+1:]...)
		}
	}
	return d
}

func removeExactListAll[T DescriptorValueTypes](d []T, v T) []T {
	ret := make([]T, 0, len(d))
	for _, i := range d {
		if i != v {
			ret = append(ret, i)
		}
	}
	return ret
}

func removeInsensitiveListOnce(d []string, v string) []string {
	return removeExactListOnce(d, strings.ToLower(v))
}

func removeInsensitiveListAll(d []string, v string) []string {
	return removeExactListAll(d, strings.ToLower(v))
}

func removeExactMap[T DescriptorValueTypes](d map[T]bool, v T) {
	delete(d, v)
}

func removeInsensitiveMap(d map[string]bool, v string) {
	delete(d, strings.ToLower(v))
}

func hasExactList[T DescriptorValueTypes](d []T, e T) bool {
	for _, v := range d {
		if v == e {
			return true
		}
	}
	return false
}

func hasInsensitiveList(d []string, e string) bool {
	return hasExactList(d, strings.ToLower(e))
}

func hasExactMap[T DescriptorValueTypes](d map[T]bool, e T) bool {
	_, ok := d[e]
	return ok
}

func hasInsensitiveMap(d map[string]bool, e string) bool {
	return hasExactMap(d, strings.ToLower(e))
}
