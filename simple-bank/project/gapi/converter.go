package gapi

import (
	db "github.com/ahmad-khatib0/go/simple-bank/project/db/sqlc"
	"github.com/ahmad-khatib0/go/simple-bank/project/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
