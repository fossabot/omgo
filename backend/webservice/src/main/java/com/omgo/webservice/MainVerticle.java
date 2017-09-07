package com.omgo.webservice;

import com.omgo.webservice.etcd.Services;
import io.grpc.ManagedChannel;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.http.HttpMethod;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.handler.BodyHandler;
import io.vertx.ext.web.handler.CookieHandler;
import io.vertx.ext.web.handler.SessionHandler;
import io.vertx.ext.web.sstore.LocalSessionStore;
import proto.DBServiceGrpc;
import proto.Db;
import proto.common.Common;

import java.util.ArrayList;
import java.util.List;

public class MainVerticle extends AbstractVerticle {
    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString("info.version"));

        setupServices();
        startApiService();
        testDBservice();
    }

    private void startApiService() {
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

    /**
     * setup service pool
     */
    private void setupServices() {
        List<String> endpoints = new ArrayList<>();
        JsonArray endpointsJA = config().getJsonArray("etcd.host", new JsonArray().add("http://localhost:2379"));
        for (int i = 0; i < endpointsJA.size(); i++) {
            String endpoint = endpointsJA.getString(i);
            endpoints.add(endpoint);
        }

        LOGGER.info("etcd host:" + endpoints);
        Services.getInstance().init(endpoints);

        String root = config().getString("service.root", "backends");

        LOGGER.info("service root:" + root);

        List<String> serviceNames = new ArrayList<>();
        JsonArray namesJA = config().getJsonArray("service.names", new JsonArray().add("dbservice"));
        for (int i = 0; i < namesJA.size(); i++) {
            String name = namesJA.getString(i);
            serviceNames.add(name);
        }

        LOGGER.info("service names:" + serviceNames);

        Services.ServicePool servicePool = Services.getInstance().createServicePool(vertx, root, serviceNames);
        LOGGER.info("service pool created");
    }

    private void testDBservice() {
        Services.ServicePool pool = Services.getInstance().getServicePool();
        ManagedChannel channel = pool.getChannel(pool.getServicePath("dbservice"));
        DBServiceGrpc.DBServiceVertxStub stub = DBServiceGrpc.newVertxStub(channel);

        Common.UserInfo.Builder userInfoBuilder = Common.UserInfo.newBuilder();
        userInfoBuilder.setEmail("masterg@yeah.net");
        Db.DB.UserExtendInfo.Builder extendInfoBuilder = Db.DB.UserExtendInfo.newBuilder();
        extendInfoBuilder.setInfo(userInfoBuilder.build())
            .setSecret("g3st4p0");

        stub.userLogin(extendInfoBuilder.build(), res -> {
            if (res.succeeded()) {
                LOGGER.info(res.result());
            } else {
                LOGGER.warn(res.cause());
            }
        });
    }
}
