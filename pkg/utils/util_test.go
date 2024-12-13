package utils

import "testing"

func Tets_ToIdentifier(t *testing.T) {
	want := "AcFun-123_p1"
	got := ToIdentifier("123")
	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}

	want = "AcFun-123_p2"
	got = ToIdentifier("123", "2")
	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}
}
