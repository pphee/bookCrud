package usecases

import (
	"book/schoolcrud/proto"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

// MockSchoolRepository is a mock of ISchoolRepository interface
type MockSchoolRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSchoolRepositoryMockRecorder
}

type MockSchoolRepositoryMockRecorder struct {
	mock *MockSchoolRepository
}

func NewMockSchoolRepository(ctrl *gomock.Controller) *MockSchoolRepository {
	mock := &MockSchoolRepository{ctrl: ctrl}
	mock.recorder = &MockSchoolRepositoryMockRecorder{mock}
	return mock
}

func (m *MockSchoolRepository) EXPECT() *MockSchoolRepositoryMockRecorder {
	return m.recorder
}

func (m *MockSchoolRepositoryMockRecorder) CreateSchool(ctx interface{}, school interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "CreateSchool", reflect.TypeOf((*MockSchoolRepository)(nil).CreateSchool), ctx, school)
}

func (m *MockSchoolRepositoryMockRecorder) ListSchools(ctx interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "ListSchools", reflect.TypeOf((*MockSchoolRepository)(nil).ListSchools), ctx)
}

func (m *MockSchoolRepositoryMockRecorder) GetSchool(ctx interface{}, id interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "GetSchool", reflect.TypeOf((*MockSchoolRepository)(nil).GetSchool), ctx, id)
}

func (m *MockSchoolRepositoryMockRecorder) UpdateSchool(ctx interface{}, school interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "UpdateSchool", reflect.TypeOf((*MockSchoolRepository)(nil).UpdateSchool), ctx, school)
}

func (m *MockSchoolRepositoryMockRecorder) DeleteSchool(ctx interface{}, id interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "DeleteSchool", reflect.TypeOf((*MockSchoolRepository)(nil).DeleteSchool), ctx, id)
}

func (m *MockSchoolRepository) CreateSchool(ctx context.Context, school *proto.School) (*proto.SchoolId, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSchool", ctx, school)
	ret0, _ := ret[0].(*proto.SchoolId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockSchoolRepository) ListSchools(ctx context.Context) ([]*proto.School, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSchools", ctx)
	ret0, _ := ret[0].([]*proto.School)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockSchoolRepository) GetSchool(ctx context.Context, id string) (*proto.School, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchool", ctx, id)
	ret0, _ := ret[0].(*proto.School)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockSchoolRepository) UpdateSchool(ctx context.Context, school *proto.School) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSchool", ctx, school)
	ret0, _ := ret[0].(error)
	return ret0
}

func (m *MockSchoolRepository) DeleteSchool(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSchool", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

func TestCreateSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockSchoolRepository(mockCtrl)
	uc := NewSchoolUsecases(mockRepo)

	ctx := context.Background()
	testSchool := &proto.School{
		Id:       primitive.NewObjectID().Hex(),
		SchoolId: "1",
		Name:     "Test School",
		Address:  "123 Test Street",
		Phone:    "123-456-7890",
	}

	mockRepo.EXPECT().CreateSchool(ctx, testSchool).Return(&proto.SchoolId{Id: "1"}, nil)

	schoolId, err := uc.CreateSchool(ctx, testSchool)

	assert.NoError(t, err)
	assert.NotNil(t, schoolId)
	assert.Equal(t, "1", schoolId.Id)
}

func TestListSchools(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockSchoolRepository(mockCtrl)
	uc := NewSchoolUsecases(mockRepo)

	ctx := context.Background()
	testSchools := []*proto.School{
		{
			Id:       primitive.NewObjectID().Hex(),
			SchoolId: "1",
			Name:     "Test School",
			Address:  "123 Test Street",
			Phone:    "123-456-7890",
		},
		{
			Id:       primitive.NewObjectID().Hex(),
			SchoolId: "2",
			Name:     "Test School 2",
			Address:  "456 Test Street",
			Phone:    "123-456-7890",
		},
	}

	mockRepo.EXPECT().ListSchools(ctx).Return(testSchools, nil)

	schools, err := uc.ListSchools(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, schools)
	assert.Equal(t, testSchools, schools)
}

func TestGetSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockSchoolRepository(mockCtrl)
	uc := NewSchoolUsecases(mockRepo)

	ctx := context.Background()
	testSchool := &proto.School{
		Id:       primitive.NewObjectID().Hex(),
		SchoolId: "1",
		Name:     "Test School",
		Address:  "123 Test Street",
		Phone:    "123-456-7890",
	}

	mockRepo.EXPECT().GetSchool(ctx, testSchool.Id).Return(testSchool, nil)

	school, err := uc.GetSchool(ctx, testSchool.Id)

	assert.NoError(t, err)
	assert.NotNil(t, school)
	assert.Equal(t, testSchool, school)
}

func TestUpdateSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockSchoolRepository(mockCtrl)
	uc := NewSchoolUsecases(mockRepo)

	ctx := context.Background()
	testSchool := &proto.School{
		Id:       primitive.NewObjectID().Hex(),
		SchoolId: "1",
		Name:     "Test School",
		Address:  "123 Test Street",
		Phone:    "123-456-7890",
	}

	mockRepo.EXPECT().UpdateSchool(ctx, testSchool).Return(nil)

	err := uc.UpdateSchool(ctx, testSchool)

	assert.NoError(t, err)
}

func TestDeleteSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockSchoolRepository(mockCtrl)
	uc := NewSchoolUsecases(mockRepo)

	ctx := context.Background()
	schoolId := primitive.NewObjectID().Hex()

	mockRepo.EXPECT().DeleteSchool(ctx, schoolId).Return(nil)

	err := uc.DeleteSchool(ctx, schoolId)

	assert.NoError(t, err)
}
