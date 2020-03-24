package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Que -.
type Que struct {
	Server    string
	Sendtopic string
	Resvtopic string
	Qos       int
	Retained  bool
	Clientid  string
	Username  string
	Password  string
	Client    MQTT.Client
	Callback  MQTT.MessageHandler
}

// Connect - .
func (q *Que) Connect() error {
	connOpts := MQTT.NewClientOptions().AddBroker(q.Server).SetClientID(q.Clientid).SetCleanSession(true)
	if q.Username != "" {
		connOpts.SetUsername(q.Username)
		if q.Password != "" {
			connOpts.SetPassword(q.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	q.Client = client

	if q.Callback != nil {
		if token := q.Client.Subscribe(q.Resvtopic, byte(q.Qos), q.Callback); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		fmt.Printf("[MQTT] Subscribe to %s\n", q.Sendtopic)
	}

	fmt.Printf("[MQTT] Connected to %s\n", q.Server)

	return nil
}

// Publish - .
func (q *Que) Publish(message string) error {
	if q.Client != nil {
		token := q.Client.Publish(q.Sendtopic, byte(q.Qos), q.Retained, message)
		if token == nil {
			return token.Error()
		}
		fmt.Printf("[MQTT] Sent to %s\n", q.Sendtopic)
	}

	return nil
}

// SetSubscribe - .
func (q *Que) SetSubscribe(callback MQTT.MessageHandler) error {
	if callback != nil {
		if token := q.Client.Subscribe(q.Resvtopic, byte(q.Qos), callback); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		fmt.Printf("[MQTT] Subscribe to %s\n", q.Sendtopic)
	}

	return nil
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("[MQTT] Received to %s [Received Message: %s]\n", message.Topic(), message.Payload())
}

func main() {
	hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://localhost:1883", "The full URL of the MQTT server to connect to")
	sendtopic := flag.String("sendtopic", "MQTT/Client/Update/TEST", "Topic to publish the messages on")
	resvtopic := flag.String("resvtopic", "MQTT/+/Update/#", "Topic to publish the messages on")
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	q := &Que{
		Server:    *server,
		Sendtopic: *sendtopic,
		Resvtopic: *resvtopic,
		Qos:       *qos,
		Retained:  *retained,
		Clientid:  *clientid,
		Username:  *username,
		Password:  *password,
		Callback:  onMessageReceived,
	}

	err := q.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if err := q.SetSubscribe(onMessageReceived); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for {
		time.Sleep(5000 * time.Millisecond)
		if err := q.Publish("test massage"); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}
