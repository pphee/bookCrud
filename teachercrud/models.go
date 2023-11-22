package models

import (
	"crypto/sha256"
	"fmt"
)

type Teacher struct {
	HashedID  string `firestore:"hashedId,omitempty"`
	Name      string `firestore:"name"`
	FirstName string `firestore:"firstName"`
	LastName  string `firestore:"lastName"`
	Subject   string `firestore:"subject"`
	Email     string `firestore:"email"`
	Phone     string `firestore:"phone"`
}

func (t *Teacher) HashID() string {
	h := sha256.New()
	h.Write([]byte(t.HashedID))
	return fmt.Sprintf("%x", h.Sum(nil))
}
