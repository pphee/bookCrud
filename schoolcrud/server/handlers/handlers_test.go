package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"

	pb "book/schoolcrud/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockISchoolUsecases struct {
	ctrl     *gomock.Controller
	recorder *MockISchoolUsecasesMockRecorder
}

type MockSchoolService_ListSchoolServer struct {
	pb.SchoolService_ListSchoolServer
	ctrl     *gomock.Controller
	recorder *MockSchoolService_ListSchoolServerMockRecorder
}

type MockSchoolService_ListSchoolServerMockRecorder struct {
	mock *MockSchoolService_ListSchoolServer
}

type MockISchoolUsecasesMockRecorder struct {
	mock *MockISchoolUsecases
}

func NewMockISchoolUsecases(ctrl *gomock.Controller) *MockISchoolUsecases {
	mock := &MockISchoolUsecases{ctrl: ctrl}
	mock.recorder = &MockISchoolUsecasesMockRecorder{mock}
	return mock
}

func NewMockSchoolService_ListSchoolServer(ctrl *gomock.Controller) *MockSchoolService_ListSchoolServer {
	mock := &MockSchoolService_ListSchoolServer{ctrl: ctrl}
	mock.recorder = &MockSchoolService_ListSchoolServerMockRecorder{mock}
	return mock
}

func (m *MockSchoolService_ListSchoolServer) EXPECTSTREAM() *MockSchoolService_ListSchoolServerMockRecorder {
	return m.recorder
}

func (m *MockISchoolUsecases) EXPECT() *MockISchoolUsecasesMockRecorder {
	return m.recorder
}

func (m *MockISchoolUsecasesMockRecorder) CreateSchool(ctx context.Context, school *pb.School) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "CreateSchool", reflect.TypeOf((*MockISchoolUsecases)(nil).CreateSchool), ctx, school)
}

func (m *MockISchoolUsecasesMockRecorder) ListSchools(ctx context.Context) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "ListSchools", reflect.TypeOf((*MockISchoolUsecases)(nil).ListSchools), ctx)
}

func (m *MockISchoolUsecasesMockRecorder) GetSchool(ctx context.Context, id string) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "GetSchool", reflect.TypeOf((*MockISchoolUsecases)(nil).GetSchool), ctx, id)
}

func (m *MockISchoolUsecasesMockRecorder) UpdateSchool(ctx context.Context, school *pb.School) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "UpdateSchool", reflect.TypeOf((*MockISchoolUsecases)(nil).UpdateSchool), ctx, school)
}

func (m *MockISchoolUsecasesMockRecorder) DeleteSchool(ctx context.Context, id string) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "DeleteSchool", reflect.TypeOf((*MockISchoolUsecases)(nil).DeleteSchool), ctx, id)
}

func (m *MockSchoolService_ListSchoolServer) Send(school *pb.School) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", school)
	ret0, _ := ret[0].(error)
	return ret0
}

func (m *MockISchoolUsecases) CreateSchool(ctx context.Context, school *pb.School) (*pb.SchoolId, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSchool", ctx, school)
	ret0, _ := ret[0].(*pb.SchoolId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockISchoolUsecases) ListSchools(ctx context.Context) ([]*pb.School, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSchools", ctx)
	ret0, _ := ret[0].([]*pb.School)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockISchoolUsecases) GetSchool(ctx context.Context, id string) (*pb.School, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchool", ctx, id)
	ret0, _ := ret[0].(*pb.School)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockISchoolUsecases) UpdateSchool(ctx context.Context, school *pb.School) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSchool", ctx, school)
	ret0, _ := ret[0].(error)
	return ret0
}

func (m *MockISchoolUsecases) DeleteSchool(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSchool", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

func TestCreateSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSchoolUsecases := NewMockISchoolUsecases(mockCtrl)
	handler := NewHandlers(mockSchoolUsecases)

	ctx := context.Background()
	testSchool := &pb.School{
		Id:       primitive.NewObjectID().Hex(),
		SchoolId: "test-school-id",
		Name:     "test-name",
		Address:  "test-address",
		Phone:    "test-phone",
	}
	expectedSchoolId := &pb.SchoolId{Id: testSchool.Id}

	mockSchoolUsecases.EXPECT().
		CreateSchool(ctx, testSchool).
		Return(expectedSchoolId, nil).
		Times(1)

	schoolId, err := handler.CreateSchool(ctx, testSchool)

	assert.NoError(t, err)
	assert.Equal(t, expectedSchoolId, schoolId)
}

func TestGetSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSchoolUsecases := NewMockISchoolUsecases(mockCtrl)
	handler := NewHandlers(mockSchoolUsecases)

	ctx := context.Background()
	testSchoolId := &pb.SchoolId{Id: primitive.NewObjectID().Hex()}
	testSchool := &pb.School{
		Id:       testSchoolId.Id,
		SchoolId: "1",
		Name:     "BKK school",
		Address:  "BKK",
		Phone:    "0872893601",
	}

	mockSchoolUsecases.EXPECT().
		GetSchool(ctx, testSchoolId.Id).
		Return(testSchool, nil).
		Times(1)

	school, err := handler.GetSchool(ctx, testSchoolId)

	assert.NoError(t, err)
	assert.Equal(t, testSchool, school)
}

func TestUpdateSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSchoolUsecases := NewMockISchoolUsecases(mockCtrl)
	handler := NewHandlers(mockSchoolUsecases)

	ctx := context.Background()
	testSchool := &pb.School{
		Id:       primitive.NewObjectID().Hex(),
		SchoolId: "2",
		Name:     "Yala school",
		Address:  "Yala",
		Phone:    "0892183101",
	}

	mockSchoolUsecases.EXPECT().
		UpdateSchool(ctx, testSchool).
		Return(nil).
		Times(1)

	_, err := handler.UpdateSchool(ctx, testSchool)

	assert.NoError(t, err)
}

func TestDeleteSchool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSchoolUsecases := NewMockISchoolUsecases(mockCtrl)
	handler := NewHandlers(mockSchoolUsecases)

	ctx := context.Background()
	testSchoolId := &pb.SchoolId{Id: primitive.NewObjectID().Hex()}

	mockSchoolUsecases.EXPECT().
		DeleteSchool(ctx, testSchoolId.Id).
		Return(nil).
		Times(1)

	_, err := handler.DeleteSchool(ctx, testSchoolId)

	assert.NoError(t, err)
}
