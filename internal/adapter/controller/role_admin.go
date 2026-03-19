package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

func (c *AuthServiceController) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.RoleInfo, error) {
	if req.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	result, err := c.role.CreateRole(ctx, input.CreateRoleInput{
		Code:        req.Role.Code,
		Name:        req.Role.Name,
		Description: req.Role.Description,
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return roleViewToProto(result), nil
}
