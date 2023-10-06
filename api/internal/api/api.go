package api

import (
	"context"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peonii/feta/internal/models"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger
	db     *pgxpool.Pool
	socket *socketio.Server
	redis  *redis.Client

	userRepo *models.UserRepo
}

func MakeAPI(ctx context.Context, logger *zap.Logger, pg *pgxpool.Pool, ws *socketio.Server, redis *redis.Client) (*api, error) {
	ur := models.NewUserRepo(ctx, pg)

	return &api{
		logger: logger,
		db:     pg,
		socket: ws,
		redis:  redis,

		userRepo: ur,
	}, nil
}

func (a *api) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Got request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)
	})
}

func (a *api) makeRouter(r *mux.Router) *mux.Router {

	// Health
	r.HandleFunc("/v1/health", a.checkHealthHandler).Methods("GET")

	// User
	r.HandleFunc("/v1/users", a.createUser).Methods("POST")

	r.Use(a.loggingMiddleware)

	return r
}

func (a *api) makeWsRouter() {
	// todo
}

func (a *api) Server() *http.Server {
	a.makeWsRouter()

	root := mux.NewRouter()

	root.Handle("/ws", a.socket)
	apiR := root.PathPrefix("/api").Subrouter()
	root.Handle("/", a.makeRouter(apiR))

	return &http.Server{
		Addr:    ":3001",
		Handler: root,
	}
}
