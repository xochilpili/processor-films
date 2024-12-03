package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/models"
)

type Api struct {
	config *config.Config
	logger *zerolog.Logger
	r      *resty.Client
}

func NewApi(config *config.Config, logger *zerolog.Logger) *Api {
	r := resty.New()
	return &Api{
		config: config,
		logger: logger,
		r:      r,
	}
}

func (a *Api) FetchTorrents(ctx context.Context, params models.FilterParams) ([]models.Torrent, error) {
	var result models.GenericResponse[models.Torrent]
	url := fmt.Sprintf("%s/%s/", a.config.TorrentApiUrl, params.Provider)
	a.logger.Info().Msgf("requesting torrents for %s to %s", params.Term, url)

	queryParams := map[string]string{
		"term": params.Term,
		"res":  params.Resolution,
	}

	res, err := a.r.R().SetHeader("Content-Type", "application/json").SetQueryParams(queryParams).SetDebug(a.config.Debug).Get(url)
	if err != nil {
		a.logger.Err(err).Msgf("error ocurred while fetching torrents for: %s, resolution: %s", params.Term, params.Resolution)
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		a.logger.Err(err).Msgf("error ocurred while decode json response for torrents: %s, resolution: %s", params.Term, params.Resolution)
		return nil, err
	}
	return result.Data, nil
}

func (a *Api) AddTorrent(ctx context.Context, magnetLink string) error {
	url := fmt.Sprintf("%s/api/v2/torrents/add", a.config.TransmissionApiUrl)
	_, err := a.r.R().SetFormData(map[string]string{
		"urls": magnetLink,
	}).Post(url)
	if err != nil {
		a.logger.Err(err).Msg("error while adding new torrent from magnet")
		return err
	}
	return nil
}

func (a *Api) GetSubtitles(ctx context.Context, title string) ([]models.Subtitle, error) {
	var result models.GenericResponse[models.Subtitle]
	a.logger.Info().Msgf("requesting subtitles for %s to %s", title, a.config.SubtitlerApiUrl)
	res, err := a.r.R().SetHeader("Content-Type", "application/json").SetQueryParam("term", title).SetDebug(a.config.Debug).Get(a.config.SubtitlerApiUrl)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (a *Api) GetTorrentMetadata(ctx context.Context, torrent *models.Torrent) (*models.TorrentMetadata, error) {
	var result models.TorrentMetadata
	a.logger.Info().Msgf("fetching torrent metadata for %s to %s", torrent.Title, a.config.TorrentMetadataApiUrl)
	res, err := a.r.R().
		SetHeader("Content-Type", "application/json").
		SetDebug(a.config.Debug).
		SetBody(map[string]interface{}{
			"query": torrent.Magnet,
		}).
		Post(a.config.TorrentMetadataApiUrl)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
