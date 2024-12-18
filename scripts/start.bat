@echo off

:prompt_env_files
set /p choice="Do you want to use the default .env files? (y/n): "
if /i "%choice%"=="y" (
    echo Using default .env files...
    copy /Y configs\Docker.env.template configs\Docker.env
    copy /Y frontend\.env.template frontend\.env
) else if /i "%choice%"=="n" (
    echo Please create your .env files in 'configs\Docker.env' and 'frontend\.env', you can use configs\Docker.env.template and frontend\.env.template as templates. Do not forget to copy configs\Docker.env to configs\.env.
    exit /b 0
) else (
    echo Invalid choice. Please enter y or n.
    goto prompt_env_files
)

copy /Y configs\Docker.env configs\.env

echo .env files are ready, starting the application using docker...

docker compose -f deployments\docker-compose.yml up --build -V