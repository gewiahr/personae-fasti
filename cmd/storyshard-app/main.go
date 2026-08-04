package main

import (
	"log"
	"personae-fasti/api"
	"personae-fasti/configs"
	"personae-fasti/internal/handler"
	bunrepo "personae-fasti/internal/repo/bun"
	"personae-fasti/internal/service"
)

func main() {
	config := configs.InitConfig()

	storage, err := bunrepo.NewBunStorage(config.DB)
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	playerRepo := storage.PlayerRepo()
	gameRepo := storage.GameRepo()
	recordRepo := storage.RecordRepo()
	entitiesRepo := storage.EntitiesRepo()
	questRepo := storage.QuestRepo()
	logRepo := storage.LogRepo()
	appRepo := storage.AppRepo()

	logService := service.NewLogService(logRepo)
	authService := service.NewAuthService(config.Auth, logService, playerRepo)
	gameService := service.NewGameService(playerRepo, gameRepo, recordRepo)
	playerService := service.NewPlayerService(playerRepo, gameRepo)
	recordService := service.NewRecordService(playerRepo, gameRepo, recordRepo)
	entitiesService := service.NewEntitiesService(entitiesRepo, recordRepo)
	questService := service.NewQuestService(questRepo, recordRepo)
	appService := service.NewAppService(appRepo)

	authHandler := handler.NewAuthHandler(authService)
	gameHandler := handler.NewGameHandler(gameService)
	playerHandler := handler.NewPlayerHandler(playerService)
	recordHandler := handler.NewRecordHandler(recordService)
	entitiesHandler := handler.NewEntitiesHandler(entitiesService)
	questHandler := handler.NewQuestHandler(questService)
	appHandler := handler.NewAppHandler(appService)

	privateApi := api.ConfigServer(
		config,
		authService,
		logService,
	)

	privateApi.SetHandlers(authHandler, gameHandler, playerHandler, recordHandler, entitiesHandler, questHandler, appHandler)

	log.Println("API server running on ", privateApi.Server.Addr)
	if err := privateApi.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}
