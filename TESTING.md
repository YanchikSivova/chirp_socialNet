
### AccountService

- `utils/hash_test.go` проверяет хеширование пароля и bcrypt-проверку.
- `utils/token_test.go` проверяет генерацию access token, парсинг temporary token и формат verification code.

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\AccountService
go test -count=1 ./...
```

### MessageService

- `util/parseToken_test.go` проверяет парсинг JWT и невалидные токены.

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\MessageService
go test -count=1 ./...
```

### Gateway

- `gateway_e2e_test.go` проверяет проксирование авторизованного запроса через gateway с помощью `httptest`.
- `public_auth_e2e_test.go` проверяет публичный маршрут `/auth/login`: метод, путь, тело запроса и возврат ответа.
- `docker_account_e2e_test.go` по умолчанию пропускается и запускается только против Docker Compose.

Запуск локальных тестов:

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\Gateway
go test -count=1 ./...
```

### PostService

- `handler/parse_test.go` проверяет функции для `X-User-Id`, `:id`, `offset` и `limit`.

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\PostService
go test -count=1 ./...
```

### NotificationService

- `handler/parse_test.go` проверяет функции для `X-User-Id` и notification `:id`.

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\NotificationService
go test -count=1 ./...
```

### FeedService

- `handler/parse_test.go` проверяет парсинг feed `:id`.

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\FeedService
go test -count=1 ./...
```

### Docker E2E тест

`Gateway/docker_account_e2e_test.go`.

 Поднять Docker Compose

```powershell
cd C:\Users\Thunderobot\Desktop\socialNet\Gateway
$env:RUN_DOCKER_E2E="1"
go test -count=1 -run TestDockerE2ELoginAndGetCurrentUser -v
```
