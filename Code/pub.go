package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("Device-pub")

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)

	//we are going to try connecting for max 10 times to the server if the connection fails.
	for i := 0; i < 10; i++ {
		if token := c.Connect(); token.Wait() && token.Error() == nil {
			break
		} else {
			fmt.Println(token.Error())
			time.Sleep(1 * time.Second)
		}
	}

	if token := c.Subscribe("mqtt", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		var message string
		fmt.Print(">> ")
		// create a new bffer reader.
		reader := bufio.NewReader(os.Stdin)
		// read a string.
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if strings.Compare(message, "\n") > 0 {
			// if there is a message, publish it.
			token := c.Publish("mqtt", 0, false, message)
			if strings.Compare(message, "bye\n") == 0 {
				// if message == "bye" then exit the shell.
				break
			}
			token.Wait()
		}
	}

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("mqtt"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

}
