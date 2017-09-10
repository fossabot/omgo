package com.omgo.webservice;

import com.omgo.webservice.etcd.Services;
import com.omgo.webservice.handler.LoginHandler;
import com.omgo.webservice.handler.RegisterHandler;
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
import io.vertx.ext.auth.User;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.Session;
import io.vertx.ext.web.handler.BodyHandler;
import io.vertx.ext.web.handler.CookieHandler;
import io.vertx.ext.web.handler.FormLoginHandler;
import io.vertx.ext.web.handler.SessionHandler;
import io.vertx.ext.web.sstore.LocalSessionStore;
import proto.DBServiceGrpc;
import proto.Db;
import proto.common.Common;

import java.util.ArrayList;
import java.util.List;

public class MainVerticle extends AbstractVerticle {
    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    private ManagedChannel grpcChannel;

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

        // login
        LoginHandler loginHandler = new LoginHandler(vertx, grpcChannel);
        loginHandler.register(router, ApiConstant.getApiPath(ApiConstant.API_LOGIN));

        // register
        RegisterHandler registerHandler = new RegisterHandler(vertx);
        registerHandler.register(router, ApiConstant.getApiPath(ApiConstant.API_REGISTER));

        // start service
        HttpServer server = vertx.createHttpServer();
        server.requestHandler(router::accept).listen(8080);
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

        grpcChannel = servicePool.getChannel(servicePool.getServicePath("dbservice"));
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
