package gapi

import (
	"fmt"

	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/pb"
	"github.com/haniifac/simplebank/token"
	"github.com/haniifac/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	pb.UnimplementedSimpleBankServer
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
