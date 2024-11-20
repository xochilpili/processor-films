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

func (a *Api) FetchTorrents(ctx context.Context, term string) ([]models.TorrentItem, error) {
	var result models.GenericResponse[models.TorrentItem]
	res, err := a.r.R().SetHeader("Content-Type", "application/json").SetQueryParam("term", term).SetDebug(a.config.Debug).Get(a.config.TorrentApiUrl)
	if err != nil {
		a.logger.Err(err).Msgf("error ocurred while fetching torrents for: %s", term)
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		a.logger.Err(err).Msgf("error ocurred while decode json response for torrents: %s", term)
		return nil, err
	}
	return result.Data, nil
}

func(a *Api) AddTorrent(ctx context.Context, magnetLink string) (error){
	url := fmt.Sprintf("%s/api/v2/torrents/add", a.config.TransmissionApiUrl)
	_, err := a.r.R().SetFormData(map[string]string{
		"urls": magnetLink,
	}).Post(url)
	if err != nil{
		a.logger.Err(err).Msg("error while adding new torrent from magnet")
		return err
	}
	return nil
}



