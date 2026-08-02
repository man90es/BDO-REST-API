package logger

import "time"

type Log struct {
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	InstanceID string    `json:"instanceId" bson:"instanceId"`
	Level      string    `json:"level" bson:"level"`
	Message    string    `json:"message" bson:"message"`
}
