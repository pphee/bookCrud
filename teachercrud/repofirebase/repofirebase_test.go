package repofirebase_test

import (
	models "book/teachercrud"
	"book/teachercrud/repofirebase"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strconv"
	"testing"
)

func setupFirestoreTest(t *testing.T) (context.Context, *firestore.Client, repofirebase.ITeacherRepository) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:7060")
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "test-project")
	assert.NoError(t, err, "Firestore client should be initialized without error")

	collectionRef := client.Collection("teachers")
	repo := repofirebase.NewTeacherRepository(collectionRef) // Assuming this returns an ITeacherRepository

	return ctx, client, repo
}

func TestCreateTeacher(t *testing.T) {
	ctx, client, repo := setupFirestoreTest(t)
	defer client.Close()

	testTeacher := &models.Teacher{
		Name:      "Test Name",
		FirstName: "TestFirstName",
		LastName:  "TestLastName",
		Subject:   "TestSubject",
		Email:     "test@example.com",
		Phone:     "1234567890",
	}
	testTeacher.HashedID = testTeacher.HashID() // Generate the HashedID

	err := repo.CreateTeacher(ctx, testTeacher)
	assert.NoError(t, err)
}

func TestGetTeacherByID(t *testing.T) {
	ctx, client, repo := setupFirestoreTest(t)
	defer client.Close()

	collectionRef := client.Collection("teachers")

	testTeacher := &models.Teacher{
		Name:      "Test Name",
		FirstName: "TestFirstName",
		LastName:  "TestLastName",
		Subject:   "TestSubject",
		Email:     "test@example.com",
		Phone:     "1234567890",
	}
	testTeacher.HashedID = testTeacher.HashID() // Generate the HashedID

	_, err := collectionRef.Doc(testTeacher.HashedID).Set(ctx, testTeacher)
	assert.NoError(t, err)

	retrievedTeacher, err := repo.GetTeacherByID(ctx, testTeacher.HashedID)
	assert.NoError(t, err, "GetTeacherByID should not produce an error")
	assert.NotNil(t, retrievedTeacher, "Retrieved teacher should not be nil")
	assert.Equal(t, testTeacher.HashedID, retrievedTeacher.HashedID, "Retrieved teacher HashedID should match")
	assert.Equal(t, testTeacher.Name, retrievedTeacher.Name, "Retrieved teacher Name should match")
	assert.Equal(t, testTeacher.FirstName, retrievedTeacher.FirstName, "Retrieved teacher FirstName should match")
	assert.Equal(t, testTeacher.LastName, retrievedTeacher.LastName, "Retrieved teacher LastName should match")
	assert.Equal(t, testTeacher.Subject, retrievedTeacher.Subject, "Retrieved teacher Subject should match")
	assert.Equal(t, testTeacher.Email, retrievedTeacher.Email, "Retrieved teacher Email should match")
	assert.Equal(t, testTeacher.Phone, retrievedTeacher.Phone, "Retrieved teacher Phone should match")
}

func TestGetAllTeachers(t *testing.T) {
	ctx, client, repo := setupFirestoreTest(t)
	defer client.Close()

	collectionRef := client.Collection("teachers")

	iter := collectionRef.Documents(ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		_, err = doc.Ref.Delete(ctx)
		assert.NoError(t, err, "Error clearing collection")
	}

	repo = repofirebase.NewTeacherRepository(collectionRef)

	for i := 0; i < 3; i++ {
		testTeacher := &models.Teacher{
			Name:      "Test Name" + strconv.Itoa(i),
			FirstName: "TestFirstName" + strconv.Itoa(i),
			LastName:  "TestLastName" + strconv.Itoa(i),
			Subject:   "TestSubject" + strconv.Itoa(i),
			Email:     "test" + strconv.Itoa(i) + "@example.com",
			Phone:     "1234567890",
		}
		testTeacher.HashedID = testTeacher.HashID()
		_, err := collectionRef.Doc(testTeacher.HashedID).Set(ctx, testTeacher)
		assert.NoError(t, err)
	}

	teachers, err := repo.GetAllTeachers(ctx)
	assert.NoError(t, err, "GetAllTeachers should not produce an error")
	assert.Len(t, teachers, 3, "Should retrieve 3 teachers")

}

func TestUpdateTeacher(t *testing.T) {
	ctx, client, repo := setupFirestoreTest(t)
	defer client.Close()

	collectionRef := client.Collection("teachers")

	// Create a test teacher
	testTeacher := &models.Teacher{
		Name:      "pee",
		FirstName: "phee",
		LastName:  "dans",
		Subject:   "math",
		Email:     "pee@dans.com",
		Phone:     "1234567890",
	}
	testTeacher.HashedID = testTeacher.HashID() // Generate HashedID
	_, err := collectionRef.Doc(testTeacher.HashedID).Set(ctx, testTeacher)
	assert.NoError(t, err, "Test teacher should be created without error")

	// Update the test teacher
	updatedTeacher := &models.Teacher{
		Name:      "jack",
		FirstName: "queen",
		LastName:  "king",
		Subject:   "computer",
		Email:     "jack@gmail.com",
		Phone:     "0987654321",
	}
	err = repo.UpdateTeacher(ctx, testTeacher.HashedID, updatedTeacher)
	assert.NoError(t, err, "Updating a teacher should not produce an error")

	// Retrieve the updated teacher to verify changes
	doc, err := collectionRef.Doc(testTeacher.HashedID).Get(ctx)
	assert.NoError(t, err, "Retrieving updated teacher should not produce an error")

	var retrievedTeacher models.Teacher
	err = doc.DataTo(&retrievedTeacher)
	assert.NoError(t, err, "DataTo should not produce an error")

	// Validate the updates
	assert.Equal(t, updatedTeacher.Name, retrievedTeacher.Name, "Teacher name should be updated")
	assert.Equal(t, updatedTeacher.FirstName, retrievedTeacher.FirstName, "Teacher first name should be updated")
	assert.Equal(t, updatedTeacher.LastName, retrievedTeacher.LastName, "Teacher last name should be updated")
	assert.Equal(t, updatedTeacher.Subject, retrievedTeacher.Subject, "Teacher subject should be updated")
	assert.Equal(t, updatedTeacher.Email, retrievedTeacher.Email, "Teacher email should be updated")
	assert.Equal(t, updatedTeacher.Phone, retrievedTeacher.Phone, "Teacher phone should be updated")
}

func TestDeleteTeacher(t *testing.T) {
	ctx, client, repo := setupFirestoreTest(t)
	defer client.Close()

	collectionRef := client.Collection("teachers")

	testTeacher := &models.Teacher{
		Name:      "Test Name",
		FirstName: "TestFirstName",
		LastName:  "TestLastName",
		Subject:   "TestSubject",
		Email:     "test@example.com",
		Phone:     "1234567890",
	}
	testTeacher.HashedID = testTeacher.HashID() // Generate HashedID
	_, err := collectionRef.Doc(testTeacher.HashedID).Set(ctx, testTeacher)
	assert.NoError(t, err, "Test teacher should be created without error")

	err = repo.DeleteTeacher(ctx, testTeacher.HashedID)
	assert.NoError(t, err, "Deleting a teacher should not produce an error")

	_, err = collectionRef.Doc(testTeacher.HashedID).Get(ctx)

	assert.NotNil(t, err, "Retrieving a deleted teacher should produce an error")
	Status, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, Status.Code(), "Error code should be NotFound for a deleted teacher")
}
