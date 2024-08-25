# Демонстрационный сервис, отображающий данные о заказах

Запуск веб-сервиса:
```docker compose up --build```

Публикация сообщений в nats:
```docker exec -it delivery_service bash```

```cd testing && go run publish-script.go```

Веб-интерфейс:
http://localhost:8080

Завершение работы:
```docker compose stop```