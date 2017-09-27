package com.omgo.webservice.handler;

import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;

public class TestHandler extends BaseHandler {
    public TestHandler(Vertx vertx) {
        super(vertx);
    }

    @Override
    public void initRoute(Router router, String path) {
        super.initRoute(router, path);

        route.handler(routingContext -> {
            HttpServerRequest request = super.getRequest(routingContext);
            HttpServerResponse response = super.getResponse(routingContext);

            if (!isSessionValid(routingContext)) {
                routingContext.fail(401);
                return;
            }

            JsonObject rsp = new JsonObject();
            rsp.put("foo", "bar");
            JsonObject headerJson = getHeaderJson(request);
            response.write(rsp.encode()).end();
        });
    }

}
