package controller

import (
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/AlbinaKonovalova/auth-service/internal/ports/input"
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

func createUserResultToProto(r input.CreateUserResult) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		Id:       r.ID.String(),
		Email:    r.Email,
		IsActive: r.IsActive,
		Roles:    r.Roles,
	}
}

func listUsersResultToProto(r input.ListUsersResult) *pb.ListUsersResponse {
	users := make([]*pb.AdminUserInfo, len(r.Users))
	for i, u := range r.Users {
		users[i] = &pb.AdminUserInfo{
			Id:       u.ID.String(),
			Email:    u.Email,
			IsActive: u.IsActive,
			Roles:    u.Roles,
		}
	}
	return &pb.ListUsersResponse{
		Users:   users,
		Total:   int32(r.Total),
		Page:    int32(r.Page),
		PerPage: int32(r.PerPage),
	}
}
