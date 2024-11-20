package processor

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/database"
	"github.com/xochilpili/processor-films/internal/models"
	"github.com/xochilpili/processor-films/internal/services"
)

type ApiService interface {
	FetchTorrents(ctx context.Context, term string) ([]models.TorrentItem, error)
	AddTorrent(ctx context.Context, magnetLink string) (error)
}


type DatabaseService interface {
	Connect() error
	Close() error
	GetFilms(ctx context.Context, table string, columns []string, provider string) ([]models.FilmItem, error)
}

type Processor struct {
	config       *config.Config
	logger       *zerolog.Logger
	dbService    DatabaseService
	apiService   ApiService

}

func New(config *config.Config, logger *zerolog.Logger) *Processor {
	apiService := services.NewApi(config, logger)
	db := database.New(config, logger)
	return &Processor{
		config:       config,
		logger:       logger,
		dbService:    db,
		apiService:   apiService,

	}
}

func (p *Processor) Process(ctx context.Context, opType models.OperationType, provider string) error {
	err := p.dbService.Connect()
	if err != nil {
		p.logger.Fatal().Err(err).Msg("error while connecting to db")
		return err
	}
	defer p.dbService.Close()

	films, err := p.dbService.GetFilms(ctx, opType.String(), []string{"id", "provider", "title", "year"}, provider)
	if err != nil {
		p.logger.Fatal().Err(err).Msgf("error while getting %s from db", opType.String())
		return err
	}
	
	for _, film := range films {
		p.logger.Info().Msgf("processing film: %s", film.Title)
		tor, err := p.apiService.FetchTorrents(ctx, film.Title)
		if err != nil {
			p.logger.Fatal().Err(err).Msgf("error while getting torrents for: %s", film.Title)
			return err
		}
		if len(tor) == 0 {
			p.logger.Info().Msgf("no torrents found for: %s", film.Title)
			continue
		}
	}
	return nil
}
