import "@typespec/http";

using Http;

model SensorData {
  sensorId: int64,
  name: string,
  lastSeen: utcDateTime,
}

model CreateSensorRequest {
  sensorId: int64,
  name: string,
}

model UpdateSensorRequest {
  name: string,
}

model GetAllSensorResponse {
  total: int64,
  data: SensorData[],
}
