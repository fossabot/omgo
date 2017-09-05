package com.omgo.webservice;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;

public class MainVerticle extends AbstractVerticle {
    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString("info.version"));

        HttpServer server = vertx.createHttpServer();

        Router router = Router.router(vertx);

        Route route1 = router.route("/some/path").handler(routingContext -> {
            HttpServerResponse response = routingContext.response();
            response.setChunked(true);
            response.write("route1\n");
            routingContext.vertx().setTimer(5000, tid -> routingContext.next());
        });

        Route route2 = router.route("/some/path").handler(routingContext -> {
            HttpServerResponse response = routingContext.response();
            response.write("route2\n");
            routingContext.vertx().setTimer(5000, tid -> routingContext.next());
        });


        Route route3 = router.route("/some/path").handler(routingContext -> {
            HttpServerResponse response = routingContext.response();
            response.write("route3");
            routingContext.response().end();
        });

        server.requestHandler(router::accept).listen(8080);
    }
}
