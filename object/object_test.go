package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	hello1Key := testHashKey(t, hello1)
	hello2Key := testHashKey(t, hello2)
	diff1Key := testHashKey(t, diff1)
	diff2Key := testHashKey(t, diff2)

	if hello1Key != hello2Key {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1Key != diff2Key {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1Key == diff1Key {
		t.Errorf("strings with different content have same hash keys")
	}
}

func testHashKey(t *testing.T, object Hashable) HashKey {
	t.Helper()

	key, err := object.HashKey()
	if err != nil {
		t.Fatalf("HashKey() returned error: %v", err)
	}

	return key
}
