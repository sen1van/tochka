.PHONY: default backend api frontend

default: backend api frontend

backend:
	cd backend && make

api:
	cd api && make

frontend:
	cd frontend && make

run: default
	cd builds && ./backend

lint:
	cd backend && golangci-lint run
