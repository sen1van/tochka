# README.MD

- Бэкэнд: /backend
- Описание API: /api

api описывается с помощью typespec и компелируется в openapi

backend написан на go преимущественно без внешних импортов (sqlite из-за необходимости)



## Makefile:
- make (default) - компиляция всего проекта
- make run - компеляция и запуск
- make backend - компеляция бэкэнд (из Makefile в /backend)
