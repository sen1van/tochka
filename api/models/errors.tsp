import "@typespec/http";

using Http;

@error
model ErrorBadRequest {
  message: string;
}
