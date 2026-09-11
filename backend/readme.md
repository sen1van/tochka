# API

типы данных обозначены:
- int: 12345
- float: 1.2345 (not used)
- string: "Ivan Ivanov"
- time: "2025-01-01T00:00:00Z" [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)

```
GET /api/ping
ANSW:
  200 {"message": "pong"}
```
  
 ``` 
POST /api/telemetry
BODY:
  json: {
    "sensor_id": 12345,
    "try_number": 12345,
    "value": 12345,
    "timestamp": "2025-01-01T00:00:00Z"
  }
ANSW:
  201 Created
```

```
GET /api/telemetry?limit={results by timestamp, default 20}&offset={offset, default 0}
ANSW:
  json: {
    "total": 12345, 
    "data": [
      {
        "sensor_id": 12345,
        "try_number": 12345,
        "value": 12345,
        "timestamp": "2025-01-01T00:00:00Z"
        },
    ...]
    }
```

```
GET /api/sensors/{sensor_id}/telemetry?limit={results by timestamp, default 20}&offset={offset, default 0}
ANSW:
  json: {
    "total": 12345, 
    "data": [
      {
        "sensor_id": 12345,
        "try_number": 12345,
        "value": 12345,
        "timestamp": "2025-01-01T00:00:00Z"
        },
    ...]
    }
```

```
POST /api/sensors
BODY:
  json: {
    "sensor_id": 12345,
    "name": "ivan ivanov"
  }
ANSW:
  201 Created
```

```
PUT /api/sensors/{sensor_id}
BODY:
  json: {
    "name": "ivan ivanov"
  }
ANSW:
  204 No Content
```

```
DELETE /api/sensors/{sensor_id}
ANSW:
  204 No Content
```

```
GET /api/sensors
ANSW:
  json: {
    "total": 12345,
    "data": [
      {
        "sensor_id": 12345,
        "name": "Ivan Ivanov",
        "last_seen": "2025-01-01T00:00:00Z"
      }
    ]
  }
```
