package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"bdo-rest-api/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)
var initialised = false
var mongoCollection *mongo.Collection

func InitLogger() {
	mongoURI := config.GetMongoDB()

	if len(mongoURI) > 0 {
		client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))

		if err != nil {
			panic(err)
		}

		mongoCollection = client.Database("bdo-rest-api").Collection("logs")
	}

	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	initialised = true

	Info(fmt.Sprintf("API initialised, configuration loaded:\n%v", config.SprintfConfig()))
}

func writeToMongo(level, message string) {
	if mongoCollection == nil {
		return
	}

	mongoCollection.InsertOne(context.TODO(), Log{
		CreatedAt: time.Now(),
		Level:     level,
		Message:   message,
	})
}

func Info(message string) {
	writeToMongo("INFO", message)

	if !config.GetVerbosity() || !initialised {
		return
	}

	infoLogger.Println(message)
}

func Error(message string) {
	writeToMongo("ERROR", message)

	if !config.GetVerbosity() || !initialised {
		return
	}

	errorLogger.Println(message)
}
