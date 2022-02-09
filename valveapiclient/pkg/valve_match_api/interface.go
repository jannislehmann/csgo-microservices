package valve_match_api

type ValveMatchApiConsumerUseCase interface {
	RequestNextShareCode(steamApiKey string, steamID uint64, historyAuthenticationCode string, lastShareCode string) (string, error)
}
