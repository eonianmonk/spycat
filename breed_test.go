package spycat

import "testing"

func TestCatBreedValidator(t *testing.T) {

	cv, err := NewCatValidator()
	if err != nil {
		t.Fatalf("Failed to initialize cat validator: %s", err)
	}

	t.Run("valid-breed", func(t *testing.T) {
		err = cv.Validate("Aegean")
		if err != nil {
			t.Fatalf("Unexpected error on existing cat breed: %s", err)
		}
	})

	t.Run("invalid-breed", func(t *testing.T) {
		err = cv.Validate("ฅ^•ﻌ•^ฅ")
		if err == nil {
			t.Fatal("Expected error on unknown breed")
		}
	})
}
