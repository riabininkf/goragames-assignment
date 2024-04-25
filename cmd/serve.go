package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/riabininkf/goragames-assignment/internal/config"
	"github.com/riabininkf/goragames-assignment/internal/container"
	httpHandlers "github.com/riabininkf/goragames-assignment/internal/http"
	"github.com/riabininkf/goragames-assignment/internal/logger"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"time"

	_ "github.com/riabininkf/goragames-assignment/internal/handlers"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			var log logger.Logger
			if err := container.Fill(logger.DefName, &log); err != nil {
				return err
			}

			var cfg *config.Config
			if err := container.Fill(config.DefName, &cfg); err != nil {
				return err
			}

			gin.DefaultWriter = io.Discard

			router := gin.Default()
			router.Use(gin.Recovery())

			for _, defName := range container.GetByTag(httpHandlers.TagHandler) {
				var handler httpHandlers.Handler
				if err := container.Fill(defName, &handler); err != nil {
					return err
				}

				router.Handle(handler.Method(), handler.Path(), httpHandlers.WrapHandler(cmd.Context(), handler.Handle))
				log.Info(fmt.Sprintf("register handler %s %s", handler.Method(), handler.Path()))
			}

			srv := &http.Server{Addr: ":" + cfg.GetString("http.port"), Handler: router}

			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("error on serving", logger.Error(err))
				}
			}()

			log.Info("listening http requests on " + srv.Addr)
			<-cmd.Context().Done()
			log.Info("shutting down the http server")

			shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
			defer cancelFunc()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Error("can't gracefully shutdown http server", logger.Error(err))
				return nil
			}

			log.Info("http server is gracefully shutted down")
			return nil
		},
	})
}
