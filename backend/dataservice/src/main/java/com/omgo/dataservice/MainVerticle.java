package com.omgo.dataservice;

import com.omgo.utils.Services;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.DeploymentOptions;
import io.vertx.core.Vertx;
import io.vertx.core.VertxOptions;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.dropwizard.DropwizardMetricsOptions;
import io.vertx.ext.dropwizard.Match;
import io.vertx.ext.dropwizard.MatchType;
import io.vertx.ext.mongo.MongoClient;
import io.vertx.grpc.VertxServer;
import io.vertx.grpc.VertxServerBuilder;
import io.vertx.redis.RedisClient;
import io.vertx.redis.RedisOptions;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class MainVerticle extends AbstractVerticle {

    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    private String serviceHost;
    private int servicePort;

    private String root;
    private String selfKind;
    private String selfName;

    private String selfFullPath;

    public static void main(String[] args) {

        String cfgPath = "";
        for (int i = 0; i < args.length - 1; i++) {
            if (args[i].equals("-conf")) {
                cfgPath = args[i + 1];
                break;
            }
        }

        Vertx vertx = Vertx.vertx(new VertxOptions().setMetricsOptions(
            new DropwizardMetricsOptions()
                .setEnabled(true)
                .addMonitoredHttpServerUri(
                    new Match().setValue("/"))
                .addMonitoredHttpServerUri(
                    new Match().setValue("/api/*").setType(MatchType.REGEX))
        ));

        JsonObject configObject;
        if (!cfgPath.isEmpty()) {
            Buffer fileBuf = vertx.fileSystem().readFileBlocking(cfgPath);
            configObject = new JsonObject(fileBuf);
        } else {
            configObject = new JsonObject();
        }

        vertx.deployVerticle(new MainVerticle(), new DeploymentOptions().setConfig(configObject));
    }

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString("info.version"));
        LOGGER.info("config debug: " + config().getBoolean("debug", false));

        serviceHost = config().getString("service.host", "localhost");
        servicePort = config().getInteger("service.port", 60001);
        root = config().getString("service.root", "backends");
        selfKind = config().getString("service.kind", "dataservice");
        selfName = config().getString("service.self", "ds-0");

        selfFullPath = Services.generatePath(root, selfKind, selfName);

        LOGGER.info("service full path:" + selfFullPath);

        VertxServer rpcServer = VertxServerBuilder
            .forPort(vertx, servicePort)
            .addService(new DBServiceGrpcImpl(createMongoClient(), createRedisClient()))
            .build();

        // Start is asynchronous
        try {
            rpcServer.start();
        } catch (IOException e) {
            e.printStackTrace();
        }

        setupServices();
    }

    @Override
    public void stop() {
        try {
            super.stop();
        } catch (Exception e) {
            e.printStackTrace();
        }

        Services.getInstance().unregisterService(selfFullPath);
    }

    /**
     * create a mongodb client
     *
     * @return
     */
    private MongoClient createMongoClient() {
        return MongoClient.createShared(vertx, config().getJsonObject("mongodb.config"));
    }

    /**
     * create and setup redis client
     *
     * @return
     */
    private RedisClient createRedisClient() {
        String host = config().getString("redis.host", "localhost");
        int port = config().getInteger("redis.port", 6379);
        String encoding = config().getString("redis.encoding", "UTF-8");
        boolean tcpKeepAlive = config().getBoolean("redis.tcpKeepAlive", true);
        boolean tcpNoDelay = config().getBoolean("redis.tcpNoDelay", true);

        RedisOptions redisConfig = new RedisOptions()
            .setHost(host)
            .setPort(port)
            .setEncoding(encoding);

        redisConfig
            .setTcpKeepAlive(tcpKeepAlive)
            .setTcpNoDelay(tcpNoDelay);

        LOGGER.info(redisConfig);

        return RedisClient.create(vertx, redisConfig);
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

        // register self to service as service
        Services.getInstance().registerService(selfFullPath, String.format("%s:%d", serviceHost, servicePort));
        LOGGER.info("service registered");
    }
}
