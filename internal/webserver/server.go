package webserver

import (
	"context"
	"net/http"

	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/models"
	"github.com/xochilpili/processor-films/internal/processor"
)

type Processor interface {
	Process(ctx context.Context, opType models.OperationType, provider string) error
}

type WebServer struct {
	config    *config.Config
	logger    *zerolog.Logger
	Web       *http.Server
	ginger    *gin.Engine
	processor Processor
}

func New(config *config.Config, logger *zerolog.Logger) *WebServer {
	ginger := gin.New()
	ginger.Use(gin.Recovery())
	ginger.Use(ginlogger.SetLogger(
		ginlogger.WithSkipPath([]string{"/ping"}),
		ginlogger.WithLogger(func(ctx *gin.Context, l zerolog.Logger) zerolog.Logger {
			return logger.Output(gin.DefaultWriter).With().Logger()
		}),
	))
	httpSrv := &http.Server{
		Addr:    config.Host + ":" + config.Port,
		Handler: ginger,
	}
	processor := processor.New(config, logger)
	srv := &WebServer{
		config:    config,
		logger:    logger,
		Web:       httpSrv,
		ginger:    ginger,
		processor: processor,
	}

	srv.loadRoutes()

	return srv
}

func (w *WebServer) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{"messasge": "pong"})
}

func (w *WebServer) festivalHandler(c *gin.Context) {
	go w.processor.Process(context.Background(), models.FESTIVALS, "all")
	c.JSON(http.StatusOK, &gin.H{"message": "ok"})
}

func (w *WebServer) popularHandler(c *gin.Context) {
	go w.processor.Process(context.Background(), models.POPULAR, "all")
	c.JSON(http.StatusOK, &gin.H{"message": "ok"})
}

func (w *WebServer) loadRoutes() {
	api := w.ginger.Group("/")
	api.GET("/ping", w.pingHandler)
	process := w.ginger.Group("/process")
	{
		process.GET("/festivals", w.festivalHandler)
		process.GET("/popular", w.popularHandler)
	}
}
