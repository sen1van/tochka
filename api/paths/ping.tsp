import "@typespec/http";

import "../models/ping.tsp";

namespace tochka;

using Http;

@route("/ping")
@tag("Ping")
interface Ping {
  @doc("Check authorization")
  @get ping(): PingResponse;
}
