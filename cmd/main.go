package main

import (
	"context"

	"github.com/xochilpili/processor-films/internal/config"
	"github.com/xochilpili/processor-films/internal/logger"
	"github.com/xochilpili/processor-films/internal/models"
	"github.com/xochilpili/processor-films/internal/processor"
)

func main() {
	config := config.New()
	logger := logger.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor := processor.New(config, logger)
	processor.Process(ctx, models.POPULAR, "all") 

}
