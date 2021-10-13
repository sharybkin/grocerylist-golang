package app

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sharybkin/grocerylist-golang/internal/http/handler"
	"github.com/sharybkin/grocerylist-golang/internal/repository"
	"github.com/sharybkin/grocerylist-golang/internal/service"
	"github.com/sharybkin/grocerylist-golang/pkg/http"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run()  {
	log.SetFormatter(new(log.JSONFormatter))

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	repository := repository.NewRepository()
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	srv := new(http.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("Error occured on shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")
	return viper.ReadInConfig()
}