@echo off
go run main.go .\examples\main.dt 

if %ERRORLEVEL% NEQ 0 (
    echo Skipped running...
    exit /b %ERRORLEVEL%
)

cd examples 
java main
cd .. 