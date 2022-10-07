package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Cludch/csgo-microservices/demoparser/internal/config"
	"github.com/Cludch/csgo-microservices/demoparser/internal/database"
	"github.com/Cludch/csgo-microservices/demoparser/internal/domain/demoparser"
	"github.com/Cludch/csgo-microservices/demoparser/internal/domain/match"
	"github.com/Cludch/csgo-microservices/demoparser/internal/domain/player"
	"github.com/Cludch/csgo-microservices/demoparser/pkg/files"
	"github.com/Cludch/csgo-microservices/shared/pkg/metrics"
	"github.com/Cludch/csgo-microservices/shared/pkg/queue"
	log "github.com/sirupsen/logrus"
)

var (
	prometheusPort         = 2111
	matchdetailsTopic      = "gamedetails"
	parserVersion     byte = 1
)

var configService *config.ConfigService
var databaseService *database.DatabaseService
var queueService *queue.QueueService
var matchService *match.MatchService
var playerService *player.PlayerService

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	configService = config.NewService()
	queueService = queue.NewService()
	databaseService = database.NewService()
	matchService = match.NewService(match.NewRepositoryMongo(databaseService))
	playerService = player.NewService(player.NewRepositoryMongo(databaseService))
}

func main() {
	// Connect to AMQP broker.
	if err := queueService.Connect(configService.GetConfig().Broker.Uri); err != nil {
		log.Fatal(err)
	}
	defer queueService.Connection.Close()

	// Connect to database.
	databaseConfig := configService.GetConfig().Database
	databaseService.Connect(databaseConfig.Host, databaseConfig.Username, databaseConfig.Password, databaseConfig.Database, databaseConfig.Port)

	// Create prometheus server.
	go metrics.PrometheusServer(prometheusPort)

	// Scan for new local files
	demos, _ := files.ScanDemosDir(configService.GetConfig().Parser.DemosDir)
	for _, demo := range demos {
		m, err := matchService.CreateMatchFromManualUpload(demo.Filename, demo.MatchTime)
		if err != nil {
			msg := "unable to create manual uploaded demo for file %s"
			log.Warn(msg, demo.Filename)
		} else if m != nil {
			msg := "found demo file %s and created manual upload entity"
			log.Infof(msg, m.Filename)
		}
	}

	log.Info("starting demoparser")

	numJobs := int64(configService.GetConfig().Parser.WorkerCount)
	matchQueue := make(chan *match.Match, numJobs)

	msg := "using %d workers"
	log.Infof(msg, numJobs)

	// Start numJobs-times parallel workers.
	for w := int64(1); w <= numJobs; w++ {
		go worker(matchQueue)
	}

	// Create a loop that checks for unparsed demos.
	t := time.NewTicker(time.Minute)
	for {
		// Get non-parsed matches from the db.
		nonParsedMatches, err := matchService.GetParseableMatches(parserVersion)

		if err != nil {
			log.Fatal(err)
		}

		// Enqueue found matches.
		for _, match := range nonParsedMatches {
			matchQueue <- match
		}

		<-t.C
	}
}

// Takes a match from the channel, parses and persists it.
func worker(matches <-chan *match.Match) {
	for m := range matches {
		filename := m.Filename
		if filename == "" {
			return
		}

		parser := demoparser.NewService(configService)
		demoFile := &files.Demo{ID: m.ID, MatchTime: m.CreatedAt, Filename: filename}

		// Check if file exists. File may have gotten deleted after being parsed the first time.
		if _, err := os.Stat(filepath.Join(configService.GetConfig().Parser.DemosDir, demoFile.Filename)); errors.Is(err, os.ErrNotExist) {
			// Set demo as unavailable.
			if err := matchService.SetStatusAndFilename(m, match.Unavailable, demoFile.Filename); err != nil {
				log.Warnf("Demo file %v for match with id %v is no longer available.", demoFile.Filename, demoFile.ID)
			}
		}

		if err := parser.Parse(configService.GetConfig().Parser.DemosDir, demoFile); err != nil {
			log.Error(err)
			continue
		}

		if !parser.GameOver {
			log.Errorf("Game %v did not finish before parsing ended. The file might be incomplete.", demoFile.Filename)
			continue
		}

		result := match.CreateResult(parser.Match)
		if err := matchService.UpdateResult(m, result, parserVersion); err != nil {
			log.Error(err)
			continue
		}

		for _, t := range m.Result.Teams {
			for _, playerResult := range t.Players {
				player, err := playerService.GetPlayer(playerResult.SteamID)
				if err != nil {
					const msg = "main: unable to query player: %s"
					log.Errorf(msg, err)
					continue
				}

				playerResult.MatchRounds = byte(len(m.Result.Rounds))
				playerResult.ScoreOwnTeam = t.Wins

				// This gets the team index in the array by turning the index around.
				// There could be a smarter way, but this is a fast one.
				enemyTeamId := (t.TeamID + 1) % 2
				playerResult.ScoreEnemyTeam = m.Result.Teams[enemyTeamId].Wins
				if err := playerService.AddResult(player, playerResult); err != nil {
					log.Error(err)
				}
			}
		}

		const msg = "demoparser: finished parsing %s"
		log.Infof(msg, filename)

		publishMatchResult(result)
	}
}

func publishMatchResult(result *match.MatchResult) error {
	json, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	ch, errPublish := queueService.Publish(json, matchdetailsTopic)
	ack := <-ch
	if !ack {
		return fmt.Errorf("unable to publish match result data")
	}

	log.Infof("published match result data")
	return errPublish
}
