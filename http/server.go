package http

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Server struct {
	ctx context.Context
	rdb redis.Client
}

func NewServer() *Server {
	return &Server{
		ctx: context.Background(),
		rdb: *redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (s *Server) Run() {

}
