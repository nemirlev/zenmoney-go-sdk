## [2.0.5] - 2026-03-07

### ⚙️ Miscellaneous Tasks

- Update Go version and dependencies
- Update Go version on github actions
- Update CI workflow actions and matrix usage
## [3.0.1](https://github.com/nemirlev/zenmoney-go-sdk/compare/v3.0.0...v3.0.1) (2026-08-22)


### 🐛 Bug Fixes

* **module:** update module path to v3 ([fff1418](https://github.com/nemirlev/zenmoney-go-sdk/commit/fff141892bc2b00573d3f4f01cb086e4893cc182))


### 📚 Documentation

* replace Go Report Card badge ([1852bd5](https://github.com/nemirlev/zenmoney-go-sdk/commit/1852bd53dc4937e47cc7fdc6879c30a8b79b4412))

## [3.0.0](https://github.com/nemirlev/zenmoney-go-sdk/compare/v2.0.7...v3.0.0) (2026-08-22)


### ⚠ BREAKING CHANGES

* v2.0.7 changed Unix timestamps from int to int64 and nullable string fields to pointers.

### 🚀 Features

* add country entity type ([0f784e5](https://github.com/nemirlev/zenmoney-go-sdk/commit/0f784e519ba08ef76157d1b7ffe5480d6538fd0d))
* add optional structured diagnostics ([bd50c3b](https://github.com/nemirlev/zenmoney-go-sdk/commit/bd50c3b38468b90edcd45196e6a4403ff40cf728))
* add structured SDK diagnostics ([5bd9a99](https://github.com/nemirlev/zenmoney-go-sdk/commit/5bd9a993da303db4a9c2ddd0967adaaa7df9698a))
* expose bounded HTTP error details ([45e9215](https://github.com/nemirlev/zenmoney-go-sdk/commit/45e9215b9dfd0437f1b27db43386c417833ffe9f))
* limit successful response bodies ([f709b9b](https://github.com/nemirlev/zenmoney-go-sdk/commit/f709b9bcf4b15973c7da5b12183d72276f51a44e))


### 🐛 Bug Fixes

* validate and resolve API base URLs ([f75253b](https://github.com/nemirlev/zenmoney-go-sdk/commit/f75253b4183d9da8d5ce992934012c506b600748))


### 🧪 Testing

* audit models against tracked diff fixture ([3703fb0](https://github.com/nemirlev/zenmoney-go-sdk/commit/3703fb04497ac045701cc5c4c0e00d60b4c850d1))
* cover public API delegation ([a8824c2](https://github.com/nemirlev/zenmoney-go-sdk/commit/a8824c2e1e751f9b0419269ef8cb778d30e7e25e))


### 📚 Documentation

* improve public API documentation ([82558ca](https://github.com/nemirlev/zenmoney-go-sdk/commit/82558ca8e1fa1233cd104105b3d2807fa5c3a317))
* mark model alignment as breaking ([5f23be3](https://github.com/nemirlev/zenmoney-go-sdk/commit/5f23be38be9876ec1565aef276e19a173716c3f4))


### 👷 Continuous Integration

* add Codecov components ([694eb7b](https://github.com/nemirlev/zenmoney-go-sdk/commit/694eb7ba530c3843cd96fd42ab9bbd71d4f7b2f9))

## [2.0.7](https://github.com/nemirlev/zenmoney-go-sdk/compare/v2.0.6...v2.0.7) (2026-08-21)


### ⚠️ Breaking Changes

The model alignment in `99b2a5e` contains source-incompatible public API
changes despite being released as a patch version:

* Unix timestamps changed from `int` to `int64`, including sync cursors,
  deletion stamps, entity `Changed` fields, `Transaction.Created`, and
  subscription timestamps.
* `Company.FullTitle` and `Reminder.Comment` changed from `string` to
  `*string` so JSON `null` remains distinguishable from an empty string.
* `User.SubscriptionRenewalDate` changed from `*int` to `*int64`.

Consumers assigning timestamps to `int` variables must convert explicitly or
migrate those variables to `int64`. Consumers reading nullable strings must
check for `nil` before dereferencing them.


### 🐛 Bug Fixes

* align models with current diff schema ([99b2a5e](https://github.com/nemirlev/zenmoney-go-sdk/commit/99b2a5e998afe32433dbb8ef66571ff9305dd3fe))
* **ci:** restore stable Go 1.27 toolchain ([eafc6a8](https://github.com/nemirlev/zenmoney-go-sdk/commit/eafc6a8daf3a0d16e0d9c1a157746f554806cc59))
* expose SDK errors publicly ([c27be4d](https://github.com/nemirlev/zenmoney-go-sdk/commit/c27be4d0dcfecc062728644c57adc15ae18fc7cb))
* make HTTP retries safe and honor timeouts ([93f619e](https://github.com/nemirlev/zenmoney-go-sdk/commit/93f619e931f852c1426b8d1a1ef9117de9fa8401))
* preserve sync cursor when forcing entities ([19b6712](https://github.com/nemirlev/zenmoney-go-sdk/commit/19b6712d1b9f3eb332083a7ca5a4e2f03d2f4565))


### 🚜 Refactor

* keep force sync compatibility at API boundary ([d67863b](https://github.com/nemirlev/zenmoney-go-sdk/commit/d67863b0dff5cb30f61397b7f36ecbe2b7b235e4))


### 📚 Documentation

* update examples for corrected client behavior ([faa68fa](https://github.com/nemirlev/zenmoney-go-sdk/commit/faa68fa0c9ee2d81e69190d0359f923f5c682ea3))


### 👷 Continuous Integration

* test Go 1.26 and 1.27 release candidate ([ef06512](https://github.com/nemirlev/zenmoney-go-sdk/commit/ef06512cd296938c29fdcca63db261c7e6f734ec))

## [2.0.6](https://github.com/nemirlev/zenmoney-go-sdk/compare/v2.0.5...v2.0.6) (2026-08-21)


### 🐛 Bug Fixes

* **ci:** keep CodeQL compatible with Go toolchain ([4df1563](https://github.com/nemirlev/zenmoney-go-sdk/commit/4df1563e4d1b2728fe647b816ae090a58fc4e739))
* **ci:** migrate releases to release-please ([46ca8c0](https://github.com/nemirlev/zenmoney-go-sdk/commit/46ca8c03848ab7b831270b713e0b1a1c10ef197d))


### 👷 Continuous Integration

* always use latest Go patch release ([648fc54](https://github.com/nemirlev/zenmoney-go-sdk/commit/648fc54c66159b6090a637c8b190036615a073f1))


### ⚙️ Miscellaneous Tasks

* update Go dependencies ([93e1110](https://github.com/nemirlev/zenmoney-go-sdk/commit/93e1110ce397ce2b3b65789a273f0955e6646a6b))
* update Go to 1.27 ([bf1fdb3](https://github.com/nemirlev/zenmoney-go-sdk/commit/bf1fdb3cccfeef2c6b046e6b2867e5995b3fa590))

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
