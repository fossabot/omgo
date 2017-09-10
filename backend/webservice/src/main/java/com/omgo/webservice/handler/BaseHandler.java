package com.omgo.webservice.handler;

import io.vertx.core.Vertx;
import io.vertx.core.http.HttpMethod;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;

public class BaseHandler {
    protected String MIME_JSON = "application/json";
    protected String CONTENT_TYPE = "content-type";

    protected Logger LOGGER;
    protected Vertx vertx;

    protected Route route;

    public BaseHandler(Vertx vertx) {
        this.vertx = vertx;
        LOGGER = LoggerFactory.getLogger(this.getClass());
    }

    public void register(Router router, String path) {
        route = router.route(httpMethod(), path)
            .consumes(consumes())
            .produces(produces());
    }

    protected HttpMethod httpMethod() {
        return HttpMethod.GET;
    }

    protected String consumes() {
        return MIME_JSON;
    }

    protected String produces() {
        return MIME_JSON;
    }
}
