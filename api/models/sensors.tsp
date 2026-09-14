import "@typespec/http";

import "./telemetry.tsp";

using Http;

model SensorData {
  sensorId: int64,
  name: string,
  lastTelemetry: TelemetryData | null,
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
