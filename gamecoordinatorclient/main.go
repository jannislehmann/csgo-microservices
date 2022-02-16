package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Cludch/csgo-microservices/gamecoordinatorclient/internal/config"
	"github.com/Cludch/csgo-microservices/gamecoordinatorclient/internal/gamecoordinator"
	"github.com/Cludch/csgo-microservices/shared/pkg/metrics"
	"github.com/Cludch/csgo-microservices/shared/pkg/queue"
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"

	pb "github.com/Cludch/csgo-microservices/gamecoordinatorclient/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	prometheusPort = 2112
	port           = 50052
	sharecodeTopic = "sharecode"
	publishTopic   = "gamedetails"
)

var configService *config.ConfigService
var queueService *queue.QueueService
var gamecoordinatorService *gamecoordinator.GamecoordinatorService

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	configService = config.NewService()
	queueService = queue.NewService()
	gamecoordinatorService = gamecoordinator.NewService()
}

func main() {
	// Connect to AMQP broker.
	if err := queueService.Connect(configService.GetConfig().Broker.Uri); err != nil {
		log.Fatal(err)
	}
	defer queueService.Connection.Close()
	if err := queueService.CreateQueue(publishTopic); err != nil {
		log.Fatal(err)
	}

	steamConfig := configService.GetConfig().Steam
	gamecoordinatorService.Connect(steamConfig.Username, steamConfig.Password, steamConfig.TwoFactorSecret)

	go consumeMessages()

	// Create prometheus server.
	go metrics.PrometheusServer(prometheusPort)

	// Create loopback gRPC server.
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register RPC handler.
	pb.RegisterMatchDetailQueryServiceServer(s, gamecoordinator.NewGamecoordinatorApiHandler(gamecoordinatorService))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to grpc serve: %v", err)
	}
}

func consumeMessages() {
	log.Debugf("ceating queue consumer for %s", sharecodeTopic)
	channel, err := queueService.Consume(sharecodeTopic)
	if err != nil {
		log.Fatalf("unable to consume topic %s: %v", sharecodeTopic, err)
	}
	go func() {
		for d := range channel {
			var sc *share_code.ShareCodeData
			err := json.Unmarshal(d.Body, &sc)
			if err != nil {
				log.Error("unable to unmarshal received data into share code data")
				if err := d.Reject(false); err != nil {
					log.Errorf("unable to reject message, %v", err)
				}
			}

			log.Debugf("requesting match details for %s", sc.Encoded)

			select {
			case details := <-gamecoordinatorService.RequestMatchDetails(sc):
				const msg = "received response for %s"
				log.Debugf(msg, sc.Encoded)
				log.Printf("received match details for: %s", sc.Encoded)
				if err := publishMatchDetails(details); err != nil {
					log.Errorf("unable to publish match details for %d %v", details.MatchId, err)
					if err := d.Nack(false, true); err != nil {
						log.Errorf("unable to nack message, %v", err)
					}
				} else {
					log.Debugf("published match details. acknowledging message")
					if err := d.Ack(false); err != nil {
						log.Errorf("unable to ack message, %v", err)
					}
				}
			case <-time.After(15 * time.Second):
				const msg = "failed to receive response for %s"
				log.Debugf(msg, sc.Encoded)
				// Requeue the message as this could be an issue with the gc.
				if err := d.Nack(false, true); err != nil {
					log.Errorf("unable to nack message, %v", err)
				}
			}
		}
	}()
}

func publishMatchDetails(matchDetails *gamecoordinator.MatchDetails) error {
	json, err := json.Marshal(matchDetails)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	ch, errPublish := queueService.Publish(json, publishTopic)
	ack := <-ch
	if !ack {
		return fmt.Errorf("unable to publish match details for %d", matchDetails.MatchId)
	}

	log.Infof("published match details for %d", matchDetails.MatchId)
	return errPublish
}
