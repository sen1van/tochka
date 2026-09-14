import "@typespec/http";

using Http;

model TelemetryData {
  sensorId: int64;
  tryNumber: int64;
  value: int64;
  timestamp: utcDateTime;
}

model TelemetryResponse {
  @statusCode statusCode: 200;
  total: int64;
  data: TelemetryData[];
}
