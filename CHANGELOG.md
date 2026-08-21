## [2.0.5] - 2026-03-07

### ⚙️ Miscellaneous Tasks

- Update Go version and dependencies
- Update Go version on github actions
- Update CI workflow actions and matrix usage
## [2.0.4] - 2025-01-25

### 🐛 Bug Fixes

- Update SDK imports to v2
## [2.0.3] - 2025-01-22

### 🐛 Bug Fixes

- Update SDK import paths to v2
## [2.0.2] - 2025-01-18

### 🐛 Bug Fixes

- *(module)* Update module path to v2
## [2.0.1] - 2025-01-16

### 🐛 Bug Fixes

- Update branch filters for workflows
- Simplify workflow triggers

### ⚙️ Miscellaneous Tasks

- Update changelog workflow for better auth
- Simplify changelog update message
## [2.0.0] - 2025-01-16

### 🚀 Features

- Add new sync methods and improve API client
- Remove example environment file
- Add examples for ZenMoney SDK usage
- Update README.md
- Add suggestion functionality for transactions
- Add new dependencies and tests

### 🐛 Bug Fixes

- Update release args for GoReleaser

### 🚜 Refactor

- Update module name to zenmoney-go-sdk
- Rename main.go to client.go
- Rename package from zenapi to zenmoney
- [**breaking**] Restructure project and enhance API client
- Move diff request to single file
- *(tests)* Clean up test code and improve error handling

### 🧪 Testing

- Add unit tests for custom error handling

### ⚙️ Miscellaneous Tasks

- Update Go version to 1.23.4
- Update workflow actions to latest versions
## [1.3.3] - 2024-10-10

### 🚀 Features

- Added comments to data structures in system.go
- Added comments to data structures in user.go

### ⚙️ Miscellaneous Tasks

- Update config goreleaser
- Update go releaser
## [1.3.2] - 2024-05-11

### 🚀 Features

- Добавил 'deletion' в response

### ⚙️ Miscellaneous Tasks

- Добавил автоматическую генерацию changelog
## [1.3.1] - 2024-04-29

### 🚀 Features

- Функция SyncSince для получения данных с последней синхронизации
- Добавил сборку, тестыб релизы с помощью github actions
- Обновил версию golang
- Badge codecov
- Добавил конфигурацию для релизера golang

### 🐛 Bug Fixes

- Ошибка в Errorf
- Название токена github
## [1.2.1] - 2023-12-16

### 🚀 Features

- Добавлены новые поля в структуру Transaction
## [1.2.0] - 2023-12-16

### 💼 Other

- Убрал тег omitempty в структурах.
- Поменял тип цвета в структуре tag
## [1.1.0] - 2023-12-06

### 💼 Other

- Обновлен конструктор NewClient для принятия токена в качестве аргумента
## [1.0.0] - 2023-12-05
