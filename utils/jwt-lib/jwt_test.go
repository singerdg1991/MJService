package jwtlib

import "testing"

func TestEncrypt(t *testing.T) {
	token, err := Encrypt(JwtPayload{
		ID:       1,
		Username: "TestUser",
		Data:     "Test",
	})
	if err != nil {
		t.Errorf("Can't return jwt because: %s", err)
	}
	if len(token) == 0 {
		t.Error("Token is invalid")
	}
}
