import "@typespec/http";

import "../models/errors.tsp";
import "../models/telemetry.tsp";
import "../models/sensors.tsp";
import "../models/utils.tsp";

namespace tochka;

using Http;

@tag("Sensors")
@route("/sensors")
interface Sensors {
  @doc("create a new sensor")
  @post newSensor(
    @body data: CreateSensorRequest
  ): {@statusCode statusCode: 201} | ErrorBadRequest | ErrorConflict;

  @doc("get all sensors")
  @get getAllSensors(
  ): GetAllSensorResponse;

  @doc("delete sensor")
  @route("/{sensorId}")
  @delete deleteSensor(
    @path sensorId: int64
  ): {@statusCode statusCode: 204} | ErrorBadRequest | ErrorNotFound;

  @doc("update sensor")
  @route("/{sensorId}")
  @patch updateSensor(
    @path sensorId: int64,
    @body data: UpdateSensorRequest
  ): {@statusCode statusCode: 204} | ErrorBadRequest | ErrorNotFound;

  @doc("get sensor telemetry")
  @route("/{sensorId}/telemetry")
  @get getSensorTelemetry(
    @path sensorId: int64,
    @query query: PaginationQuery
  ): TelemetryResponse;
}
