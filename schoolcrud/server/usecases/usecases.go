package usecases

import (
	pb "book/schoolcrud/proto"
	"book/schoolcrud/server/repo"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type ISchoolUsecases interface {
	CreateSchool(ctx context.Context, school *pb.School) (*pb.SchoolId, error)
	ListSchools(ctx context.Context) ([]*pb.School, error)
	GetSchool(ctx context.Context, id string) (*pb.School, error)
	UpdateSchool(ctx context.Context, school *pb.School) error
	DeleteSchool(ctx context.Context, id string) error
}

type schoolUsecases struct {
	repo repo.ISchoolRepository
}

func NewSchoolUsecases(repo repo.ISchoolRepository) ISchoolUsecases {
	return &schoolUsecases{repo: repo}
}

func (s *schoolUsecases) CreateSchool(ctx context.Context, school *pb.School) (*pb.SchoolId, error) {
	schoolId, err := s.repo.CreateSchool(ctx, school)
	if err != nil {
		log.Printf("Error creating school: %v", err)
		return nil, status.Errorf(codes.Internal, "Error creating school")
	}
	return schoolId, nil
}

func (s *schoolUsecases) ListSchools(ctx context.Context) ([]*pb.School, error) {
	schools, err := s.repo.ListSchools(ctx)
	if err != nil {
		log.Printf("Error listing schools: %v", err)
		return nil, status.Errorf(codes.Internal, "Error listing schools")
	}
	return schools, nil
}

func (s *schoolUsecases) GetSchool(ctx context.Context, id string) (*pb.School, error) {
	school, err := s.repo.GetSchool(ctx, id)
	if err != nil {
		log.Printf("Error getting school with ID %s: %v", id, err)
		return nil, status.Errorf(codes.Internal, "Error getting school")
	}
	return school, nil
}

func (s *schoolUsecases) UpdateSchool(ctx context.Context, school *pb.School) error {
	err := s.repo.UpdateSchool(ctx, school)
	if err != nil {
		log.Printf("Error updating school with ID %s: %v", school.Id, err)
		return status.Errorf(codes.Internal, "Error updating school")
	}
	return nil
}

func (s *schoolUsecases) DeleteSchool(ctx context.Context, id string) error {
	err := s.repo.DeleteSchool(ctx, id)
	if err != nil {
		log.Printf("Error deleting school with ID %s: %v", id, err)
		return status.Errorf(codes.Internal, "Error deleting school")
	}
	return nil
}
