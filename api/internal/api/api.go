package api

import (
	"context"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger
	db     *pgxpool.Pool
	socket *socketio.Server
	redis  *redis.Client
}

func MakeAPI(ctx context.Context, logger *zap.Logger, pg *pgxpool.Pool, ws *socketio.Server, redis *redis.Client) (*api, error) {
	return &api{
		logger: logger,
		db:     pg,
		socket: ws,
		redis:  redis,
	}, nil
}

func (a *api) makeRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/health", a.checkHealthHandler).Methods("GET")

	return r
}

func (a *api) makeWsRouter() {
	// todo
}

func (a *api) Server() *http.Server {
	a.makeWsRouter()

	root := mux.NewRouter()

	root.Handle("/ws", a.socket)
	root.Handle("/api", a.makeRouter())

	return &http.Server{
		Addr:    ":3001",
		Handler: root,
	}
}
