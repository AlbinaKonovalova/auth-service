package controller

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AuthServiceController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	result, err := c.user.CreateUser(ctx, input.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Roles:    req.Roles,
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return createUserResultToProto(result), nil
}

func (c *AuthServiceController) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var isActive *bool
	if req.IsActive != nil {
		v := req.GetIsActive()
		isActive = &v
	}

	result, err := c.user.ListUsers(ctx, dto.UserListFilters{
		IsActive: isActive,
		Role:     req.Role,
		Page:     int(req.Page),
		PerPage:  int(req.PerPage),
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return listUsersResultToProto(result), nil
}

func (c *AuthServiceController) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.AdminUserInfo, error) {
	id, err := parseUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	result, err := c.user.GetUser(ctx, id)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return getUserResultToProto(result), nil
}

func (c *AuthServiceController) ActivateUser(ctx context.Context, req *pb.ActivateUserRequest) (*pb.ActivateUserResponse, error) {
	id, err := parseUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := c.user.ActivateUser(ctx, id); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.ActivateUserResponse{Success: true}, nil
}

func (c *AuthServiceController) DeactivateUser(ctx context.Context, req *pb.DeactivateUserRequest) (*pb.DeactivateUserResponse, error) {
	id, err := parseUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := c.user.DeactivateUser(ctx, id); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.DeactivateUserResponse{Success: true}, nil
}
