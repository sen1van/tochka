# API

int: 12345
float: 1.2345 (not used)
string: "Ivan Ivanov"
time: "2025-01-01T00:00:00Z"


GET /api/v0/ping
ANSW:
  204 No Content
  
POST /api/v0/telemetry
BODY:
  json: {
    "sensor_id": 12345,
    "try_number": 12345,
    "value": 12345,
    "timestamp": "2025-01-01T00:00:00Z"
  }
ANSW:
  201 Created

GET /api/v0/telemetry?limit={last N results by timestamp, default 20}&offset={offset, default 0}
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


GET /api/v0/sensors/{sensor_id}/telemetry?limit={last N results by timestamp, default 20}&offset={offset, default 0}
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


POST /api/v0/sensors
BODY:
  json: {
    "sensor_id": 12345,
    "name": "ivan ivanov"
  }
ANSW:
  201 Created


PATCH /api/v0/sensors/{sensor_id}
BODY:
  json: {
    "name": "ivan ivanov"
  }
ANSW:
  204 No Content


DELETE /api/v0/sensors/{sensor_id}
ANSW:
  204 No Content


GET /api/v0/sensors
ANSW:
  json: [
    {
      "sensor_id": 12345,
      "name": "Ivan Ivanov",
      "last_seen": "2025-01-01T00:00:00Z"
    }
  ]
