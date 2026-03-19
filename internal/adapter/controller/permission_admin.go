package controller

import (
	"context"

	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

func (c *AuthServiceController) ListPermissions(ctx context.Context, _ *pb.ListPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	permissions, err := c.permission.ListPermissions(ctx)
	if err != nil {
		return nil, c.domainErrToStatus(err)
	}

	return listPermissionsResultToProto(permissions), nil
}
