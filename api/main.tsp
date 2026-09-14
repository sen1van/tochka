import "@typespec/http";

import "./paths/ping.tsp";
import "./paths/telemetry.tsp";
import "./paths/sensors.tsp";


using Http;

@service(#{ title: "tochka" })
@server("https://example.try", "Server")
@route("/api")
@useAuth(ApiKeyAuth<ApiKeyLocation.header, "X-API-Token">)
namespace tochka;
