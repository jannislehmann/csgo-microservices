package match

import (
	"errors"
	"time"

	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
)

type MatchService struct {
	repo Repository
}

func NewService(r Repository) *MatchService {
	return &MatchService{
		repo: r,
	}
}

func (s *MatchService) UpdateDownloadInformationForOutcomeId(id uint64, matchTime time.Time, url string) error {
	m, err := s.GetMatchByValveOutcomeId(id)
	if err != nil {
		return err
	}

	if m.DownloadURL != "" || m.Status != Created {
		return errors.New("download information already exists")
	}

	m.Status = Downloadable
	m.Time = matchTime
	m.DownloadURL = url

	if err := m.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateDownloadInformation(m)
}

func (s *MatchService) GetMatch(id entity.ID) (*Match, error) {
	return s.repo.Find(id)
}

func (s *MatchService) GetMatchByFilename(filename string) (*Match, error) {
	return s.repo.FindByFilename(filename)
}

func (s *MatchService) GetAll() ([]*Match, error) {
	return s.repo.List()
}

func (s *MatchService) GetAllParsed() ([]*Match, error) {
	return s.repo.ListParsedMatches()
}

func (s *MatchService) GetMatchByValveId(id uint64) (*Match, error) {
	return s.repo.FindByValveId(id)
}

func (s *MatchService) GetMatchByValveOutcomeId(id uint64) (*Match, error) {
	return s.repo.FindByValveOutcomeId(id)
}

func (s *MatchService) GetMatchByFaceitId(id string) (*Match, error) {
	return s.repo.FindByFaceitId(id)
}

func (s *MatchService) GetDownloadableMatches() ([]*Match, error) {
	return s.repo.ListDownloadableMatches()
}

func (s *MatchService) GetValveMatchesMissingDownloadUrl() ([]*Match, error) {
	return s.repo.ListValveMatchesMissingDownloadUrl()
}

func (s *MatchService) SetStatusAndFilename(m *Match, st Status, f string) error {
	m.Status = st
	m.Filename = f

	if err := m.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateStatusAndFilename(m)
}

func (s *MatchService) UpdateStatus(m *Match, st Status) error {
	m.Status = st

	if err := m.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateStatus(m)
}

func (s *MatchService) CreateDownloadableMatchFromFaceitId(faceitMatchId string, downloadUrl string, startTime time.Time) (*Match, error) {
	dbMatch, err := s.GetMatchByFaceitId(faceitMatchId)
	if err != nil && !errors.Is(err, entity.ErrNotFound) || dbMatch != nil {
		return dbMatch, err
	}

	m, _ := NewMatch(Faceit)
	m.FaceitMatchId = faceitMatchId
	m.DownloadURL = downloadUrl
	m.Status = Downloadable
	m.Time = startTime
	return m, s.repo.Create(m)
}

func (s *MatchService) CreateMatchFromSharecode(sc *share_code.ShareCodeData) (*Match, error) {
	dbMatch, err := s.GetMatchByValveId(sc.MatchID)
	if err != nil && !errors.Is(err, entity.ErrNotFound) || dbMatch != nil {
		return dbMatch, err
	}

	m, _ := NewMatch(MatchMaking)
	m.ShareCode = sc
	return m, s.repo.Create(m)
}

func (s *MatchService) CreateMatchFromManualUpload(filename string, matchTime time.Time) (*Match, error) {
	dbMatch, err := s.GetMatchByFilename(filename)
	if err != nil && !errors.Is(err, entity.ErrNotFound) || dbMatch != nil {
		return nil, nil
	}

	m, _ := NewMatch(Manual)
	m.Filename = filename
	m.Time = matchTime
	m.Status = Downloaded
	return m, s.repo.Create(m)
}

func (s *MatchService) GetParseableMatches(parserVersion byte) ([]*Match, error) {
	downloaded, errD := s.repo.ListDownloadedMatches()
	if errD != nil {
		return nil, errD
	}

	parsed, errP := s.repo.ListParsedMatches()
	if errP != nil {
		return nil, errP
	}

	parseable := make([]*Match, len(downloaded))
	copy(parseable, downloaded)

	for _, match := range parsed {
		if match.Result.ParserVersion < parserVersion {
			parseable = append(parseable, match)
		}
	}

	return parseable, nil
}

func (s *MatchService) UpdateResult(m *Match, r *MatchResult, parserVersion byte) error {
	m.Result = r
	m.Result.ParserVersion = parserVersion
	err := s.repo.UpdateResult(m)
	if err != nil {
		return err
	}

	if m.Status != Parsed {
		return s.UpdateStatus(m, Parsed)
	}

	return nil
}
