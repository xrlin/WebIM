package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang-master"
	"log"
	"os"
)

var options *MQTT.ClientOptions

var MQTTClient MQTT.Client

func GetMQTTClient() (client MQTT.Client, err error) {
	if MQTTClient == nil {
		MQTTClient = MQTT.NewClient(options)
		token := MQTTClient.Connect()
		if token.Wait() && token.Error() != nil {
			err = token.Error()
			MQTTClient = nil
			return
		}
	}
	client = MQTTClient
	return
}

func createLogger(fileName string) *log.Logger {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return log.New(file, "My app Name ", log.Lshortfile)
}

var MQTTLogger = createLogger("mqtt_client.log")

func init() {
	options = MQTT.NewClientOptions()
	options.AddBroker("ws://localhost:1883")
	options.SetAutoReconnect(true)
	options.SetClientID("webim-server")
	MQTT.ERROR = MQTTLogger
	MQTT.WARN = MQTTLogger
	MQTT.CRITICAL = MQTTLogger
	MQTT.DEBUG = MQTTLogger
	GetMQTTClient()
}
