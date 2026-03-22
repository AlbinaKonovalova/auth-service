package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

func (c *AuthServiceController) ListPermissions(ctx context.Context, _ *pb.ListPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	permissions, err := c.permission.ListPermissions(ctx)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return listPermissionsResultToProto(permissions), nil
}

func (c *AuthServiceController) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.PermissionInfo, error) {
	if req.Permission == nil {
		return nil, status.Error(codes.InvalidArgument, "permission is required")
	}

	result, err := c.permission.CreatePermission(ctx, input.CreatePermissionInput{
		Code:        req.Permission.Code,
		Description: req.Permission.Description,
	})
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return permissionViewToProto(result), nil
}

func (c *AuthServiceController) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*pb.DeletePermissionResponse, error) {
	if err := c.permission.DeletePermission(ctx, req.PermissionCode); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.DeletePermissionResponse{Success: true}, nil
}
