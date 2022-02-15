package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Cludch/csgo-microservices/shared/pkg/metrics"
	"github.com/Cludch/csgo-microservices/shared/pkg/queue"
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"github.com/Cludch/csgo-microservices/valveapiclient/internal/config"
	"github.com/Cludch/csgo-microservices/valveapiclient/internal/database"
	"github.com/Cludch/csgo-microservices/valveapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/valveapiclient/internal/domain/valve_api"
	"github.com/Cludch/csgo-microservices/valveapiclient/pkg/api_client"
	"github.com/Cludch/csgo-microservices/valveapiclient/pkg/valve_match_api"
	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	prometheusPort = 2112
	port           = 50051
	topic          = "sharecode"
)

var configService *config.ConfigService
var userService *user.UserService
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

	userService = user.NewService(user.NewRepositoryPostgres(db), valve_match_api.New(&api_client.HttpApiClient{}))
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
	pb.RegisterValveMatchApiServiceServer(s, new(valve_api.ValveMatchApiHandler))
	pb.RegisterUserServiceServer(s, user.NewUserHandler(userService))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to grpc serve: %v", err)
	}
}

func query() {
	// Create a loop that checks for new share codes each minute.
	t := time.NewTicker(time.Minute)
	for {
		users, err := userService.GetUsersWithApiEnabled()

		if err != nil {
			log.Fatal(err)
		}

		// Iterate all csgo users and request the next share code for the latest share code.
		for _, u := range users {
			sc, err := userService.QueryLatestShareCode(u)
			if err != nil {
				log.Error(err)
			}

			if sc == nil {
				continue
			}

			if err = publishShareCode(sc); err != nil {
				const msg = "unable to publish share code %s"
				log.Errorf(msg, err)
				continue
			}

			if err = userService.UpdateLatestShareCode(u, sc); err != nil {
				const msg = "unable to update user latest share code %s"
				log.Errorf(msg, err)
				continue
			}
		}

		<-t.C
	}
}

func publishShareCode(sc *share_code.ShareCodeData) error {
	json, err := json.Marshal(sc)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	ch, errPublish := queueService.Publish(json, topic)
	ack := <-ch
	if !ack {
		return fmt.Errorf("unable to publish share code data for %v", sc.Encoded)
	}

	log.Infof("published share code data for %v", sc.Encoded)
	return errPublish
}
