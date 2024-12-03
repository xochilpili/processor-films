package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/models"
)

type Database struct {
	config *config.Config
	logger *zerolog.Logger
	db     *sql.DB
}

func New(config *config.Config, logger *zerolog.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Connect() error {
	psqlConn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", d.config.Database.Host, d.config.Database.Username, d.config.Database.Password, d.config.Database.Name)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return err
	}
	d.db = db
	err = d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (p *Database) Ping() error {
	return p.db.Ping()
}

func (p *Database) Close() error {
	err := p.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Database) GetOlderFilms(ctx context.Context, table string) ([]models.FilmItem, error) {
	var sqlStmt string = fmt.Sprintf("select (id, provider, title, year, genres, created_at, processed_at) from %s where processed_at between current_timestamp + interval '30 days' and current_timestamp + interval '15 days' and processed = 0 limit 10", table)
	rows, err := p.db.QueryContext(ctx, sqlStmt)
	if err != nil {
		return nil, err
	}
	var films []models.FilmItem
	for rows.Next() {
		film := models.FilmItem{}
		if err := rows.Scan(&film.Id, &film.Provider, &film.Title, &film.Year, &film.Genres); err != nil {
			p.logger.Err(err).Msg("error while fetching film from databse")
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (p *Database) GetFilms(ctx context.Context, table string, columns []string, provider string) ([]models.FilmItem, error) {
	cols := strings.Join(columns, ",")
	var sqlStmt string = fmt.Sprintf(`select %s from %s where processed = 0 and processed_at is null limit 10`, cols, table)
	if provider != "all" && provider != "" {
		sqlStmt = fmt.Sprintf(`select %s from %s where provider='%s'`, cols, table, provider)
	}
	rows, err := p.db.QueryContext(ctx, sqlStmt)
	if err != nil {
		return nil, err
	}
	var films []models.FilmItem
	for rows.Next() {
		film := models.FilmItem{}
		if err := rows.Scan(&film.Id, &film.Provider, &film.Title, &film.Year); err != nil {
			p.logger.Err(err).Msg("error while fetching film from databse")
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (p *Database) UpdateProcess(ctx context.Context, table string, id int) {
	var sqlStmt string = fmt.Sprintf("update %s set processed_at = current_timestamp where id = $1", table)
	_, err := p.db.ExecContext(ctx, sqlStmt, id)
	if err != nil {
		p.logger.Err(err).Msgf("error while updating processed time for id: %d", id)
	}
}

func (p *Database) ProcessedFilm(ctx context.Context, table string, id int) {
	var sqlStmt string = fmt.Sprintf("update %s set processed = 1 where id = $1", table)
	_, err := p.db.ExecContext(ctx, sqlStmt, id)
	if err != nil {
		p.logger.Err(err).Msgf("error while deleting film id: %d", id)
	}
}
