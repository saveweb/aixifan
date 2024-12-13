package utils

import "testing"

func Test_CheckIAItemExist(t *testing.T) {
	want := true
	exist, err := CheckIAItemExist("BiliBili-BV14o4y1T7JZ_p1-ZJ1T6VB")
	if err != nil {
		t.Fatal(err)
	}
	if exist != want {
		t.Fatalf("want %v, got %v", want, exist)
	}

	want = false
	exist, err = CheckIAItemExist("BiliBili-BV14o4y1T7JZ_p1-ZJ1T6VB_not_exist")
	if err != nil {
		t.Fatal(err)
	}
	if exist != want {
		t.Fatalf("want %v, got %v", want, exist)
	}
}
