package bimap

import "testing"

func assertEqual[T comparable](t *testing.T, want, got T, msg string) {
	if want != got {
		t.Errorf("want '%v', got '%v': %s", want, got, msg)
	}
}

func TestReadFromZero(t *testing.T) {
	var m Bimap[int, string]

	assertEqual(t, 0, m.Len(), "get length")
	assertEqual(t, false, m.ContainsForward(1), "contains forward")
	assertEqual(t, false, m.ContainsReverse("foo"), "contains reverse")
	m.Each(func(key int, value string) {
		t.Errorf("loop each function called: key %q, value %q", key, value)
	})

	_, ok := m.GetForward(1)
	assertEqual(t, false, ok, "get forward")
	_, ok = m.GetReverse("foo")
	assertEqual(t, false, ok, "get reverse")
}

func TestWrite(t *testing.T) {
	var m Bimap[int, string]
	assertEqual(t, false, m.ContainsForward(1), "contains before add?")
	m.Add(1, "foo")
	assertEqual(t, true, m.ContainsForward(1), "contains after add?")
	m.RemoveForward(1)
	assertEqual(t, false, m.ContainsForward(1), "contains after remove?")
}
