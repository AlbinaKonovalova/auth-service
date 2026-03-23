# Auth Service

HTTP/JSON сервис аутентификации и авторизации для админки.

Текущий статус: MVP для фронта.

Сейчас сервис умеет:
- login
- refresh access token через httpOnly cookie
- logout
- current user (`/auth/me`)
- healthcheck
- swagger UI/spec
- seed initial admin

### Схема каталогов

```text
auth-service/
├── cmd/                                        # entrypoints
│   ├── auth-service/                           # main package сервиса
│   └── seed/                                   # seed initial admin
│
├── internal/
│   ├── adapter/                                # входные/выходные адаптеры
│   │   ├── controller/                         # proto/gateway controllers
│   │   ├── repository/                         # реализации output ports
│   │   │   ├── postgres/                       # postgres repositories
│   │   │   └── mail/                           # smtp mail adapter
│   │   ├── auth/                               # technical auth/security adapters: password hashing, jwt, token hash
│   │   └── system/                             # technical system adapters: clock, uuid and other local runtime providers
│   │
│   ├── config/                                 # config loader + runtime configs
│   │   └── modules/                            # отдельный модуль на каждую подсистему (argon2, auth, cookie, cors, db, log, smtp, ...)
│   │
│   ├── domain/                                 # доменное ядро
│   │   ├── entity/                             # сущности и агрегаты
│   │   ├── value/                              # value objects
│   │   ├── service/                            # domain services / policies / factories over multiple entities
│   │   └── dto/                                # доменные dto
│   │
│   ├── infrastructure/                         # frameworks & drivers
│   │   ├── app/                                # bootstrap / DI / lifecycle
│   │   ├── http/                               # http server + gateway + swagger
│   │   ├── postgres/                           # db connection / pooling
│   │   └── middleware/
│   │       └── http/                           # http middleware
│   │
│   ├── ports/                                  # входящие и исходящие контракты
│   │   ├── input/                              # usecase interfaces
│   │   └── output/                             # repository/service interfaces
│   │
│   └── usecase/                                # прикладные сценарии
│       ├── auth/                               # login / refresh / logout / me
│       ├── password_reset/                     # request / confirm reset
│       ├── user/                               # create / get / list / activate / deactivate
│       ├── access/                             # assignment usecases: user roles and role permissions
│       ├── role/                               # list / create / delete roles
│       ├── permission/                         # list / create / delete permissions
│       └── common/                             # shared application helpers, not place for domain logic
│
├── pkg/
│   ├── proto/                                  # исходные .proto контракты
│   │   └── authservice/
│   │       └── v1/
│   │
│   └── authservice/                            # generated code from proto
│       └── v1/                                 # *.pb.go, *.pb.gw.go, swagger.json
│
├── migrations/                                 # goose migrations
├── config/                                     # runtime yaml config
├── docker/                                     # docker files
├── buf.yaml                                    # buf lint/breaking config
├── buf.gen.yaml                                # buf generate config
├── buf.work.yaml                               # buf workspace config
├── Makefile                                    # команды разработки
├── go.mod                                      # go module
├── go.sum                                      # locked deps
└── README.md                                   # документация проекта
```
