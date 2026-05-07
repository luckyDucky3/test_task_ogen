# Ogen x-description task

Проект показывает доработку генератора `ogen`: у методов в сгенерированном Go-коде рядом со стандартным `description` добавляется комментарий из OpenAPI extension `x-description`.

## Что внутри

- `third_party/ogen` - локальная копия `github.com/ogen-go/ogen` с правкой генератора.
- `open-api.yaml` - OpenAPI-спецификация CRUD API для ресурса `Book`.
- `api` - код, сгенерированный измененным `ogen`.
- `main.go` - in-memory реализация сервера и клиентский CRUD-сценарий через `httptest`.

## Команды

```powershell
Push-Location .\third_party\ogen
& 'C:\Users\maxle\sdk\go1.26.2\bin\go.exe' run ./cmd/ogen --target ..\..\api --package api --clean ..\..\open-api.yaml
Pop-Location
& 'C:\Users\maxle\sdk\go1.26.2\bin\go.exe' test ./...
& 'C:\Users\maxle\sdk\go1.26.2\bin\go.exe' run .
```
