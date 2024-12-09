# Онлайн библеотека песен

## Запуск
Необходимо задать натройки для работу приложения в .env-файле. 
По умолчанию используются дефолтные настройки для запуска программы. 

Для вызова API при добавлении новой песни необходио в env-файле указать
``` go
# --API--
API_CALL=true
API_HOST=необходимый хост
API_PORT=необходимый порт
```

Для запуска необходимо использовать команду 
``` go
docker-compose up
```

## Возмоности 
Получение данных библиотеки с фильтрацией по всем полям и пагинацией
Получение текста песни с пагинацией по куплетам
Удаление песни
Изменение данных песни
Добавление новой песни в формате

Для реализации данных функций и представления swagger-документации используются следующие маршруты в функции NewRouter: 
``` go 
func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.WithLogging(h.Log))

	r.Put("/songs/{id}", http.HandlerFunc(h.Update))
	r.Get("/songs", http.HandlerFunc(h.GetAll))
	r.Get("/songs/{song}", http.HandlerFunc(h.GetText))
	r.Post("/songs", http.HandlerFunc(h.AddSong))
	r.Delete("/songs/{id}", http.HandlerFunc(h.Delete))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
```

## Миграции
Для работы с миграциями исользуется пакет goose.
При необходимости для установки данного пакета используйте команду:
``` go
make install-deps
```

## Swagger
После запуска приложения при стандартных настройках в env-файле swagger-документация будет доступна по адресу http://localhost:8080/swagger/index.html
