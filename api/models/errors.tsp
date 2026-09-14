import "@typespec/http";

using Http;

@error
model ErrorBadRequest {
  @statusCode _: 400;
  message: string;
}

@error
model ErrorUnauthorized {
  @statusCode _: 401;
  message: string;
}

@error
model ErrorNotFound {
  @statusCode _: 404;
  message: string;
}

@error
model ErrorConflict {
  @statusCode _: 409;
  message: string;
}
