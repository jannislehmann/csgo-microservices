package faceit_api

// MatchResponse contains information about the latest match.
type MatchDetailsResponse struct {
	DemoUrl   []string `json:"demo_url"`
	StartTime int64    `json:"started_at"`
	Status    string   `json:"status"`
}

// MatchResponse contains information about the latest match.
type PlayerMatchHistoryResponse struct {
	Result []*PlayerMatchHistoryEntry `json:"items"`
}

type PlayerMatchHistoryEntry struct {
	MatchId string `json:"match_id"`
}

// InvalidFaceitApiCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type InvalidFaceitApiCredentials struct{}

func (e *InvalidFaceitApiCredentials) Error() string {
	return "Invalid faceit api credentials"
}

// InvalidFaceitApiCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type FaceitApiConnectionIssues struct{}

func (e *FaceitApiConnectionIssues) Error() string {
	return "Too many requests or unavailable api"
}
