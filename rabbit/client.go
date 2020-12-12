package rabbit

import (
	"github.com/streadway/amqp"
	"github.com/ui-kreinhard/go-rabbit-tunnel/chamqp"
	"log"
	"net"
	"syscall"
)

type RabbitTunnelClient struct {
	channel    *chamqp.Channel	
	localTunIp string
}

func NewRabbitTunnelClient(channel *chamqp.Channel, localTunIp string) *RabbitTunnelClient {
	return &RabbitTunnelClient{
		channel,
		localTunIp,
	}
}

func (r *RabbitTunnelClient) Publish(dest string, netPackage []byte) error {
	exchangeName := "tunnel"
	routingKey := dest
	return r.channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        netPackage,
		})
}

func send(ip string, p []byte) {
	ipA := net.ParseIP(ip)
	var err error
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr:
			[4]byte{ipA[0], ipA[1], ipA[2], ipA[3]},
	}
	err = syscall.Sendto(fd, p, 0, &addr)
	if err != nil {
		log.Fatal("Sendto:", err)
	}
}

func (r *RabbitTunnelClient) Listen() {
	ipFrameChannel := make(chan amqp.Delivery)
	errChan := make(chan error)

	exchangeName := "tunnel"
	routingKey := r.localTunIp
	queue := r.localTunIp + "_tun"

	r.channel.ExchangeDeclare(exchangeName, "topic", false, false, false, false, nil, errChan)
	r.channel.QueueDeclare(queue, false, true, false, false, nil, nil, nil)
	r.channel.QueueBind(queue, routingKey, exchangeName, false, nil, nil)
	r.channel.Consume(queue, "", true, false, false, false, nil, ipFrameChannel, nil)
	for {
		select {
		case msg := <-ipFrameChannel:
			frame := msg.Body
			log.Println("Received frame via rabbit", frame)
			send(r.localTunIp, frame)	
		case err := <-errChan:
			log.Printf("Failed to listen for requests: %v", err)
		}

	}
}