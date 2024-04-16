package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Budhiarta/bank-film-BE/pkg/bootsrapper"
	"github.com/Budhiarta/bank-film-BE/pkg/config"
	smtp "github.com/Budhiarta/bank-film-BE/pkg/utils/smtp/impl"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/Budhiarta/bank-film-BE/pkg/database"
)

func init() {
	if os.Getenv("ENV") == "production" {
		return
	}

	//	load env variables from .env file for local development
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func main() {
	env := config.LoadConfig()

	ctx := context.Background()
	db, err := database.Connect(
		env["DB_HOST"],
		env["DB_PORT"],
		env["DB_USER"],
		env["DB_PASS"],
		env["DB_NAME"],
		5,
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	e := echo.New()
	mailer, err := smtp.InitSMTP(smtp.Config{
		Host:      "smtp.gmail.com",
		Port:      587,
		Username:  "budhitest1@gmail.com",
		Password:  "npzubnpjxsauqigg",
		From:      "budhitest1@gmail.com",
		QueueSize: 1,
		Workers:   1,
		Retry: smtp.RetryConfig{
			Delay: 5 * time.Second,
			Max:   1,
		},
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	mailer.StartMailWorker(ctx)

	bootsrapper.InitController(e, db, env, mailer)

	e.Logger.Fatal(e.Start(":" + env["PORT"]))

}
