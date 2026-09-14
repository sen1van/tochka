# README.MD

- Бэкэнд: /backend
- Описание API: /api
- Фронтэнд: /frontend (сделано ИИ)

api описывается с помощью typespec и компелируется в openapi

backend написан на go преимущественно без внешних импортов (sqlite из-за необходимости)

frontend написан на vue, при помощии ИИ по openapi.yaml

## Makefile:
- make (default) - компиляция всего проекта
- make run - компеляция и запуск
- make backend - компеляция бэкэнд (из Makefile в /backend)
- make frontend - компеляция фронтэнд (из Makefile в /frontend)
- make api - компеляция api (из Makefile в /api)
