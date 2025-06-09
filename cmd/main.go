package main

import (
	"flag"
	"log"
	"net/http"

	"ozon_test/config"
	"ozon_test/internal/api"
	"ozon_test/internal/services"
	"ozon_test/internal/storage"
)

func main() {
	storageType := flag.String("storage", "inmemory", "Storage type: inmemory or postgres")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	var postStorage storage.PostStorage
	var commentStorage storage.CommentStorage

	switch *storageType {
	case "inmemory":
		postStorage = storage.NewInMemoryPostStorage()
		commentStorage = storage.NewInMemoryCommentStorage()
	case "postgres":
		if err := storage.ApplyMigrations(cfg); err != nil {
			log.Fatalf("Ошибка применения миграций: %v", err)
		}
		pool, err := storage.CreateDBPool(cfg)
		if err != nil {
			log.Fatalf("Ошибка создания пула подключений: %v", err)
		}
		postStorage = storage.NewPostgresPostStorage(pool)
		commentStorage = storage.NewPostgresCommentStorage(pool)
	default:
		log.Fatal("Неизвестный тип хранилища")
	}

	postService := services.NewPostService(postStorage)
	commentService := services.NewCommentService(commentStorage, postStorage)

	postHandler := api.NewPostHandler(postService)
	commentHandler := api.NewCommentHandler(commentService)

	http.HandleFunc("/posts", postHandler.GetAllPosts)
	http.HandleFunc("/posts/create", postHandler.CreatePost)
	http.HandleFunc("/posts/disable-comments", postHandler.DisableComments)
	http.HandleFunc("/comments", commentHandler.GetComments)
	http.HandleFunc("/comments/create", commentHandler.CreateComment)

	serverAddr := ":8080"
	log.Printf("Сервер запускается на %s...", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
