package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
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
	mongoURI := viper.GetString("mongo")

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

	configPrintOut := fmt.Sprintf("\tPort:\t\t%v\n", viper.GetInt("port")) +
		fmt.Sprintf("\tProxies:\t%v\n", viper.GetStringSlice("proxy")) +
		fmt.Sprintf("\tVerbosity:\t%v\n", viper.GetBool("verbose")) +
		fmt.Sprintf("\tCache TTL:\t%v\n", viper.GetDuration("cacheTTL")) +
		fmt.Sprintf("\tRate limit:\t%v/min\n", viper.GetInt64("rateLimit")) +
		fmt.Sprintf("\tMongoDB:\t%v", viper.GetString("mongo"))

	Info(fmt.Sprintf("API initialised, configuration loaded:\n%v", configPrintOut))
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

	if !viper.GetBool("verbose") || !initialised {
		return
	}

	infoLogger.Println(message)
}

func Error(message string) {
	writeToMongo("ERROR", message)

	if !viper.GetBool("verbose") || !initialised {
		return
	}

	errorLogger.Println(message)
}
