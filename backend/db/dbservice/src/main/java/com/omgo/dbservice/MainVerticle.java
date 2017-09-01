package com.omgo.dbservice;

import com.omgo.dbservice.etcd.Services;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Future;
import io.vertx.core.file.FileSystem;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.jdbc.JDBCClient;
import io.vertx.ext.sql.SQLClient;
import io.vertx.grpc.VertxServer;
import io.vertx.grpc.VertxServerBuilder;
import io.vertx.redis.RedisClient;
import io.vertx.redis.RedisOptions;
import io.vertx.redis.RedisTransaction;
import proto.Db;

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

        serviceHost = config().getString("service.host", "localhost");
        servicePort = config().getInteger("service.port", 60001);

        VertxServer rpcServer = VertxServerBuilder
            .forPort(vertx, servicePort)
            .addService(new DBServiceGrpcImpl(createSQLClient(), createRedisClient()))
            .build();

        // Start is asynchronous
        try {
            rpcServer.start();
        } catch (IOException e) {
            e.printStackTrace();
        }

        setupServices();

        testRedis();
    }

    /**
     * create and setup SQL client
     *
     * @return
     */
    private SQLClient createSQLClient() {
        String url = config().getString("sql.url", "jdbc:mysql://localhost:3306/master");
        int maxPoolSize = config().getInteger("sql.maxPoolSize", 10);
        String username = config().getString("sql.username", "driver");
        String password = config().getString("sql.password", "mysql");
        String host = config().getString("sql.host", "localhost");
        int port = config().getInteger("sql.port", 3306);
        String database = config().getString("sql.database", "master");
        String charset = config().getString("sql.charset", "UTF-8");

        JsonObject dataSourceProperty = new JsonObject()
            .put("databaseName", database)
            .put("portNumber", port)
            .put("serverName", host)
            .put("cachePrepStmts", true)
            .put("prepStmtCacheSize", 250)
            .put("prepStmtCacheSqlLimit", 2048)
            .put("useServerPrepStmts", true)
            .put("useLocalSessionState", true)
            .put("useLocalTransactionState", true)
            .put("rewriteBatchedStatements", true)
            .put("cacheResultSetMetadata", true)
            .put("cacheServerConfiguration", true)
            .put("elideSetAutoCommits", true)
            .put("maintainTimeStats", false);

        JsonObject mySQLConnectionConfig = new JsonObject()
            .put("provider_class", "io.vertx.ext.jdbc.spi.impl.HikariCPDataSourceProvider")
            .put("driverClassName", "com.mysql.cj.jdbc.Driver")
            .put("jdbcUrl", url)
            .put("maxPoolSize", maxPoolSize)
            .put("username", username)
            .put("password", password)
            .put("charset", charset)
            .put("datasource", dataSourceProperty);

        LOGGER.info(mySQLConnectionConfig);

        return JDBCClient.createShared(vertx, mySQLConnectionConfig);
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

        redisClient = RedisClient.create(vertx, redisConfig);
        return redisClient;
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
        String selfKind = config().getString("service.kind", "dbservice");
        String selfName = config().getString("service.self", "dbs-0");

        LOGGER.info("service root:" + root);

        List<String> serviceNames = new ArrayList<>();
        JsonArray namesJA = config().getJsonArray("service.names", new JsonArray().add("snowflake"));
        for (int i = 0; i < namesJA.size(); i++) {
            String name = namesJA.getString(i);
            serviceNames.add(name);
        }

        LOGGER.info("service names:" + serviceNames);

        Services.ServicePool servicePool = Services.getInstance().createServicePool(vertx, root, serviceNames);
        LOGGER.info("service pool created");

        // register self to etcd as service
        servicePool.registerService(Services.generatePath(root, selfKind, selfName), String.format("%s:%d", serviceHost, servicePort));
        LOGGER.info("service registered");
    }

    private static String TAG = "--------->  ";
    private RedisClient redisClient;

    private void testRedis() {
        RedisTransaction transaction = redisClient.transaction();
        transaction.multi(event -> {
            transaction.hgetall("user:100000", getAllEvent -> {
                if (getAllEvent.succeeded() && "QUEUED".equals(getAllEvent.result())) {
                    transaction.expire("user:100000", 10, expireEvent -> {
                        if (expireEvent.succeeded() && "QUEUED".equals(expireEvent.result())) {
                            transaction.exec(execEvent -> System.out.println(execEvent.result()));
                        } else {
                            transaction.discard(de -> {

                            });
                        }
                    });
                } else {
                    transaction.discard(discardEvent -> {
                    });
                }
            });
        });
    }

    private void testComposition() {
        Future<String> lastFuture = Future.future();
        lastFuture.setHandler(res -> {
           if (res.succeeded()) {
               LOGGER.info(TAG + "complete! " + res.result());
           } else {
               LOGGER.info(TAG + "complete! " + res.cause());
           }
        });

        Future<String> sqlFuture = getFailFuture("shit");
        sqlFuture.setHandler(res -> {
            if (res.succeeded()) {

            } else {
                Db.DB.StatusCode code = Db.DB.StatusCode.valueOf(res.cause().getMessage());
                LOGGER.info("wow " + code.toString());
            }
        });
    }

    private Future<String> getOkFuture(String msg) {
        Future<String> future = Future.future();

        JsonObject dataJson = new JsonObject();
        dataJson.put("hello", "world");
        redisClient.hmset("testkey:1", dataJson, res -> {
            if (res.succeeded()) {
                future.complete(msg);
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }

    private Future<String> getFailFuture(String msg) {
        Future<String> future = Future.future();

        redisClient.hgetall("testkey:2", res -> {
            if (res.succeeded()) {
                JsonObject jsonObject = res.result();
                if (jsonObject.isEmpty()) {
                    future.fail(Db.DB.StatusCode.STATUS_INVALID_PARAM.toString());
                } else {
                    future.complete(msg);
                }
            } else {
                future.fail(Db.DB.StatusCode.STATUS_INVALID_PARAM.toString());
            }
        });

        return future;
    }
}
