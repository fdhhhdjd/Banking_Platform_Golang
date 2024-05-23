package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/auth"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"google.golang.org/grpc/metadata"
)

func AuthorizeUser(c context.Context, accessibleRoles []string) (*auth.Payload, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(constants.AuthorizationHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != constants.AuthorizationTypeBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	JwtMaker, err := auth.GetJWTMaker()
	if err != nil {
		return nil, fmt.Errorf("JWT invalid: %s", err)
	}

	payload, err := JwtMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
