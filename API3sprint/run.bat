@echo off
setlocal

:: Устанавливаем переменные окружения
set CONFIG_PATH=%~dp0config\local.yaml
set GIN_MODE=release

:: Переходим в директорию проекта
cd /d %~dp0

:: Проверяем существование конфига
if not exist "%CONFIG_PATH%" (
    echo Error: Config file not found at %CONFIG_PATH%
    pause
    exit /b 1
)

:: Создаем папку storage если не существует
if not exist "storage" mkdir storage

:: Запускаем приложение
echo Starting server with config: %CONFIG_PATH%
go run cmd/url-shortener/main.go

endlocal