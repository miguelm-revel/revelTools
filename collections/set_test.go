package collections

import (
	"encoding/json"
	"testing"
)

func TestSet_UnmarshalJSON_Null(t *testing.T) {
	var s Set[int] // nil map
	if err := json.Unmarshal([]byte("null"), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Fatalf("expected nil set after unmarshalling null, got %#v", s)
	}
}

func TestSet_MarshalJSON_NilPointer(t *testing.T) {
	var ps *Set[int] = nil
	b, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(b) != "null" {
		t.Fatalf("expected %q, got %q", "null", string(b))
	}
}

func TestSet_MarshalJSON_NilMapValue(t *testing.T) {
	var s Set[int] = nil
	b, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(b) != "null" {
		t.Fatalf("expected %q, got %q", "null", string(b))
	}
}

func TestSet_UnmarshalJSON_Array_Ints(t *testing.T) {
	var s Set[int]
	if err := json.Unmarshal([]byte(`[1,2,3]`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatalf("expected non-nil set")
	}
	if s.Len() != 3 {
		t.Fatalf("expected len=3, got %d", s.Len())
	}
	for _, v := range []int{1, 2, 3} {
		if !s.Has(v) {
			t.Fatalf("expected set to contain %d", v)
		}
	}
}

func TestSet_UnmarshalJSON_Deduplicates(t *testing.T) {
	var s Set[int]
	if err := json.Unmarshal([]byte(`[1,1,2,2,2,3]`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 3 {
		t.Fatalf("expected len=3 after dedup, got %d", s.Len())
	}
	for _, v := range []int{1, 2, 3} {
		if !s.Has(v) {
			t.Fatalf("expected set to contain %d", v)
		}
	}
}

func TestSet_UnmarshalJSON_ExpectedArrayError(t *testing.T) {
	var s Set[int]
	err := json.Unmarshal([]byte(`{"a":1}`), &s)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// Mensaje exacto de tu implementación
	if err.Error() != "Set: expected JSON array" {
		t.Fatalf("expected error %q, got %q", "Set: expected JSON array", err.Error())
	}
}

func TestSet_UnmarshalJSON_ClearsExistingSet(t *testing.T) {
	s := NewSet([]int{9, 10, 11})
	if s.Len() != 3 {
		t.Fatalf("precondition failed: expected len=3, got %d", s.Len())
	}

	// Unmarshal debe limpiar el contenido anterior.
	if err := json.Unmarshal([]byte(`[1,2]`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Len() != 2 {
		t.Fatalf("expected len=2, got %d", s.Len())
	}
	if s.Has(9) || s.Has(10) || s.Has(11) {
		t.Fatalf("expected old elements to be cleared, got %#v", s)
	}
	if !s.Has(1) || !s.Has(2) {
		t.Fatalf("expected new elements to exist")
	}
}

func TestSet_MarshalJSON_Array_Strings(t *testing.T) {
	s := NewSet([]string{"a", "b", "c"})

	b, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No asumimos orden. Decodificamos el JSON a slice y verificamos contenido.
	var got []string
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("marshal output is not valid JSON array: %v; bytes=%s", err, string(b))
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 elements in marshalled array, got %d; bytes=%s", len(got), string(b))
	}

	seen := map[string]bool{}
	for _, v := range got {
		seen[v] = true
	}

	for _, v := range []string{"a", "b", "c"} {
		if !seen[v] {
			t.Fatalf("expected marshalled JSON to contain %q; got=%v; bytes=%s", v, got, string(b))
		}
	}
}

func TestSet_JSON_RoundTrip_Ints(t *testing.T) {
	orig := NewSet([]int{1, 2, 3, 4, 5})

	b, err := json.Marshal(&orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded Set[int]
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unexpected error: %v; bytes=%s", err, string(b))
	}

	if decoded.Len() != orig.Len() {
		t.Fatalf("expected len=%d, got %d", orig.Len(), decoded.Len())
	}

	for it := range orig.Iter() {
		if !decoded.Has(it) {
			t.Fatalf("expected decoded to contain %v", it)
		}
	}
}
