# Product Management

Сервис для управления складом продуктов

## Архитектура 
В проекте применялся принцип Чистой архитектуры (Clean Architecture) с четким разделением на слои зоны ответственности:
   - Repository:работа, действия с БД
   - Usecase: бизнес логика сервиса
   - Rest: Обработчики запросов

## Технологии
- Golang
- Postgres
- Docker
- Swagger
- dig (Dependency Injection, был использован для 
  автоматизации, упрощения внедрения зависимостей, 
  поддержки масшабирования, 
  сохранение лаконичности и простоты кода во избежание "Матрешки" ручного свзывания)
- gin (web framework)
- golang-migrate (DB migration)
- godotenv (env variables)
- Linter (анализ кода на ошибки, стиль и перформанс)

## Запуск

1. Поднять Docker контейнер в директории deployment
```bash
cd deployment && docker-compose up -d
```
(P.S при желании можно поменять настройки docker-compose под свои нужды, но потом учесть их в .env!)

2. Создаем файл .env в корневой дериктории проекта. В корневой папке имеется шаблон .env.example для настройки окружения переменных.
Копируем и вставляем все поля в .env, дополняем поля host, port, dsn(если что то меняли в docker-compose.yml) file оставляем как есть т.к там путь к миграциям

3. После всех настроек, запускаем проект
```bash
go run cmd/app/main.go
```
4. Далее переходим в браузер, в строке поиска вводим 
```http://{host}:{port}/swagger/index.html``` 
где host и port исходя из своей конфиграции env. Пример: 
```http://localhost:8080/swagger/index.html```

5. End-Поинты
```http request
POST        /products // добавить товар
GET         /products // получить/смотреть все товары
GET         /products/:id // получить товар по id
PUT         /products/:id // изменить товар
DELETE      /products/:id // удалить/архивировать товар
PUT         /products/:id // восстановить товар

GET         /products/health // проверка работоспособности сервиса
```