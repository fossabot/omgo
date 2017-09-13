package com.omgo.dbservice;

import com.omgo.dbservice.etcd.Services;
import com.omgo.dbservice.model.ModelConverter;
import com.omgo.dbservice.model.Utils;
import io.grpc.ManagedChannel;
import io.vertx.core.AsyncResult;
import io.vertx.core.Future;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.sql.SQLClient;
import io.vertx.ext.sql.SQLConnection;
import io.vertx.redis.RedisClient;
import io.vertx.redis.RedisTransaction;
import proto.Db;
import proto.SnowflakeOuterClass;
import proto.SnowflakeServiceGrpc;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ThreadLocalRandom;

public class DbOperator {

    // SQL query constants
    private static final String QUERY_USERINFO_USN = "SELECT * FROM user WHERE usn=?";
    private static final String QUERY_USERINFO_UID = "SELECT * FROM user WHERE uid=?";
    private static final String QUERY_USERINFO_EMAIL = "SELECT * FROM user WHERE email=?";

    // redis config
    private static final int REDIS_USER_EXPIRE_DURATION = 24 * 60 * 60;

    // clients
    private SQLClient sqlClient;
    private RedisClient redisClient;

    public DbOperator(SQLClient sqlClient, RedisClient redisClient) {
        this.sqlClient = sqlClient;
        this.redisClient = redisClient;
    }


    /**
     * get a unique user id from snowflake service
     * a random step (1 ~ 1000) will be added to snowflake's userid
     * this step is guarantee to be even, so the userid will always be odd
     *
     * @return
     */
    public Future<JsonObject> generateUniqueUserId() {
        Future<JsonObject> future = Future.future();

        Services.ServicePool servicePool = Services.getInstance().getServicePool();
        if (servicePool == null) {
            future.fail("service pool not initialized");
            return future;
        }

        ManagedChannel channel = servicePool.getChannel(servicePool.getServicePath(Services.SERVICE_SNOWFLAKE));
        if (channel == null) {
            future.fail("service not found: snowflake");
            return future;
        }

        // make a random user id increase step, and make sure increment is even
        // so the user id maintain odd
        long randomStep = ThreadLocalRandom.current().nextLong(1, 1000);
        if (randomStep % 2 != 0) {
            randomStep++;
        }

        SnowflakeServiceGrpc.SnowflakeServiceVertxStub stub = SnowflakeServiceGrpc.newVertxStub(channel);
        SnowflakeOuterClass.Snowflake.Param param = SnowflakeOuterClass.Snowflake.Param.newBuilder()
            .setName("userid")
            .setStep(randomStep)
            .build();

        stub.next2(param, res -> {
            if (res.succeeded()) {
                SnowflakeOuterClass.Snowflake.Value value = res.result();
                JsonObject jsonObject = new JsonObject();
                jsonObject.put(ModelConverter.KEY_UID, value.getValue());
                future.complete(jsonObject);
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }

    /**
     * Remove user info from redis
     *
     * @param usn
     * @return
     */
    public Future<JsonObject> removeUserInfoRedis(long usn) {
        Future<JsonObject> future = Future.future();
        if (usn == 0L) {
            future.fail("invalid usn");
        } else {
            redisClient.del(AccountUtils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    future.complete(new JsonObject());
                } else {
                    future.fail(res.cause());
                }
            });
        }
        return future;
    }

    /**
     * Query user info in redis
     *
     * @param usn user serial number
     * @return Future
     */
    public Future<JsonObject> queryUserInfoRedis(long usn) {
        Future<JsonObject> future = Future.future();

        if (usn == 0L) {
            future.fail("invalid usn");
        } else {
            redisClient.hgetall(AccountUtils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    if (!res.result().isEmpty()) {
                        future.complete(res.result());
                    } else {
                        future.fail(String.format("user %d not found in redis", usn));
                    }
                } else {
                    future.fail(res.cause());
                }
            });
        }

        return future;
    }

    /**
     * Query user info in MySQL
     *
     * @param json
     * @return Future
     */
    public Future<JsonObject> queryUserInfoSQL(JsonObject json) {
        Future<JsonObject> future = Future.future();

        Db.DB.UserEntry userEntry = ModelConverter.json2UserEntry(json);

        long usn = userEntry.getUsn();
        long uid = userEntry.getUid();
        String email = userEntry.getEmail();

        if (usn == 0L && uid == 0L && Utils.isEmptyString(email)) {
            future.fail("invalid query key");
        } else {
            sqlClient.getConnection(connRes -> {
                if (connRes.succeeded()) {
                    SQLConnection connection = connRes.result();

                    String query = "";
                    JsonArray params = new JsonArray();
                    if (usn != 0L) {
                        params.add(usn);
                        query = QUERY_USERINFO_USN;
                    } else if (uid != 0L) {
                        params.add(uid);
                        query = QUERY_USERINFO_UID;
                    } else if (Utils.isNotEmptyString(email)) {
                        params.add(email);
                        query = QUERY_USERINFO_EMAIL;
                    }

                    connection.queryWithParams(query, params, queryRes -> {
                        if (queryRes.succeeded()) {
                            List<JsonObject> results = queryRes.result().getRows();
                            if (results != null && results.size() > 0) {
                                future.complete(results.get(0));
                            } else {
                                future.fail("query success with no result");
                            }
                        } else {
                            future.fail(queryRes.cause());
                        }

                        connection.close();
                    });

                } else {
                    future.fail(connRes.cause());
                }
            });
        }

        return future;
    }

    /**
     * Update user info in redis
     *
     * @param userInfoJson
     * @return Future
     */
    public Future<JsonObject> updateUserInfoRedis(JsonObject userInfoJson) {
        Future<JsonObject> future = Future.future();
        if (userInfoJson == null) {
            future.fail("invalid userinfo(null)");
        } else {
            long usn = userInfoJson.getLong(ModelConverter.KEY_USN);
            String key = AccountUtils.getRedisKey(usn);
            RedisTransaction transaction = redisClient.transaction();
            transaction.multi(event -> {
                transaction.hmset(key, userInfoJson, hmsetEvent -> {
                    if (isRedisTransactionSucceed(hmsetEvent)) {
                        transaction.expire(key, REDIS_USER_EXPIRE_DURATION, expireEvent -> {
                            if (isRedisTransactionSucceed(expireEvent)) {
                                transaction.exec(execEvent -> {
                                    if (execEvent.succeeded()) {
                                        future.complete(userInfoJson);
                                    } else {
                                        future.fail(execEvent.cause());
                                    }
                                });
                            } else {
                                transaction.discard(discard -> {
                                    future.fail(expireEvent.cause());
                                });
                            }
                        });
                    } else {
                        transaction.discard(discard -> {
                            future.fail(hmsetEvent.cause());
                        });
                    }
                });
            });
        }

        return future;
    }

    /**
     * Update user info in MySQL
     *
     * @param userJson
     * @return Future
     */
    public Future<JsonObject> updateUserInfoSQL(JsonObject userJson) {
        Future<JsonObject> future = Future.future();

        long usn = userJson.getLong(ModelConverter.KEY_USN);
        if (usn == 0L) {
            future.fail("invalid usn");
            return future;
        }

        Set<String> updatableKeys = ModelConverter.getUserUpdatableMapKeySet();

        String SQL_UPDATE = "UPDATE user SET ";

        List<String> columnNameList = new ArrayList<>();
        JsonArray params = new JsonArray();

        Map<String, Object> map = userJson.getMap();
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            String key = entry.getKey();
            if (!updatableKeys.contains(key)) {
                continue;
            }
            params.add(entry.getValue());
            columnNameList.add(key + "=?");
        }

        if (columnNameList.size() == 0) {
            future.fail("update user info failed, invalid user info");
            return future;
        }

        SQL_UPDATE += String.join(",", columnNameList);
        SQL_UPDATE += " WHERE usn=?";

        params.add(usn);

        // update
        String finalSQL_UPDATE = SQL_UPDATE;
        sqlClient.getConnection(res -> {
            if (res.succeeded()) {
                SQLConnection connection = res.result();
                connection.updateWithParams(finalSQL_UPDATE, params, sqlRes -> {
                    if (sqlRes.succeeded()) {
                        connection.queryWithParams(QUERY_USERINFO_USN, new JsonArray().add(usn), queryRes -> {
                            if (queryRes.succeeded()) {
                                List<JsonObject> rows = queryRes.result().getRows();
                                if (rows.size() > 0) {
                                    future.complete(rows.get(0));
                                } else {
                                    future.fail("query after update failed");
                                }
                            } else {
                                future.fail("query after update failed");
                            }
                        });
                    } else {
                        future.fail(sqlRes.cause());
                    }
                });
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }

    /**
     * Insert userInfo to SQL
     *
     * @param userJson
     * @return
     */
    public Future<JsonObject> insertUserInfoSQL(JsonObject userJson) {
        Future<JsonObject> future = Future.future();

        long uid = userJson.getLong(ModelConverter.KEY_UID);
        String token = userJson.getString(ModelConverter.KEY_TOKEN);
        // query
        String queryQuery = ModelConverter.SQLQueryQueryUid(uid);
        // insert
        String insertQuery = ModelConverter.SQLQueryInsert(userJson);
        sqlClient.getConnection(res -> {
            if (res.succeeded()) {
                SQLConnection connection = res.result();
                connection.execute(insertQuery, insertRes -> {
                    if (insertRes.succeeded()) {
                        connection.query(queryQuery, queryRes -> {
                            if (queryRes.succeeded()) {
                                if (queryRes.result() != null && queryRes.result().getRows().size() > 0) {
                                    JsonObject result = queryRes.result().getRows().get(0);
                                    result.put(ModelConverter.KEY_TOKEN, token);
                                    future.complete(result);
                                } else {
                                    future.fail("query error");
                                }
                            } else {
                                future.fail(queryRes.cause());
                            }
                        });
                    } else {
                        future.fail(insertRes.cause());
                    }
                });
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }

    private boolean isRedisTransactionSucceed(AsyncResult<String> result) {
        if (result.succeeded() && "QUEUED".equals(result.result())) {
            return true;
        } else {
            return false;
        }
    }
}
