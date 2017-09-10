package com.omgo.webservice.handler;

import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;

public class RegisterHandler extends BaseHandler {
    public RegisterHandler(Vertx vertx) {
        super(vertx);
    }

    @Override
    public void register(Router router, String path) {
        super.register(router, path);

        route.handler(routingContext -> {
            HttpServerRequest request = routingContext.request();

            LOGGER.info("handling request: " + request.uri());
            JsonObject jsonObject = new JsonObject();
            jsonObject.put("uid", 1000);

            String email = request.headers().get("email");
            String password = request.headers().get("password");

            LOGGER.info("email:" + email);
            LOGGER.info("password:" + password);

            HttpServerResponse response = routingContext.response();
            // enable chunked responses because we will be adding data as
            // we execute over other handlers. This is only required once and
            // only if several handlers do output.
            response.setChunked(true);

            response.putHeader("content-type", "application/json");
            response.write(jsonObject.encode()).end();
        });
    }
}
