package com.omgo.webservice;

import com.codahale.metrics.MetricRegistry;
import com.codahale.metrics.SharedMetricRegistries;
import com.omgo.utils.ConfigUtils;
import com.omgo.utils.Services;
import com.omgo.utils.Utils;
import com.omgo.webservice.handler.*;
import io.prometheus.client.CollectorRegistry;
import io.prometheus.client.dropwizard.DropwizardExports;
import io.prometheus.client.exporter.common.TextFormat;
import io.prometheus.client.hotspot.DefaultExports;
import io.prometheus.client.vertx.MetricsHandler;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.DeploymentOptions;
import io.vertx.core.Vertx;
import io.vertx.core.VertxOptions;
import io.vertx.core.http.HttpServer;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.dropwizard.DropwizardMetricsOptions;
import io.vertx.ext.dropwizard.Match;
import io.vertx.ext.dropwizard.MatchType;
import io.vertx.ext.dropwizard.MetricsService;
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

        Vertx vertx = Vertx.vertx(new VertxOptions().setMetricsOptions(
            new DropwizardMetricsOptions()
                .setEnabled(true)
//                .addMonitoredHttpServerUri(
//                    new Match().setValue("/"))
                .addMonitoredHttpServerUri(
                    new Match().setValue("/api/*").setType(MatchType.REGEX))
        ));

        String cfgPath = ConfigUtils.extractConfigPath(args);
        JsonObject configObject = ConfigUtils.loadConfigFromPath(vertx, cfgPath, new JsonObject());
        vertx.deployVerticle(new MainVerticle(), new DeploymentOptions().setConfig(configObject));
    }

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString(ConfigUtils.INFO_VERSION));

        Utils.DEBUG = config().getBoolean(ConfigUtils.DEBUG, false);
        Utils.STANDALONE = config().getBoolean(ConfigUtils.STANDALONE, true);

        setupServices();
        startApiService();
    }

    private void startApiService() {
        Router router = Router.router(vertx);

        // create http server
        HttpServer server = vertx.createHttpServer();

        // Cookies, sessions and request bodies
        router.route().handler(CookieHandler.create());
        router.route().handler(BodyHandler.create());
        SessionStore store = LocalSessionStore.create(
            vertx,
            config().getString(ConfigUtils.SESSION_MAP),
            config().getLong(ConfigUtils.SESSION_EXPIRE, 24 * 60 * 60 * 1000L));
        router.route().handler(SessionHandler.create(store));

        // login
        LoginHandler loginHandler = new LoginHandler(vertx, dataCenters);
        loginHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_LOGIN));

        // register
        RegisterHandler registerHandler = new RegisterHandler(vertx, dataCenters);
        registerHandler.setRoute(router, ApiConstant.getApiPath(ApiConstant.API_REGISTER));

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
        MetricRegistry metricRegistry = SharedMetricRegistries.getOrCreate("webservice.metrics");
        CollectorRegistry.defaultRegistry.register(new DropwizardExports(metricRegistry));
        router.get(config().getString(ConfigUtils.METRICS_PATH, "/metrics")).handler(new MetricsHandler());

        // this is just an example about how to get dropwizard metrics
        MetricsService metricsService = MetricsService.create(vertx);
        router.get("/go").handler(res -> {
            // set up server
            JsonObject metrics = metricsService.getMetricsSnapshot(server);


            res.response()
                .setStatusCode(200)
                .putHeader("Content-Type", TextFormat.CONTENT_TYPE_004)
                .end(metrics.encode());
        });

        // start http server
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
        JsonArray endpointsJA = config().getJsonArray(ConfigUtils.ETCD_HOST, new JsonArray().add("http://localhost:2379"));
        for (int i = 0; i < endpointsJA.size(); i++) {
            String endpoint = endpointsJA.getString(i);
            endpoints.add(endpoint);
        }

        LOGGER.info("etcd host:" + endpoints);
        Services.getInstance().init(endpoints);

        String root = config().getString(ConfigUtils.SERVICE_ROOT, "backends");

        LOGGER.info("service root:" + root);

        List<String> serviceTypes = new ArrayList<>();
        JsonArray typesArray = config().getJsonArray(ConfigUtils.SERVICE_TYPES, new JsonArray().add("dataservice"));
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
