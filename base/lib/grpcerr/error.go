package grpcerr

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(code codes.Code, msg string, details ...proto.Message) error {
	status := status.New(code, msg)
	status, _ = status.WithDetails(details...)
	return status.Err()
}
