package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Cludch/csgo-microservices/demodownloader/internal/config"
	"github.com/Cludch/csgo-microservices/demodownloader/internal/domain/download"
	"github.com/Cludch/csgo-microservices/demodownloader/pkg/downloader"
	pb "github.com/Cludch/csgo-microservices/demodownloader/proto"
	"github.com/Cludch/csgo-microservices/shared/pkg/metrics"
	"github.com/Cludch/csgo-microservices/shared/pkg/queue"
	shared "github.com/Cludch/csgo-microservices/shared/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	prometheusPort    = 2111
	port              = 50051
	publishTopic      = "downloaded-demo"
	matchdetailsTopic = "gamedetails"
)

var configService *config.ConfigService
var downloaderService *downloader.DownloaderService
var queueService *queue.QueueService

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	configService = config.NewService()
	queueService = queue.NewService()

	downloaderService := downloader.New(&http.Client{})
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
	pb.RegisterUserServiceServer(s, download.NewDownloadHandler(downloaderService, configService.GetConfig().Downloader.DemosDir))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to grpc serve: %v", err)
	}
}

// MatchResponse contains information about the latest match.
type MatchDownloadDetails struct {
	ID        string    `json:"id"`
	DemoUrl   string    `json:"demoUrl"`
	StartTime time.Time `json:"startedAt"`
}

func consumeMessages() {
	log.Debugf("ceating queue consumer for %s", matchdetailsTopic)
	channel, err := queueService.Consume(matchdetailsTopic)
	if err != nil {
		log.Fatalf("unable to consume topic %s: %v", matchdetailsTopic, err)
	}
	go func() {
		for d := range channel {
			var matchDetails *shared.MatchDetails
			err := json.Unmarshal(d.Body, &matchDetails)
			if err != nil {
				log.Error("unable to unmarshal received data into share code data")
				if err := d.Reject(false); err != nil {
					log.Errorf("unable to reject message, %v", err)
				}
			}

			log.Debugf("downloading demo for %s", matchDetails.DownloadUrl)

			downloaderService.DownloadDemo(matchDetails.DownloadUrl, configService.GetConfig().Downloader.DemosDir, matchDetails.MatchTime)
		}
	}()
}

func publishDownloadedDemo(details *MatchDownloadDetails) error {
	json, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	ch, errPublish := queueService.Publish(json, publishTopic)
	ack := <-ch
	if !ack {
		return fmt.Errorf("unable to publish faceit match details for %v", details.ID)
	}

	log.Infof("published faceit match details for %v", details.ID)
	return errPublish
}
