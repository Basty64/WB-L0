# Демонстрационный сервис, отображающий данные о заказах


Сборка с базой данных:
``docker run --name delivery_service -p 5434:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=delivery_service -d postgres:14``

Накат миграций:
``goose postgres "postgresql://admin:admin@127.0.0.1:5434/delivery_service?sslmode=disable" up``