package controller

import (
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

func loginResultToProto(result dto.LoginResult) *pb.LoginResponse {
	return &pb.LoginResponse{
		AccessToken: result.AccessToken,
		User:        currentUserToProto(result.User),
	}
}

func currentUserToProto(u dto.CurrentUser) *pb.UserInfo {
	return &pb.UserInfo{
		Id:          u.ID.String(),
		Roles:       u.Roles,
		Permissions: u.Permissions,
	}
}

func refreshResultToProto(result dto.RefreshResult) *pb.RefreshResponse {
	return &pb.RefreshResponse{
		AccessToken: result.AccessToken,
		User:        currentUserToProto(result.User),
	}
}

func adminUserViewToProto(v domainservice.AdminUserView) *pb.AdminUserInfo {
	return &pb.AdminUserInfo{
		Id:       v.ID.String(),
		Email:    v.Email,
		IsActive: v.IsActive,
		Roles:    v.Roles,
	}
}

func createUserResultToProto(v domainservice.AdminUserView) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		Id:       v.ID.String(),
		Email:    v.Email,
		IsActive: v.IsActive,
		Roles:    v.Roles,
	}
}

func userListToProto(v domainservice.UserList) *pb.ListUsersResponse {
	users := make([]*pb.AdminUserInfo, len(v.Items))
	for i, item := range v.Items {
		users[i] = adminUserViewToProto(item)
	}

	return &pb.ListUsersResponse{
		Users:   users,
		Total:   int32(v.Total),
		Page:    int32(v.Page),
		PerPage: int32(v.PerPage),
	}
}

func getUserResultToProto(v domainservice.AdminUserView) *pb.AdminUserInfo {
	return adminUserViewToProto(v)
}

func getUserRolesResultToProto(roles []domainservice.RoleView) *pb.GetUserRolesResponse {
	pbRoles := make([]*pb.RoleInfo, len(roles))
	for i, r := range roles {
		pbRoles[i] = roleViewToProto(r)
	}
	return &pb.GetUserRolesResponse{Roles: pbRoles}
}

func roleViewToProto(r domainservice.RoleView) *pb.RoleInfo {
	return &pb.RoleInfo{
		Id:          r.ID.String(),
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
	}
}

func listRolesResultToProto(roles []domainservice.RoleView) *pb.ListRolesResponse {
	pbRoles := make([]*pb.RoleInfo, len(roles))
	for i, r := range roles {
		pbRoles[i] = roleViewToProto(r)
	}
	return &pb.ListRolesResponse{Roles: pbRoles}
}

func listPermissionsResultToProto(permissions []domainservice.PermissionView) *pb.ListPermissionsResponse {
	pbPerms := make([]*pb.PermissionInfo, len(permissions))
	for i, p := range permissions {
		pbPerms[i] = permissionViewToProto(p)
	}
	return &pb.ListPermissionsResponse{Permissions: pbPerms}
}

func permissionViewToProto(p domainservice.PermissionView) *pb.PermissionInfo {
	return &pb.PermissionInfo{
		Id:          p.ID.String(),
		Code:        p.Code,
		Description: p.Description,
	}
}

func getRolePermissionsResultToProto(permissions []domainservice.PermissionView) *pb.GetRolePermissionsResponse {
	pbPerms := make([]*pb.PermissionInfo, len(permissions))
	for i, p := range permissions {
		pbPerms[i] = permissionViewToProto(p)
	}
	return &pb.GetRolePermissionsResponse{Permissions: pbPerms}
}
