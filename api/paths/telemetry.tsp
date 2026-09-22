import "@typespec/http";

import "../models/errors.tsp";
import "../models/telemetry.tsp";
import "../models/utils.tsp";

namespace tochka;

using Http;

@route("/telemetry")
@tag("Telemetry")
interface Telemetry {
  @doc("create a new telemetry stamp")
  @post newTelemetryStamp(
    @body data: TelemetryData
  ): {@statusCode statusCode: 201} | ErrorBadRequest;

  @doc("get telemetry from all sensors")
  @get getAllTelemetry(
    @query query: PaginationQuery
  ): TelemetryResponse;
}
