@echo off
set CONFIG_PATH=%~dp0config\local.yaml
cd /d %~dp0
go run cmd/url-shortener/main.go