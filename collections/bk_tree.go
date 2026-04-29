package collections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"regexp"
	"slices"
)

const (
	MATCH    = 0
	GAP      = 1
	MISMATCH = 1
)

type matrix [][]int

func newMatrix(n, m int) matrix {
	nm := make([][]int, n)
	for i := range nm {
		nm[i] = make([]int, m)
	}
	return nm
}

func score(seq1, seq2 string) int {
	n, m := len(seq1), len(seq2)
	nm := newMatrix(n+1, m+1)
	for i := range seq1 {
		nm[i][0] = GAP * i
	}
	for i := range seq2 {
		nm[0][i] = GAP * i
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			hor := nm[i-1][j] + GAP
			ver := nm[i][j-1] + GAP
			var diag int
			if seq1[i-1] == seq2[j-1] {
				diag = MATCH
			} else {
				diag = MISMATCH
			}
			diag += nm[i-1][j-1]
			nm[i][j] = min(hor, ver, diag)
		}
	}
	return nm[n][m]
}

type String interface {
	~string
}

type bkTree struct {
	term     string
	children map[int]*bkTree
	deleted  bool
}

func (b *bkTree) iter() iter.Seq[string] {
	return func(yield func(string) bool) {
		if !yield(b.term) {
			return
		}
		for _, child := range b.children {
			for el := range child.iter() {
				if !yield(el) {
					return
				}
			}
		}
	}
}

func (b *bkTree) search(term string, fuzziness int) bool {
	d0 := score(b.term, term)
	for _, subTerm := range regexp.MustCompile(`[\w-]`).Split(term, -1) {
		if s := score(b.term, subTerm); s < d0 {
			d0 = s
		}
	}
	if d0 <= fuzziness && !b.deleted {
		return true
	}
	low := d0 - fuzziness
	high := d0 + fuzziness
	for dist, node := range b.children {
		if dist >= low && dist <= high {
			if result := node.search(term, fuzziness); result {
				return true
			}
		}
	}
	return false
}

// BKTree is a Burkhard-Keller tree for approximate string matching.
// It supports fuzzy membership tests using edit distance, controlled by the Fuzziness field.
type BKTree struct {
	root      *bkTree
	Fuzziness int
	len       int
}

// NewBKTree creates a BKTree pre-populated with the given string values.
func NewBKTree(values []string) BKTree {
	tree := BKTree{}
	for _, value := range values {
		tree.Add(value)
	}
	return tree
}

// Into returns all non-deleted terms in the tree as a slice.
func (b *BKTree) Into() []string {
	return slices.Collect(b.Iter())
}

// Add inserts a term into the tree. Duplicate terms are silently ignored.
func (b *BKTree) Add(term string) {
	b.len++
	q := &bkTree{
		term:     term,
		children: make(map[int]*bkTree),
	}
	if b.root == nil {
		b.root = q
		return
	}
	curr := b.root
	for {
		k := score(curr.term, term)
		if k == 0 {
			return
		}
		if child, ok := curr.children[k]; ok {
			curr = child
		} else {
			curr.children[k] = q
			return
		}
	}
}

// Has reports whether term is present in the tree within the configured Fuzziness distance.
func (b *BKTree) Has(term string) bool {
	return b.root.search(term, b.Fuzziness)
}

func (b *bkTree) deleteExact(term string) bool {
	k := score(b.term, term)
	if k == 0 {
		if b.deleted {
			return false
		}
		b.deleted = true
		return true
	}
	child, ok := b.children[k]
	if !ok {
		return false
	}
	return child.deleteExact(term)
}

// Del soft-deletes term from the tree. Has no effect if the term is not present.
func (b *BKTree) Del(term string) {
	b.len--
	if b.root != nil {
		b.root.deleteExact(term)
	}
}

// Len returns the number of terms tracked by the tree.
func (b *BKTree) Len() int {
	return b.len
}

// Iter returns a forward iterator over all non-deleted terms in the tree.
func (b *BKTree) Iter() iter.Seq[string] {
	return func(yield func(string) bool) {
		for element := range b.root.iter() {
			if !yield(element) {
				return
			}
		}
	}
}

// Iter2 returns an indexed iterator over all non-deleted terms in the tree.
func (b *BKTree) Iter2() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		idx := 0
		for element := range b.root.iter() {
			if !yield(idx, element) {
				return
			}
			idx++
		}
	}
}

// UnmarshalJSON populates the tree from a JSON array of strings.
// Fuzziness is set to 2 after decoding.
func (b *BKTree) UnmarshalJSON(bts []byte) error {
	if bytes.Equal(bts, []byte("null")) {
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(bts))

	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '[' {
		return fmt.Errorf("set: expected JSON array")
	}

	if b == nil {
		*b = *new(BKTree)
	}

	for dec.More() {
		var term string
		if err = dec.Decode(&term); err != nil {
			return err
		}
		b.Add(term)
	}

	_, err = dec.Token()
	b.Fuzziness = 2
	return err
}

// MarshalJSON serializes the tree as a JSON array of strings.
// Returns "null" if the tree is empty.
func (b *BKTree) MarshalJSON() ([]byte, error) {
	if b.root == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer
	buf.WriteByte('[')

	first := true
	for v := range b.Iter() {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		s, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.Write(s)
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}
