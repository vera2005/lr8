# build.ps1

# Запуск golangci-lint
Write-Host "Запуск golangci-lint..."
golangci-lint run ./...

# Проверка статуса выполнения линтера
if ($LASTEXITCODE -ne 0) {
    Write-Host "Ошибка при выполнении golangci-lint."
    exit $LASTEXITCODE
}

# Компиляция проекта
Write-Host "Компиляция проекта..."
go build -o myapp.exe .

# Проверка статуса компиляции
if ($LASTEXITCODE -ne 0) {
    Write-Host "Ошибка при компиляции проекта."
    exit $LASTEXITCODE
}

Write-Host "Скрипт завершен."
