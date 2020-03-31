package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"github.com/mirzaakhena/mirserver/messagebroker"
	"github.com/mirzaakhena/mirserver/utils/config"
	"github.com/mirzaakhena/mirserver/utils/log"
)

// MBProducer is
type MBProducer struct {
	MessageBroker messagebroker.IMessageBrokerProducer
	// DefaultMessage string
	ServerID string
}

func main() {

	cf := config.NewSimpleConfig("config", "./")

	log.GetLog().WithFile(cf.GetString("log.path", "./log"), "mirserver", "api", 2)

	nsqURL := cf.GetString("messagebroker.nsqd_url", "")

	log.GetLog().Info(nil, "NSQ address mirserver api: %s", nsqURL)

	serverID, _ := uuid.NewV4()

	x := MBProducer{
		MessageBroker: messagebroker.NewProducer(nsqURL),
		ServerID:      serverID.String(),
	}

	router := gin.Default()

	router.GET("/test", x.CallProducer)

	router.Run(":8080")
}

// CallProducer is
func (mb *MBProducer) CallProducer(c *gin.Context) {

	message := c.DefaultQuery("message", "Hello")

	rawData := map[string]interface{}{
		"message": message,
	}

	// publish to message broker
	byteData, _ := json.Marshal(rawData)
	if err := mb.MessageBroker.Publish("Test_Topic", byteData); err != nil {
		log.GetLog().Error(nil, "Failed to publish to message broker %s", err.Error())
	}

	log.GetLog().Info(nil, "raw data %v", string(byteData))

	c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Sending message %s from %s", message, mb.ServerID)})
}
