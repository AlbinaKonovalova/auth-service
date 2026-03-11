package controller

import (
	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
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
