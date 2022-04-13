package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Cludch/csgo-microservices/faceitapiclient/internal/config"
	"github.com/Cludch/csgo-microservices/faceitapiclient/internal/database"
	"github.com/Cludch/csgo-microservices/faceitapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/faceitapiclient/pkg/faceit_api"
	pb "github.com/Cludch/csgo-microservices/faceitapiclient/proto"
	"github.com/Cludch/csgo-microservices/shared/pkg/metrics"
	"github.com/Cludch/csgo-microservices/shared/pkg/queue"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	prometheusPort = 2111
	port           = 50051
	topic          = "matchdownload"
)

var configService *config.ConfigService
var userService *user.UserService
var faceitApiConsumer *faceit_api.FaceitApiConsumerService
var queueService *queue.QueueService

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	configService = config.NewService()
	queueService = queue.NewService()

	dbConfig := configService.GetConfig().Database
	db, err := database.Connect(dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	faceitApiConsumer := faceit_api.New(&http.Client{})
	userService = user.NewService(user.NewRepositoryPostgres(db), faceitApiConsumer)
}

func main() {
	// Connect to AMQP broker.
	if err := queueService.Connect(configService.GetConfig().Broker.Uri); err != nil {
		log.Fatal(err)
	}
	defer queueService.Connection.Close()
	if err := queueService.CreateQueue(topic); err != nil {
		log.Fatal(err)
	}

	// Query for new share codes.
	go query()

	// Create prometheus server.
	go metrics.PrometheusServer(prometheusPort)

	// Create loopback gRPC server.
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register RPC handler.
	pb.RegisterUserServiceServer(s, user.NewUserHandler(userService))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to grpc serve: %v", err)
	}
}

func query() {
	// Create a loop that checks for new share codes each minute.
	t := time.NewTicker(time.Minute)
	faceitAPIKey := configService.GetConfig().Faceit.APIKey
	for {
		users, err := userService.GetUsersWithApiEnabled()

		if err != nil {
			log.Fatal(err)
		}

		// Iterate all registered faceit users and request the download urls of the last finished matches.
		for _, u := range users {
			playerMatchHistory, err := faceitApiConsumer.GetPlayerMatchHistory(faceitAPIKey, u.ID)
			if err != nil {
				log.Error(err)
			}

			if playerMatchHistory.Result == nil {
				continue
			}

			for _, matchHistory := range playerMatchHistory.Result {
				matchId := matchHistory.MatchId
				matchDetails, err := faceitApiConsumer.GetMatchDetails(faceitAPIKey, matchId)

				if err != nil {
					log.Error(err)
				}

				if matchDetails == nil || matchDetails.Status != "FINISHED" {
					continue
				}

				downloadUrl := matchDetails.DemoUrl[0]
				startTime := time.Unix(matchDetails.StartTime, 0)
				details := &MatchDownloadDetails{
					ID:        matchId,
					DemoUrl:   downloadUrl,
					StartTime: startTime,
				}
				if err = publishFaceitMatchDetails(details); err != nil {
					const msg = "unable to publish faceit match details %s"
					log.Errorf(msg, err)
					continue
				}

				log.Infof("created downloadable faceit match for id %s", matchId)
			}
		}

		<-t.C
	}
}

// MatchResponse contains information about the latest match.
type MatchDownloadDetails struct {
	ID        string    `json:"id"`
	DemoUrl   string    `json:"demoUrl"`
	StartTime time.Time `json:"startedAt"`
}

func publishFaceitMatchDetails(details *MatchDownloadDetails) error {
	json, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	ch, errPublish := queueService.Publish(json, topic)
	ack := <-ch
	if !ack {
		return fmt.Errorf("unable to publish faceit match details for %v", details.ID)
	}

	log.Infof("published faceit match details for %v", details.ID)
	return errPublish
}
