default:
	cd backend && make

run: default
	cd builds && ./backend
