package cmd

import (
	"context"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peonii/feta/internal/api"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func APICmd(ctx context.Context) *cobra.Command {
	var apiCmd = &cobra.Command{
		Use: "api",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := zap.NewProduction()
			if os.Getenv("ENV") != "production" {
				logger, err = zap.NewDevelopment()
			}

			if err != nil {
				return err
			}

			pg, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
			if err != nil {
				logger.Error("Failed to connect to database",
					zap.Error(err),
				)
				return err
			}

			ws := socketio.NewServer(&engineio.Options{
				Transports: []transport.Transport{
					&polling.Transport{
						CheckOrigin: allowOriginFunc,
					},
					&websocket.Transport{
						CheckOrigin: allowOriginFunc,
					},
				},
			})

			rdb := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			})

			a, err := api.MakeAPI(ctx, logger, pg, ws, rdb)
			if err != nil {
				logger.Error("Failed to make API",
					zap.Error(err),
				)

				return err
			}

			srv := a.Server()

			// go func() {
			// 	if err := ws.Serve(); err != nil {
			// 		logger.Error("Socket.io listen error",
			// 			zap.Error(err),
			// 		)
			// 	}
			// }()

			// logger.Info("Started Socket.IO server")

			go func() { _ = srv.ListenAndServe() }()

			logger.Info("Started HTTP server")

			<-ctx.Done()

			srv.Shutdown(ctx)
			ws.Close()

			return nil
		},
	}

	return apiCmd
}
