import "@typespec/http";

import "../models/errors.tsp";
import "../models/telemetry.tsp";
import "../models/sensors.tsp";

namespace tochka;

using Http;

@tag("Sensors")
@route("/sensors")
interface Sensors {
  @doc("create a new sensor")
  @post newSensor(
    @body data: CreateSensorRequest
  ): {@statusCode statusCode: 201} | ErrorBadRequest;

  @doc("get all sensors")
  @get getAllSensors(
  ): GetAllSensorResponse;

  @doc("delete sensor")
  @route("/{sensorId}")
  @delete deleteSensor(
    @path sensorId: int64
  ): {@statusCode statusCode: 204} | ErrorBadRequest;

  @doc("update sensor")
  @route("/{sensorId}")
  @put updateSensor(
    @path sensorId: int64,
    @body data: UpdateSensorRequest
  ): {@statusCode statusCode: 204} | ErrorBadRequest;

  @doc("get sensor telemetry")
  @route("/{sensorId}/telemetry")
  @get getSensorTelemetry(
    @path sensorId: int64,
    @query limit: int64,
    @query offset: int64
  ): TelemetryResponse;
}
