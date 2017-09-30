package com.omgo.webservice.handler;

import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;

public class TestHandler extends BaseHandler {
    public TestHandler(Vertx vertx) {
        super(vertx);
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response) {
        HttpServerRequest request = super.getRequest(routingContext);

        JsonObject rsp = getResponseJson();
        rsp.put("foo", "bar");
        JsonObject headerJson = getHeaderJson(request);
        response.write(rsp.encode()).end();
    }
}
