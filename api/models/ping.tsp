import "@typespec/http";

using Http;

model PingResponse {
  @statusCode statusCode: 200;
  message: "Pong";
}
