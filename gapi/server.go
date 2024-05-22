package gapi

import (
	database "github.com/fdhhhdjd/Banking_Platform_Golang/database/sqlc"
	"github.com/fdhhhdjd/Banking_Platform_Golang/pb"
)

type SimpleBankServer struct {
	pb.UnimplementedSimpleBankServer
	store database.Store
}

func NewSimpleBankServer(store database.Store) *SimpleBankServer {
	return &SimpleBankServer{store: store}
}
