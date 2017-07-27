package com.omgo.dbservice;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.asyncsql.MySQLClient;
import io.vertx.ext.sql.SQLClient;
import io.vertx.grpc.VertxServer;
import io.vertx.grpc.VertxServerBuilder;
import io.vertx.redis.RedisClient;
import io.vertx.redis.RedisOptions;

import java.io.IOException;

public class MainVerticle extends AbstractVerticle {

    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    @Override
    public void start() {
        String rpcHost = config().getString("rpc.host", "localhost");
        int rpcPort = config().getInteger("rpc.port", 60001);

        VertxServer rpcServer = VertxServerBuilder
            .forAddress(vertx, rpcHost, rpcPort)
            .addService(new DBServiceGrpcImpl(createSQLClient(), createRedisClient()))
            .build();

        // Start is asynchronous
        try {
            rpcServer.start();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private SQLClient createSQLClient() {
        String host = config().getString("sql.host", "localhost");
        int port = config().getInteger("sql.port", 3306);
        int maxPoolSize = config().getInteger("sql.maxPoolSize", 10);
        String username = config().getString("sql.username", "driver");
        String password = config().getString("sql.password", "mysql");
        String database = config().getString("sql.database", "master");
        String charset = config().getString("sql.charset", "UTF-8");
        int queryTimeout = config().getInteger("sql.queryTimeout", 10000);

        JsonObject mySQLConnectionConfig = new JsonObject()
            .put("host", host)
            .put("port", port)
            .put("maxPoolSize", maxPoolSize)
            .put("username", username)
            .put("password", password)
            .put("database", database)
            .put("charset", charset)
            .put("queryTimeout", queryTimeout);

        LOGGER.info(mySQLConnectionConfig);

        return MySQLClient.createNonShared(vertx, mySQLConnectionConfig);
    }

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

}
