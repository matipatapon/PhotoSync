package helper_test

import (
	"fmt"
	"photosync/src/helper"
	"testing"
)

func TestHasherShouldHashFile(t *testing.T) {
	expectedHash := "e0ac3601005dfa1864f5392aabaf7d898b1b5bab854f1acb4491bcd806b76b0c"
	file := []byte("file content")
	sut := helper.NewHasher()

	hash, err := sut.Hash(file)

	if hash != expectedHash {
		fmt.Printf("Expected '%s', got '%s'", expectedHash, hash)
		t.FailNow()
	}
	if err != nil {
		fmt.Print(err.Error())
		t.FailNow()
	}
}
