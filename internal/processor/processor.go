package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/database"
	"github.com/xochilpili/processor-films/internal/models"
	"github.com/xochilpili/processor-films/internal/services"
	"github.com/xochilpili/processor-films/internal/utils"
)

type ApiService interface {
	FetchTorrents(ctx context.Context, params models.FilterParams) ([]models.Torrent, error)
	AddTorrent(ctx context.Context, magnetLink string) error
	GetSubtitles(ctx context.Context, title string) ([]models.Subtitle, error)
	GetTorrentMetadata(ctx context.Context, torrent *models.Torrent) (*models.TorrentMetadata, error)
}

type DatabaseService interface {
	Connect() error
	Close() error
	GetFilms(ctx context.Context, table string, columns []string, provider string) ([]models.FilmItem, error)
	GetOlderFilms(ctx context.Context, table string) ([]models.FilmItem, error)
	UpdateProcess(ctx context.Context, table string, id int)
	ProcessedFilm(ctx context.Context, table string, id int)
}

type Processor struct {
	config     *config.Config
	logger     *zerolog.Logger
	dbService  DatabaseService
	apiService ApiService
}

func New(config *config.Config, logger *zerolog.Logger) *Processor {
	apiService := services.NewApi(config, logger)
	db := database.New(config, logger)
	return &Processor{
		config:     config,
		logger:     logger,
		dbService:  db,
		apiService: apiService,
	}
}

func (p *Processor) Process(ctx context.Context, opType models.OperationType, provider string) error {
	err := p.dbService.Connect()
	if err != nil {
		p.logger.Fatal().Err(err).Msg("error while connecting to db")
		return err
	}
	defer p.dbService.Close()

	films, err := p.dbService.GetOlderFilms(ctx, opType.String())
	if err != nil {
		p.logger.Fatal().Err(err).Msgf("error while getting older %s films from db", opType.String())
		return err
	}

	if len(films) == 0 {
		films, err = p.dbService.GetFilms(ctx, opType.String(), []string{"id", "provider", "title", "year"}, provider)
		if err != nil {
			p.logger.Fatal().Err(err).Msgf("error while getting all %s films from db", opType.String())
			return err
		}
	}

	for _, film := range films {

		if film.Provider == "yts" {
			provider = film.Provider
		} else {
			provider = "all"
		}

		var title string
		if opType.EnumIndex() == 1 {
			title = fmt.Sprintf("%s %d", film.Title, film.Year)
		} else {
			title = film.Title
		}

		p.logger.Info().Msgf("processing film: %s, type: %s", title, opType.String())

		torrentItems, err := p.apiService.FetchTorrents(ctx, models.FilterParams{Provider: provider, Term: title, Resolution: "720p"})
		if err != nil {
			p.logger.Fatal().Err(err).Msgf("error while getting torrents for: %s", title)
			return err
		}

		if len(torrentItems) == 0 {
			p.logger.Info().Msgf("no torrents found for: %s", title)
			p.dbService.UpdateProcess(ctx, opType.String(), film.Id)
			continue
		}

		torrent, strFile, spa := p.hasFileSubtitles(torrentItems)
		if spa {
			// torrent has spanish subtitles
			p.apiService.AddTorrent(ctx, torrent.Magnet)
			p.dbService.ProcessedFilm(ctx, opType.String(), film.Id)
			continue
		}

		subs, err := p.apiService.GetSubtitles(ctx, title)
		if err != nil {
			p.logger.Err(err).Msgf("error while fetching subtitles for %s", title)
			continue
		}

		if len(subs) == 0 {
			if strFile {
				p.apiService.AddTorrent(ctx, torrent.Magnet)
				p.dbService.ProcessedFilm(ctx, opType.String(), film.Id)
				p.logger.Info().Msgf("torrent %s added with file subtitles", torrent.Title)
			}
			// no torrent files and no subtitles
			p.logger.Info().Msgf("no subtitles found for %s", title)
			continue
		}

		subTorrent := p.matchSubtitles(torrentItems, subs)

		if subTorrent == nil {
			// add torrent from file and continue
			p.logger.Info().Msgf("no online subtitles matches for %s", title)
			if strFile {
				p.apiService.AddTorrent(ctx, torrent.Magnet)
				p.dbService.ProcessedFilm(ctx, opType.String(), film.Id)
				p.logger.Info().Msgf("torrent %s added with file subtitles", torrent.Title)
			}
			continue
		}

		// add sub-torrent and continue
		p.apiService.AddTorrent(ctx, subTorrent.Magnet)
		p.dbService.ProcessedFilm(ctx, opType.String(), film.Id)
		p.logger.Info().Msgf("torrent %s added with matched online subtitles", subTorrent.Title)

		// for debug proposes
		if p.config.Debug {
			fmt.Printf("film: %s\n", title)
			out, _ := json.MarshalIndent(torrentItems, "", "\t")
			fmt.Println(string(out))
		}
	}
	p.logger.Info().Msgf("processed %d items", len(films))
	return nil
}

func (p *Processor) matchSubtitles(torrents []models.Torrent, subs []models.Subtitle) *models.Torrent {
	var partialMatch bool
	var torr *models.Torrent
	for _, torrent := range torrents {
		ok := utils.Some(subs, func(s models.Subtitle) bool {
			resOk := utils.Some(s.Resolution, func(res string) bool {
				if res != "" && res == torrent.Resolution {
					return true
				}
				return false
			})
			qOk := utils.Some(s.Quality, func(qa string) bool {
				if qa != "" && qa == torrent.Quality {
					return true
				}
				return false
			})
			gOk := utils.Some(s.Group, func(g string) bool {
				// subtitles are based on whatever user's add, some of them uses yify instead of yts
				if g == "yify" && torrent.Group == "yts" {
					return true
				}
				if g != "" && g == torrent.Group {
					return true
				}
				if strings.Contains(strings.ToLower(g), strings.ToLower(torrent.Group)) {
					return true
				}
				return false
			})
			if resOk && qOk && gOk {
				p.logger.Info().Msgf("perfect matched %s", torrent.Title)
				return true
			}
			if (resOk && qOk) || (resOk && gOk) || (qOk && gOk) {
				p.logger.Warn().Msgf("partial match for %s, %t, %t, %t", torrent.Title, resOk, qOk, gOk)
				partialMatch = true
				torr = &torrent
				return true
			}
			return false
		})
		if partialMatch {
			return torr
		}
		if ok && !partialMatch {
			return &torrent
		}
	}
	return nil
}

func (p *Processor) hasFileSubtitles(torrents []models.Torrent) (*models.Torrent, bool, bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, torrent := range torrents {
		metadata, err := p.apiService.GetTorrentMetadata(ctx, &torrent)
		if err != nil {
			p.logger.Err(err).Msgf("error while receiving metadata for torrent: %s", torrent.Title)
			return nil, false, false
		}

		p.logger.Info().Msgf("%d total files found in torrent's metadata: %s", len(metadata.Data.Files), torrent.Title)

		for _, file := range metadata.Data.Files {
			extension := path.Ext(file.Name)
			if extension == ".srt" || extension == ".ass" {
				if strings.Contains("spa", file.Name) || strings.Contains("latin", file.Name) {
					return &torrent, true, true
				}
				return &torrent, true, false
			}
		}
	}
	return nil, false, false
}
