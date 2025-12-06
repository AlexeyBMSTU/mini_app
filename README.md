# Telegram Mini App

## Основные команды
- `make setup` Установка переменных окружения .env
- `make install` Установка всех зависимостей
- `make clean` Очистка зависимостей (.env, dist, node_modules, go.sum)
- `make run` Запуск приложения
- `make run-backend` Запуск бота и http сервера
- `make run-frontend` Запуск клиентской части
- `make build-frontend` Билд клиентской части

## Для полного запуска необходимо выполнить:
1. `make setup`
2. `make install`
3. `make run`

## Для запуска бота и http сервера необходимо выполнить:
1. `make setup`
2. `make install`
3. `make run-backend`

## Для запуска клиентской части необходимо выполнить:
1. `make setup`
2. `make install`
3. `make run-frontend`