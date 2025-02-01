package gapi

import (
	"context"
	"fmt"
	"strings"

	"example.com/simplebank/token"
	"google.golang.org/grpc/metadata"
)

func (server *Server) authorizeUser(Ctx context.Context) (*token.Payload, error) {
	mtdt, ok := metadata.FromIncomingContext(Ctx)

	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	values := mtdt.Get("authorization")	
	if len(values) == 0 {
		return nil, fmt.Errorf("authorization token is not provided")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization token format")
	}

	authType := strings.ToLower(fields[0])

	if authType != "bearer" {	
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)

	if err != nil {
		return nil, fmt.Errorf("access token is not valid")
	}

	return payload, nil
}