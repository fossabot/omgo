package com.omgo.dbservice;

import com.coreos.jetcd.KV;
import com.coreos.jetcd.data.ByteSequence;
import com.coreos.jetcd.data.KeyValue;
import com.coreos.jetcd.kv.GetResponse;
import com.coreos.jetcd.options.GetOption;
import com.omgo.dbservice.etcd.Services;
import io.grpc.ManagedChannel;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Future;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.jdbc.JDBCClient;
import io.vertx.ext.sql.SQLClient;
import io.vertx.grpc.VertxChannelBuilder;
import io.vertx.grpc.VertxServer;
import io.vertx.grpc.VertxServerBuilder;
import io.vertx.redis.RedisClient;
import io.vertx.redis.RedisOptions;
import proto.SnowflakeOuterClass;
import proto.SnowflakeServiceGrpc;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.concurrent.CompletableFuture;

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

//        testGRPC();

        testService();

//        vertx.executeBlocking(future -> {
//            future.complete(testETCD());
//        }, res -> {
//            LOGGER.info("result is: ", res.result());
//        });
    }

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

    private Random random = new Random();

    private Future<String> queryRedis() {
        Future<String> future = Future.future();
        if (random.nextBoolean()) {
            future.complete("redis ok");
        } else {
            future.fail("shit happens");
        }
        return future;
    }

    private Future<String> queryMySql() {
        Future<String> future = Future.future();
        if (random.nextBoolean()) {
            future.complete("mysql ok");
        } else {
            future.fail("mysql fuck up");
        }
        return future;
    }

    private void testCompose() {
        Future<String> futureRedis = queryRedis();
        futureRedis.compose(
            s1 -> {
                LOGGER.info(s1);
            },
            queryMySql().setHandler(res -> {
                if (res.succeeded()) {
                    LOGGER.info(res.result());
                } else {
                    LOGGER.info(res.cause());
                }
            })
        );
    }

    private void testGRPC() {

        ManagedChannel channel = VertxChannelBuilder
            .forAddress(vertx, "localhost", 40001)
            .usePlaintext(true)
            .build();

        SnowflakeServiceGrpc.SnowflakeServiceVertxStub stub = SnowflakeServiceGrpc.newVertxStub(channel);

        SnowflakeOuterClass.Snowflake.Param param = SnowflakeOuterClass.Snowflake.Param.newBuilder()
            .setName("userid")
            .setStep(1000)
            .build();

        stub.next2(param, res -> {
            if (res.succeeded()) {
                SnowflakeOuterClass.Snowflake.Value value = res.result();
                LOGGER.info(value.getValue());
            } else {
                LOGGER.error(res.cause());
            }
        });
    }

    private String testETCD() {
        String host = config().getString("etcd.host", "http://localhost:2379");
        LOGGER.info("etcd host:" + host);
        Services.getInstance().init(host);
        ByteSequence key = ByteSequence.fromString("root/");

        KV kvClient = Services.getInstance().getKVClient();
        if (kvClient != null) {
            try {
                ByteSequence endKey = Services.getRangeKey("root/");

                CompletableFuture<GetResponse> getFuture = kvClient.get(key, GetOption.newBuilder().withRange(endKey).build());
                GetResponse response = getFuture.get();
                List<KeyValue> results = response.getKvs();
                for (KeyValue kv : results) {
                    String snHost = kv.getValue().toStringUtf8();
                    String snKey = kv.getKey().toStringUtf8();
                    LOGGER.info(String.format("%s %s", snKey, snHost));
                }
                return results.get(0).getValue().toStringUtf8();
            } catch (Exception e) {
                e.printStackTrace();
            }
        }

        return "";
    }

    private void testService() {
        String host = config().getString("etcd.host", "http://localhost:2379");
        LOGGER.info("etcd host:" + host);
        Services.getInstance().init(host);

        List<String> serviceNames = new ArrayList<>();
        serviceNames.add("snowflake");

        Services.ServicePool servicePool = Services.getInstance().createServicePool(vertx,"backends", serviceNames);
        LOGGER.info("service pool created");

        ManagedChannel channel = servicePool.getChannel("backends/snowflake");
        if (channel != null) {
            SnowflakeServiceGrpc.SnowflakeServiceVertxStub stub = SnowflakeServiceGrpc.newVertxStub(channel);

            SnowflakeOuterClass.Snowflake.Param param = SnowflakeOuterClass.Snowflake.Param.newBuilder()
                .setName("userid")
                .setStep(1000)
                .build();

            stub.next2(param, res -> {
                if (res.succeeded()) {
                    SnowflakeOuterClass.Snowflake.Value value = res.result();
                    LOGGER.info(value.getValue());
                } else {
                    LOGGER.error(res.cause());
                }
            });
        } else {
            LOGGER.error("unable to find channel to snowflake");
        }
    }
}
