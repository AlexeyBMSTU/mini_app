.PHONY: setup run

setup:
	# Копируем .env.example в .env (если его нет)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env from .env.example"; \
	else \
		echo ".env already exists, skipping..."; \
	fi
	
	# Копируем .env в backend/.env
	@if [ -f .env ]; then \
		cp .env backend/.env; \
		echo "Copied .env to backend/.env"; \
	else \
		echo "Warning: .env not found, creating from example..."; \
		cp .env.example backend/.env; \
	fi
	
	# Копируем .env в frontend/.env
	@if [ -f .env ]; then \
		cp .env frontend/.env; \
		echo "Copied .env to frontend/.env"; \
	else \
		echo "Warning: .env not found, creating from example..."; \
		cp .env.example frontend/.env; \
	fi

run: setup
	@echo "Starting backend and frontend in parallel..."
	@cd backend && go run cmd/app/main.go &
	@sleep 2 && cd frontend && npm run build:dev && npm run dev

run-backend: setup
	@echo "Starting backend..."
	@cd backend && go run cmd/app/main.go

run-frontend: setup
	@echo "Building and starting frontend..."
	@cd frontend && npm run build:dev && npm run dev

build-frontend:
	@echo "Building frontend..."
	@cd frontend && npm run build:dev

install:
	@echo "Installing backend dependencies..."
	@cd backend && go mod tidy && go mod download
	
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install

clean:
	@echo "Cleaning up..."
	@rm -rf backend/.env backend/go.sum frontend/dist frontend/node_modules frontend/.env 
	@echo "Cleaned up (.env, dist, node_modules, go.sum) files from subdirectories"