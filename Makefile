.PHONY: default backend api

default: backend api

backend:
	cd backend && make

api:
	cd api && make


run: default
	cd builds && ./backend

lint:
	cd backend && golangci-lint run
