package repomongo

import (
	models "book/studentcrud"
	"book/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type IStudentRepository interface {
	Create(ctx context.Context, student *models.Student) (*mongo.InsertOneResult, error)
	FindByID(ctx context.Context, studentID string) (models.Student, error)
	FindAll(ctx context.Context) ([]*models.Student, error)
	Update(ctx context.Context, studentID string, updateData interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, studentID string) (*mongo.DeleteResult, error)
}

type mongoStudentRepository struct {
	collection        Collection
	encryptionService util.IEncryptionService
}

func NewStudentRepository(collection Collection, encryptionService util.IEncryptionService) IStudentRepository {
	return &mongoStudentRepository{
		collection:        collection,
		encryptionService: encryptionService,
	}
}

func (r *mongoStudentRepository) Create(ctx context.Context, student *models.Student) (*mongo.InsertOneResult, error) {
	studentFirstName, err := r.encryptionService.Encrypt(student.FirstName)
	if err != nil {
		return nil, err
	}
	student.FirstName = studentFirstName

	studentLastName, err := r.encryptionService.Encrypt(student.LastName)
	if err != nil {
		return nil, err
	}
	student.LastName = studentLastName

	result, err := r.collection.InsertOne(ctx, student)
	return result, err
}

func (r *mongoStudentRepository) FindByID(ctx context.Context, studentID string) (models.Student, error) {
	var student models.Student
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return student, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&student)
	if err != nil {
		return student, err
	}

	student.FirstName, err = r.encryptionService.Decrypt(student.FirstName)
	if err != nil {
		return student, err
	}

	student.LastName, err = r.encryptionService.Decrypt(student.LastName)
	if err != nil {
		return student, err
	}
	return student, err
}

func (r *mongoStudentRepository) FindAll(ctx context.Context) ([]*models.Student, error) {
	var students []*models.Student
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var student models.Student
		if err := cursor.Decode(&student); err != nil {
			return nil, err
		}
		student.FirstName, _ = r.encryptionService.Decrypt(student.FirstName)
		student.LastName, _ = r.encryptionService.Decrypt(student.LastName)
		students = append(students, &student)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *mongoStudentRepository) Update(ctx context.Context, studentID string, updateData interface{}) (*mongo.UpdateResult, error) {
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData}, options.Update().SetUpsert(true))
	return result, err
}

func (r *mongoStudentRepository) Delete(ctx context.Context, studentID string) (*mongo.DeleteResult, error) {
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return result, err
}