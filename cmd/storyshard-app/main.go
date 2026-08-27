package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"personae-fasti/api"
	"personae-fasti/configs"
	"personae-fasti/internal/handler"
	bunrepo "personae-fasti/internal/repo/bun"
	"personae-fasti/internal/service"
	objectstorage "personae-fasti/internal/storage"
)

func main() {
	config := configs.InitConfig()

	storage, err := bunrepo.NewBunStorage(config.DB)
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	defer storage.Close()

	migrationOnly := len(os.Args) > 1 && os.Args[1] == "migrate"
	production := strings.EqualFold(config.App.Environment, "production")
	if migrationOnly || config.App.MigrateOnStart || !production {
		if err := storage.Migrate(context.Background()); err != nil {
			log.Fatalf("database migration failed: %v", err)
		}
	}
	if migrationOnly {
		log.Println("database migrations applied successfully")
		return
	}
	playerRepo := storage.PlayerRepo()
	gameRepo := storage.GameRepo()
	recordRepo := storage.RecordRepo()
	entitiesRepo := storage.EntitiesRepo()
	questRepo := storage.QuestRepo()
	logRepo := storage.LogRepo()
	appRepo := storage.AppRepo()
	imageRepo := storage.ImageRepo()

	logService := service.NewLogService(logRepo)
	authService := service.NewAuthService(config.Auth, playerRepo)
	gameService := service.NewGameService(playerRepo, gameRepo, recordRepo)
	playerService := service.NewPlayerService(playerRepo, gameRepo)
	recordService := service.NewRecordService(playerRepo, gameRepo, recordRepo, questRepo)
	entitiesService := service.NewEntitiesService(entitiesRepo, recordRepo)
	questService := service.NewQuestService(questRepo, recordRepo)
	appService := service.NewAppService(appRepo)
	imageStorage := objectstorage.NewImageStorage(config.ImageStorage)
	imageService := service.NewImageService(imageRepo, entitiesRepo, imageStorage)

	authHandler := handler.NewAuthHandler(authService)
	gameHandler := handler.NewGameHandler(gameService)
	playerHandler := handler.NewPlayerHandler(playerService)
	recordHandler := handler.NewRecordHandler(recordService)
	entitiesHandler := handler.NewEntitiesHandler(entitiesService)
	questHandler := handler.NewQuestHandler(questService)
	appHandler := handler.NewAppHandler(appService)
	imageHandler := handler.NewImageHandler(imageService)

	privateApi := api.ConfigServer(
		config,
		authService,
		logService,
		storage.Ping,
	)

	privateApi.SetHandlers(authHandler, gameHandler, playerHandler, recordHandler, entitiesHandler, questHandler, appHandler, imageHandler)

	log.Println("API server running on ", privateApi.Server.Addr)
	if err := privateApi.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("API server stopped: %v", err)
	}
}
