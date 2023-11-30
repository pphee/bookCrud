package handlers

import (
	pb "book/schoolcrud/proto"
	"book/schoolcrud/server/usecases"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type SchoolHandlers struct {
	schoolHandler usecases.ISchoolUsecases
	pb.UnimplementedSchoolServiceServer
}

func NewHandlers(schoolHandler usecases.ISchoolUsecases) pb.SchoolServiceServer {
	return &SchoolHandlers{
		schoolHandler: schoolHandler,
	}
}

func (h *SchoolHandlers) CreateSchool(ctx context.Context, req *pb.School) (*pb.SchoolId, error) {
	log.Println("CreateSchool")

	schoolId, err := h.schoolHandler.CreateSchool(ctx, req)
	if err != nil {
		// Log the error and return a gRPC status error
		log.Printf("Error creating school: %v", err)
		return nil, status.Errorf(codes.Internal, "Error creating school: %v", err)
	}

	return schoolId, nil
}

func (h *SchoolHandlers) ListSchool(_ *emptypb.Empty, stream pb.SchoolService_ListSchoolServer) error {
	log.Println("ListSchools")

	schools, err := h.schoolHandler.ListSchools(context.Background())
	if err != nil {
		log.Printf("Error listing schools: %v", err)
		return status.Errorf(
			codes.Internal,
			"Unknown internal error: %v", err,
		)
	}

	for _, school := range schools {
		if err := stream.Send(school); err != nil {
			log.Printf("Error while streaming school data: %v", err)
			return status.Errorf(
				codes.Internal,
				"Error while streaming school data: %v", err,
			)
		}
	}

	return nil
}

func (h *SchoolHandlers) GetSchool(ctx context.Context, req *pb.SchoolId) (*pb.School, error) {
	log.Println("GetSchool")

	school, err := h.schoolHandler.GetSchool(ctx, req.Id)
	if err != nil {
		log.Printf("Error getting school: %v", err)
		return nil, status.Errorf(codes.Internal, "Error getting school: %v", err)
	}

	return school, nil
}

func (h *SchoolHandlers) UpdateSchool(ctx context.Context, req *pb.School) (*emptypb.Empty, error) {
	log.Println("UpdateSchool")

	err := h.schoolHandler.UpdateSchool(ctx, req)
	if err != nil {
		log.Printf("Error updating school: %v", err)
		return nil, status.Errorf(codes.Internal, "Error updating school: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *SchoolHandlers) DeleteSchool(ctx context.Context, req *pb.SchoolId) (*emptypb.Empty, error) {
	log.Println("DeleteSchool")

	err := h.schoolHandler.DeleteSchool(ctx, req.Id)
	if err != nil {
		log.Printf("Error deleting school: %v", err)
		return nil, status.Errorf(codes.Internal, "Error deleting school: %v", err)
	}

	return &emptypb.Empty{}, nil
}
