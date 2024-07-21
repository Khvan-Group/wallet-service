package main

import (
	"github.com/Khvan-Group/common-library/logger"
	"github.com/Khvan-Group/common-library/utils"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"wallet-service/internal/api"
	"wallet-service/internal/core/rabbitmq"
	"wallet-service/internal/database"
)

const SERVER_PORT = "SERVER_PORT"

func main() {
	start()
}

func start() {
	// init logger
	logger.InitLogger()
	logger.Logger.Info("Starting server")

	// load environments
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// init RabbitMQ
	rabbitmq.InitRabbitMQ()

	// init database
	database.InitDB()

	// init server
	port := ":" + utils.GetEnv(SERVER_PORT)
	r := mux.NewRouter()
	a := api.New()
	a.AddRoutes(r)

	rabbitmq.ConsumeRabbitMQ(a.Wallets.Service)

	logger.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
