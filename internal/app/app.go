package app

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/katenester/doc/internal/repository"
	"github.com/katenester/doc/internal/repository/postgres/config"
	"github.com/katenester/doc/internal/service"
	"github.com/katenester/doc/internal/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// Run - Building dependencies and logic
func Run() {
	// Download variables env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error initalization db password(file env) %s", err.Error())
	}
	db, err := config.NewPostgresDB(config.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("error initalization db %s", err.Error())
	}
	// Dependency injection for architecture application
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := transport.NewHandler(services)
	srv := new(transport.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server %s", err.Error())
		}
	}()

	logrus.Print("todo server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured while shutting down server %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Fatalf("error occured while closing db %s", err.Error())
	}
}
