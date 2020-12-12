package main

import (
	"github.com/ui-kreinhard/go-rabbit-tunnel/util"
	"github.com/ui-kreinhard/go-rabbit-tunnel/chamqp"
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
	"github.com/ui-kreinhard/go-rabbit-tunnel/rabbit"
	"log"
	"os"
)

func setupTunDevice(tunIp string) error {
	output, err := util.Exec("/sbin/ip", "addr", "add", tunIp + "/24", "dev", "O_O")
	log.Println(output)
	if err == nil {
		_, err = util.Exec("ip", "link", "set", "dev", "O_O", "up")
		return err
	}
	return err
}

func main() {
	rabbitUrl := os.Args[1]
	localTunIp := os.Args[2]
	
	
	conn := chamqp.Dial(rabbitUrl)
	channel := conn.Channel()
	rabbitClient := rabbit.NewRabbitTunnelClient(channel, localTunIp)
	go rabbitClient.Listen()
	
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = "O_O"

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	
	err = setupTunDevice(localTunIp)
	if err != nil {
		log.Println(err)
	}

	var frame ethernet.Frame

	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Fatal(err)
		}
		frame = frame[:n]
		
		ip := waterutil.IPv4Destination(frame)
		
		log.Println("IP dest: ", ip.String())
		log.Printf("Dst: %s\n", frame.Destination())
		log.Printf("Src: %s\n", frame.Source())
		log.Printf("Ethertype: % x\n", frame.Ethertype())
		log.Printf("Payload: % x\n", frame.Payload())

		err = rabbitClient.Publish(ip.String(), frame)
		if err != nil {
			log.Println("Publish err", err)
		}
	}
}