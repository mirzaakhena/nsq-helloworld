package main

import (
	"encoding/json"

	"github.com/mirzaakhena/mirserver/messagebroker"
	"github.com/mirzaakhena/mirserver/utils/config"
	"github.com/mirzaakhena/mirserver/utils/log"
)

// MBConsumer is
type MBConsumer struct {
}

func main() {

	cf := config.NewSimpleConfig("config", "./")

	log.GetLog().WithFile(cf.GetString("log.path", "./log"), "mirserver", "con", 2)

	nsqURL := cf.GetString("messagebroker.nsqd_url", "")

	log.GetLog().Info(nil, "NSQ address mirserver con: %s", nsqURL)

	mb := MBConsumer{}
	x := messagebroker.NewConsumer([]messagebroker.ConsumerHandler{
		messagebroker.ConsumerHandler{
			Topic:               "Test_Topic",
			FunctionHandler:     mb.TestRequest,
			NumberOfConcurrency: 1,
		},
	})

	x.Run(nsqURL)

}

// TestRequest is
func (o *MBConsumer) TestRequest(m *messagebroker.Context) error {

	var req map[string]string
	if err := json.Unmarshal(m.Message, &req); err != nil {
		log.GetLog().Error(nil, "Fail convert json. %s", err.Error())
		return nil
	}

	byteData, _ := json.Marshal(req)
	log.GetLog().Info(nil, "Receive message: %v", string(byteData))

	return nil
}
