# Демонстрационный сервис, отображающий данные о заказах


Запуск веб-сервиса:
``docker compose up -d``

Накат миграций:
``goose postgres "postgresql://admin:admin@127.0.0.1:5434/delivery_service?sslmode=disable" up``

