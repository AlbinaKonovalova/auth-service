package controller

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AuthServiceController) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	id, err := parseUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	roles, err := c.access.GetUserRoles(ctx, id)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return getUserRolesResultToProto(roles), nil
}

func (c *AuthServiceController) GetRolePermissions(ctx context.Context, req *pb.GetRolePermissionsRequest) (*pb.GetRolePermissionsResponse, error) {
	permissions, err := c.access.GetRolePermissions(ctx, req.RoleCode)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return getRolePermissionsResultToProto(permissions), nil
}

func (c *AuthServiceController) AssignPermission(ctx context.Context, req *pb.AssignPermissionRequest) (*pb.AssignPermissionResponse, error) {
	if req.Assignment == nil {
		return nil, status.Error(codes.InvalidArgument, "assignment is required")
	}

	if err := c.access.AssignPermission(ctx, input.AssignPermissionInput{
		RoleCode:       req.RoleCode,
		PermissionCode: req.Assignment.PermissionCode,
	}); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.AssignPermissionResponse{Success: true}, nil
}

func (c *AuthServiceController) RevokePermission(ctx context.Context, req *pb.RevokePermissionRequest) (*pb.RevokePermissionResponse, error) {
	if err := c.access.RevokePermission(ctx, input.RevokePermissionInput{
		RoleCode:       req.RoleCode,
		PermissionCode: req.PermissionCode,
	}); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.RevokePermissionResponse{Success: true}, nil
}

func (c *AuthServiceController) ListRoles(ctx context.Context, _ *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	roles, err := c.role.ListRoles(ctx)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return listRolesResultToProto(roles), nil
}

func (c *AuthServiceController) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	userID, err := parseUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if req.Assignment == nil {
		return nil, status.Error(codes.InvalidArgument, "assignment is required")
	}

	if err := c.access.AssignRole(ctx, input.AssignRoleInput{
		UserID:   userID,
		RoleCode: req.Assignment.RoleCode,
	}); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.AssignRoleResponse{Success: true}, nil
}

func (c *AuthServiceController) RevokeRole(ctx context.Context, req *pb.RevokeRoleRequest) (*pb.RevokeRoleResponse, error) {
	userID, err := parseUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := c.access.RevokeRole(ctx, input.RevokeRoleInput{
		UserID:   userID,
		RoleCode: req.RoleCode,
	}); err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return &pb.RevokeRoleResponse{Success: true}, nil
}
