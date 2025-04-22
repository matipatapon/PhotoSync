package password

import (
	"testing"
)

func generate_password(length int) string {
	var password string
	for i := int(1); i <= length; i++ {
		password += string(rune(49 + i%10))
	}
	return password
}

func TestPasswordFacadeTwoSamePasswordsShouldNotHaveSameHash(t *testing.T) {
	sut := PasswordFacade{}
	password := generate_password(72)
	hash1, err1 := sut.HashPassword(password)
	if err1 != nil {
		t.Fail()
	}
	hash2, err2 := sut.HashPassword(password)
	if err2 != nil {
		t.Fail()
	}
	if hash1 == hash2 {
		t.Error("Two password shouldn't have same hash!")
	}
}

func TestPasswordFacadeShouldFailWhenPasswordIsTooLong(t *testing.T) {
	sut := PasswordFacade{}
	too_long_password := generate_password(73)
	hash, err := sut.HashPassword(too_long_password)
	if hash != "" || err == nil {
		t.Fail()
	}
}
