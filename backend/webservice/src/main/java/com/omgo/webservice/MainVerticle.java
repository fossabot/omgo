package com.omgo.webservice;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.http.HttpMethod;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.handler.BodyHandler;
import io.vertx.ext.web.handler.CookieHandler;
import io.vertx.ext.web.handler.SessionHandler;
import io.vertx.ext.web.sstore.LocalSessionStore;

public class MainVerticle extends AbstractVerticle {
    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString("info.version"));

        Router router = Router.router(vertx);

        // Cookies, sessions and request bodies
        router.route().handler(CookieHandler.create());
        router.route().handler(BodyHandler.create());
        router.route().handler(SessionHandler.create(LocalSessionStore.create(vertx)));

        registerAuthRoute(router);

        HttpServer server = vertx.createHttpServer();
        server.requestHandler(router::accept).listen(8080);
    }

    private void registerAuthRoute(Router router) {
        // Simple auth service with uses a properties file for user/role info

        Route route = router.route(HttpMethod.GET, ApiConstant.getApiPath(ApiConstant.API_LOGIN))
            .consumes(Constants.MIME_JSON)
            .produces(Constants.MIME_JSON);

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
