package gamecoordinator

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/Philipp15b/go-steam/v3"
	csgo "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	"github.com/Philipp15b/go-steam/v3/protocol/steamlang"
	"github.com/Philipp15b/go-steam/v3/totp"
)

// AppID describes the csgo app / steam id.
const AppID = 730

type GamecoordinatorService struct {
	client *steam.Client
}

func NewService() *GamecoordinatorService {
	return &GamecoordinatorService{}
}

// Connect connects to the steam service and starts a connection with the CSGO gamecoordinator.
func (s *GamecoordinatorService) Connect(username, password, twoFactorSecret string, connected sync.WaitGroup) {
	connected.Add(1)
	totpInstance := totp.NewTotp(twoFactorSecret)

	myLoginInfo := new(steam.LogOnDetails)
	myLoginInfo.Username = username
	myLoginInfo.Password = password
	twoFactorCode, err := totpInstance.GenerateCode()
	if err != nil {
		log.Error(err)
	}

	myLoginInfo.TwoFactorCode = twoFactorCode

	client := steam.NewClient()
	if _, connectErr := client.Connect(); connectErr != nil {
		log.Panic(connectErr)
	}

	for event := range client.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			log.Info("connected to steam. Logging in...")
			client.Auth.LogOn(myLoginInfo)
		case *steam.LoggedOnEvent:
			log.Info("logged on")
			s.client = client
			client.Social.SetPersonaState(steamlang.EPersonaState_Invisible)

			s.connectToGamecordinator()

			connected.Done()
		case steam.DisconnectedEvent:
			log.Panic("steam client: disconnected")
		case steam.FatalErrorEvent:
			log.Fatal(e)
		}
	}
}

func (s *GamecoordinatorService) connectToGamecordinator() {
	s.client.GC.RegisterPacketHandler(s)
	s.client.GC.SetGamesPlayed(730)
	s.shakeHands()
}

// shakeHands sends a hello to the GC.
func (s *GamecoordinatorService) shakeHands() {
	// Try to avoid not being ready on instant call of connection.
	time.Sleep(5 * time.Second)

	version := uint32(1)
	s.write(uint32(csgo.EGCBaseClientMsg_k_EMsgGCClientHello), &csgo.CMsgClientHello{
		Version: &version,
	})
}

// Write sends a message to the game coordinator.
func (s *GamecoordinatorService) write(messageType uint32, msg protoreflect.ProtoMessage) {
	s.client.GC.Write(gamecoordinator.NewGCMsgProtobuf(AppID, messageType, msg))
}
