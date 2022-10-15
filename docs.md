# Документация

## Создание пользователя
Отправляем запрос на создание пользователя:
```
curl -X 'POST' \
  'http://localhost:8080/api/users' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "Data": {
    "FirstName": "Ivan",
    "LastName": "Ivanov",
    "Email": "ivanov@mail.ru",
    "Password": "12345",
    "TimeZoneOffset": 3
  }
}'
``` 

Получаем ответ вида:
```
{"ID":1} - id нового пользователя
```
## Аутентификация пользователя
```
curl -X 'POST' \
'http://localhost:8080/api/login' \
-H 'accept: application/json' \
-H 'Content-Type: application/x-www-form-urlencoded' \
-d 'username=ivanov%40mail.ru&password=12345'
```

Получаем ответ вида:
```
{"Data":{"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM"}}
```
Далее передаем полученный токен в заголовке - **Authorization**

## Создание встречи
```
curl -X 'POST' \
  'http://localhost:8080/api/events' \
  -H 'Authorization:Barrier eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "Data": {
    "From": "2022-10-12T15:00:00Z",
    "Till": "2022-10-12T16:30:00Z",
    "CreatorID": 1,
    "Participants": [
      2
    ],
    "Details": "details",
    "IsPrivate": true,
    "IsRepeat": false
  }
}'
```

Получаем ответ вида:
```
{"ID":1} - id нового эвента
```

Создание встречи с правилом повторения встречи:
```
curl -X 'POST' \
  'http://localhost:8080/api/events' \
  -H 'Authorization:Barrier eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "Data": {
    "From": "2022-10-15T15:00:00Z",
    "Till": "2022-10-15T16:30:00Z",
    "CreatorID": 1,
    "Participants": [
      2
    ],
    "Details": "details",
    "ScheduleRule": "SCHEDULER_MODE=CUSTOM;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=WEEKLY;CUSTOM_DAY_LIST=1,2,3",
    "IsPrivate": false,
    "IsRepeat": true
  }
}'
```

Получаем ответ вида:
```
{"ID":2} - id нового эвента
```

## Получение информации встречи по ID
```
curl -X 'GET' \
  -H 'Authorization:Barrier eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM' \
  'http://localhost:8080/api/events/1' \
  -H 'accept: application/json'
```

Получаем ответ вида:
```
{"Data":{"CreatedAt":"2022-10-15 18:39:09.187022 +0000 UTC","Creator":{"CreatedAt":"2022-10-15 18:23:42.134989 +0000 UTC","Email":"ivanov@mail.ru","FirstName":"Ivan","ID":1,"LastName":"Ivanov","TimeZoneOffset":3,"UpdatedAt":"2022-10-15 18:23:42.134989 +0000 UTC"},"Details":"details","From":"2022-10-12 15:00:00 +0000 UTC","ID":1,"IsPrivate":true,"IsRepeat":false,"Participants":[{"CreatedAt":"2022-10-15 18:32:58.810394 +0000 UTC","Email":"petrov@mail.ru","FirstName":"Petr","ID":2,"LastName":"Petrov","TimeZoneOffset":3,"UpdatedAt":"2022-10-15 18:32:58.810394 +0000 UTC"}],"Till":"2022-10-12 16:30:00 +0000 UTC","UpdatedAt":"2022-10-15 18:39:09.219167 +0000 UTC"}}
```

## Подтверждение присутствия на встрече
```
curl -X 'PUT' \
  'http://localhost:8080/api/invites/1' \
  -H 'Authorization:Barrier eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "Data": {
    "IsAccepted": true
  }
}'
```

Получаем ответ вида:
```
{"ID":1} # id обновленной записи 
```

## Получение эвентов пользователя в заданном интервале
```
curl -X 'GET' \
  'http://localhost:8080/api/users/1/events?from=2022-10-14T11%3A00%3A00Z&till=2022-10-20T11%3A00%3A00Z' \
  -H 'accept: application/json' \
  -H 'Authorization:Barrier eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM'
```

В результате будет возвращен список встреч пользователя. Те встречи, что он создал или те, инвайт на которые принял.

## Получение ближайшего свободного окна для нескольких пользователей

```
curl -X 'GET' \
  'http://localhost:8080/api/windows?user_ids%5B%5D=1&user_ids%5B%5D=2&win_size=15m' \
  -H 'accept: application/json' \
  -H 'Authorization:Barrier eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoiaXZhbm92QG1haWwucnUifQ.evXbemYej5uqiiMygZ3nbK-iD7ow3dM8kojbaiVrcJM
```

В результате будет возвращены отметки начала и конца ближайшего свободного окна.