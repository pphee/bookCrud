package models

type Teacher struct {
	ID        string // Firestore document ID
	FirstName string `firestore:"firstName"`
	LastName  string `firestore:"lastName"`
	Subject   string `firestore:"subject"`
	Email     string `firestore:"email"`
	Phone     string `firestore:"phone"`
}
