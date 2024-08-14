# Демонстрационный сервис, отображающий данные о заказах


Сборка с базой данных:
``docker run --name delivery_service -p 5434:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=delivery_service -d postgres:14``