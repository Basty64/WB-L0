# Delivery-service
## Демонстрационный сервис, отображающий данные о заказах


Реализованы

1. Модели хранения данных по структуре из [файла](https://github.com/Basty64/delivery-service/blob/main/docs/model.json);
2. Веб интерфейс для доступа к заказам;
2. Хранение заказов в бд;
3. Кэширование заказов in memory;
4. Добавление и обработка новых заказов через nats-streaming;
5. Валидация входящих данных;
6. Деплой через docker и docker compose.

### Запуск:

При первом запуске:
```
docker compose up postgres
```
Затем:
```
docker compose up --build
```
> Использовался docker compose второй версии без retries и поэтому при первом запуске миграции не успевают накатываться до полного запуска постгреса.

При последующих запусках:
```
docker compose up
```

Публикация сообщений в nats:

Если подключаться из терминала:
```
docker exec -it delivery_service bash
```
Затем:
```
cd testing && go run publish-script.go
```

> Если подключаться к контейнеру из десктопной версии докера, то первая команда не нужна.

Завершение работы:
```
docker compose stop
```
Веб-интерфейс:
http://localhost:8080

### Пример изображения
![Пример изображения](docs/example.png)

### Структура БД
![Структура БД](docs/bd-diagram.png)

Используемый стэк:
* Бэк -
go-1.22.6 || postgres 14 || nats-streaming
* Фронт - html || css || java-script

## Результаты тестирования с помощью утилиты vegeta:

В первом случае запросы отправляются на форму выдачи заказов, во втором без неё берется конкретный заказ по айди.

### Команда с включённым фронтендом:
```
echo "GET http://localhost:8080/" | vegeta attack -duration=60s -rate=1000 | tee results.bin | vegeta report
```
### 1000 rps
```
Requests      [total, rate, throughput]         60000, 1000.02, 1000.01
Duration      [total, attack, wait]             59.999s, 59.999s, 521.042µs
Latencies     [min, mean, 50, 90, 95, 99, max]  265.5µs, 650.343µs, 603.869µs, 768.625µs, 940.623µs, 1.894ms, 10.895ms
Bytes In      [total, mean]                     184200000, 3070.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:60000
```
### 10000 rps
```
Requests      [total, rate, throughput]         300000, 10000.16, 1191.24
Duration      [total, attack, wait]             59.112s, 30s, 29.112s
Latencies     [min, mean, 50, 90, 95, 99, max]  6.708µs, 904.718ms, 16.462µs, 602.986ms, 4.836s, 26.276s, 30.028s
Bytes In      [total, mean]                     216177120, 720.59
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           23.47%
Status Codes  [code:count]                      0:229584  200:70416 
```

### Тест с использованием фронта:

* Низкая успешность (23.47%): Это самая большая проблема. Сервис обрабатывает только 23.47% запросов успешно, а 76.53% заканчиваются ошибкой.
* Высокая задержка: Средняя задержка (904.718мс) очень большая, а 99-й процентиль (26.276 секунд) говорит о том, что в 1% случаев запросы заканчиваются с задержкой более 26 секунд.
* Малый объём входящих данных: Сервис получает незначительный объем входящих данных (216177120 байт), что может указывать на то, что фронт не передает серверу все необходимые данные.

### С отключённым:
```
echo "GET http://localhost:8080/order?orderUID=3" | vegeta attack -duration=60s -rate=1000 | tee results.bin | vegeta report
```
### 1000 rps
```
Requests      [total, rate, throughput]         60000, 1000.02, 1000.01
Duration      [total, attack, wait]             59.999s, 59.999s, 244.584µs
Latencies     [min, mean, 50, 90, 95, 99, max]  82.834µs, 1.023ms, 218.234µs, 286.806µs, 331.034µs, 10.21ms, 190.048ms
Bytes In      [total, mean]                     61020000, 1017.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:60000  
```
### 10000 rps
```
Requests      [total, rate, throughput]         597802, 9947.54, 7776.87
Duration      [total, attack, wait]             1m2s, 1m0s, 1.505s
Latencies     [min, mean, 50, 90, 95, 99, max]  50.583µs, 237.457ms, 24.386ms, 792.16ms, 1.099s, 1.89s, 2.978s
Bytes In      [total, mean]                     487200969, 814.99
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           80.14%
Status Codes  [code:count]                      0:118745  200:479057 
```
### Краткие выводы
* Высокая успешность (80.14%): Сервис обрабатывает 80.14% запросов успешно, что значительно лучше, чем при использовании фронта.
* Сравнительно низкая задержка: Средняя задержка (237.457мс) ниже, чем при использовании фронта, а 99-й процентиль (1.89 секунд) значительно меньше.
* Большой объем входящих данных: Сервис получает больше входящих данных (487200969 байт), что может указывать на то, что фронт не передает серверу все необходимые данные или передает их неэффективно.

## Результаты тестирования с помощью утилиты WRK:

```
 10 threads and 10000 connections
 
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    24.16ms   37.92ms   1.66s    97.56%
    Req/Sec    10.03k     3.67k   28.93k    76.21%
    
 5873067 requests in 1.00m, 0.96GB read
  
 Socket errors: connect 7455, read 45393, write 1, timeout 0
  
 Non-2xx or 3xx responses: 5873067
  
Requests/sec:  97730.12
Transfer/sec:     16.40MB
```
### Общие выводы
Фронт является узким местом. Результаты показывают, что фронт неэффективно обрабатывает запросы и передает данные серверу. Он является причиной низкой успешности, высокой задержки и недостаточной передачи данных.
Так как целью является показать именно бэкенд часть, то проблемы при использовании фронта не являются критичными на данном этапе. Планку в 1000 запросов уже можно было считать хорошим результатом, но этот рубеж удалось превзойти.
Сервис получился вполне рабочим и при некоторых доработках в будущем может быть значительно увеличен в производительности.

В будущем тесты могут быть скорректированы.
