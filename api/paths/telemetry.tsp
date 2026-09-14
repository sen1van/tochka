import "@typespec/http";

import "../models/errors.tsp";
import "../models/telemetry.tsp";

namespace tochka;

using Http;

@route("/telemetry")
@tag("Telemetry")
interface Telemetry {
  @doc("create a new telemetry stamp")
  @post newTelemetryStamp(
    @body data: TelemetryData[]
  ): {@statusCode statusCode: 201} | ErrorBadRequest;

  @doc("get telemetry from all sensors")
  @get getAllTelemetry(
    @query limit: int64,
    @query offset: int64
  ): TelemetryResponse;
}
