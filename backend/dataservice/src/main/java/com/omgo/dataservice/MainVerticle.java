package com.omgo.dataservice;

import com.omgo.dataservice.etcd.Services;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.json.JsonArray;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
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

    @Override
    public void start() {
        LOGGER.info("config version: " + config().getString("info.version"));
        LOGGER.info("config debug: " + config().getBoolean("debug", false));

        serviceHost = config().getString("service.host", "localhost");
        servicePort = config().getInteger("service.port", 60001);

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

    /**
     * create a mongodb client
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

        String root = config().getString("service.root", "backends");
        String selfKind = config().getString("service.kind", "dataservice");
        String selfName = config().getString("service.self", "ds-0");

        LOGGER.info("service root:" + root);

        // register self to etcd as service
        Services.getInstance().registerService(Services.generatePath(root, selfKind, selfName), String.format("%s:%d", serviceHost, servicePort));
        LOGGER.info("service registered");
    }
}
