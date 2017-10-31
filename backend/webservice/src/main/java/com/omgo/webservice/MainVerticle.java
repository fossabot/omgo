package com.omgo.webservice;

import com.codahale.metrics.MetricRegistry;
import com.codahale.metrics.SharedMetricRegistries;
import com.omgo.webservice.handler.*;
import com.omgo.webservice.service.Services;
import io.prometheus.client.CollectorRegistry;
import io.prometheus.client.dropwizard.DropwizardExports;
import io.prometheus.client.hotspot.DefaultExports;
import io.prometheus.client.vertx.MetricsHandler;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServer;
import io.vertx.core.json.JsonArray;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.handler.BodyHandler;
import io.vertx.ext.web.handler.CookieHandler;
import io.vertx.ext.web.handler.SessionHandler;
import io.vertx.ext.web.sstore.LocalSessionStore;
import io.vertx.ext.web.sstore.SessionStore;

import java.util.ArrayList;
import java.util.List;

public class MainVerticle extends AbstractVerticle {
    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);
    private Services.Pool dataCenters;

    public static void main(String[] args) {
        Vertx vertx = Vertx.vertx();
        vertx.deployVerticle(new MainVerticle());
    }

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString("info.version"));

        Utils.DEBUG = config().getBoolean("debug", false);
        Utils.STANDALONE = config().getBoolean("standalone", true);

        setupServices();
        startApiService();
    }

    private void startApiService() {
        Router router = Router.router(vertx);

        // Cookies, sessions and request bodies
        router.route().handler(CookieHandler.create());
        router.route().handler(BodyHandler.create());
        SessionStore store = LocalSessionStore.create(
            vertx,
            config().getString("session.map"),
            config().getLong("session.expire", 24 * 60 * 60 * 1000L));
        router.route().handler(SessionHandler.create(store));

        if (!Utils.STANDALONE) {
            // login
            LoginHandler loginHandler = new LoginHandler(vertx, dataCenters);
            loginHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_LOGIN));

            // register
            RegisterHandler registerHandler = new RegisterHandler(vertx, dataCenters);
            registerHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_REGISTER));
        }

        // handshake
        HandshakeHandler handshakeHandler = new HandshakeHandler(vertx);
        handshakeHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_HANDSHAKE));

        // user profile
        UserProfileHandler userProfileHandler = new UserProfileHandler(vertx, dataCenters);
        userProfileHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_USERPROFILE));

        // test
        TestHandler testHandler = new TestHandler(vertx);
        testHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_TEST));

        // metrics
        DefaultExports.initialize();
        new DropwizardExports(SharedMetricRegistries.getOrCreate("vertx.http.servers.localhost:8080")).register();
        router.get("/metrics").handler(new MetricsHandler());

        // start service
        HttpServer server = vertx.createHttpServer();
        server.requestHandler(router::accept).listen(8080);
    }

    /**
     * setup service pool
     */
    private void setupServices() {
        if (Utils.STANDALONE) {
            return;
        }

        // init etcd
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

        List<String> serviceTypes = new ArrayList<>();
        JsonArray typesArray = config().getJsonArray("service.types", new JsonArray().add("dataservice"));
        for (int i = 0; i < typesArray.size(); i++) {
            String name = typesArray.getString(i);
            serviceTypes.add(name);
        }

        LOGGER.info("service names:" + serviceTypes);

        Services.getInstance().init(endpoints);

        LOGGER.info("Services inited");

        // init dataservice

        dataCenters = Services.getInstance().getServicePool(vertx, root, "dataservice");
        dataCenters.addOnChangeListener(new Services.Pool.OnChangeListener() {
            @Override
            public void onServiceAdded(Services.Pool pool) {
                LOGGER.info("new service added");
            }

            @Override
            public void onServiceRemoved(Services.Pool pool) {
                LOGGER.info("service removed");
            }
        });

        // init agent manager
        AgentManager.getInstance().init(vertx, root, "agent");
    }
}
