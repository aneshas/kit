package main

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/tonto/kit/http"
	pb "github.com/tonto/kit/http/example/proto/account"
)

// NewAccountService creates new account service
func NewAccountService(hooks ...*twirp.ServerHooks) *Account {
	acc := Account{}

	acc.TwirpInit(
		pb.AccountPathPrefix,
		pb.NewAccountServer(&acc, twirp.ChainHooks(hooks...)),
	)

	return &acc
}

// Account twirp service
type Account struct {
	http.TwirpService
}

// Profile returns user profile
func (a *Account) Profile(ctx context.Context, req *pb.ProfileReq) (*pb.ProfileResp, error) {
	return &pb.ProfileResp{
		Name:    "John Doe",
		Email:   "john@doe.com",
		Address: "Sesame Street",
	}, nil
}
