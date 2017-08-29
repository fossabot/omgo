package com.omgo.dbservice;

import com.omgo.dbservice.etcd.Services;
import com.omgo.dbservice.model.ModelConverter;
import com.omgo.dbservice.model.Utils;
import io.grpc.ManagedChannel;
import io.vertx.core.Future;
import io.vertx.core.Handler;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.sql.SQLClient;
import io.vertx.ext.sql.SQLConnection;
import io.vertx.redis.RedisClient;
import proto.DBServiceGrpc;
import proto.Db;
import proto.SnowflakeOuterClass;
import proto.SnowflakeServiceGrpc;
import proto.common.Common;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ThreadLocalRandom;

/**
 * Database gRPC service implementation
 * <p>
 * Created by mg on 17/07/2017.
 */
public class DBServiceGrpcImpl extends DBServiceGrpc.DBServiceVertxImplBase {

    private static final String QUERY_USERINFO_USN = "SELECT * FROM user WHERE usn=?";
    private static final String QUERY_USERINFO_UID = "SELECT * FROM user WHERE uid=?";
    private static final String QUERY_USERINFO_EMAIL = "SELECT * FROM user WHERE email=?";

    private static final Logger LOGGER = LoggerFactory.getLogger(DBServiceGrpcImpl.class);

    private SQLClient sqlClient;
    private RedisClient redisClient;

    public DBServiceGrpcImpl(SQLClient sqlClient, RedisClient redisClient) {
        this.sqlClient = sqlClient;
        this.redisClient = redisClient;
    }

    @Override
    public void userQuery(Db.DB.UserKey request, Future<Common.UserInfo> response) {
        LOGGER.info("userQuery:" + request);

        // query success handler
        Handler<Common.UserInfo> successHandler = response::complete;

        // query in redis then in mysql
        Future<Common.UserInfo> redisFuture = queryUserInfoRedis(request.getUsn());
        redisFuture.setHandler(res -> {
            if (res.succeeded()) {
                LOGGER.info(String.format("redis hit for user:%d", res.result().getUsn()));
                successHandler.handle(res.result());
            } else {
                Future<Common.UserInfo> mysqlFuture = queryUserInfoSQL(request);
                mysqlFuture.setHandler(sqlRes -> {
                    if (sqlRes.succeeded()) {
                        // update redis
                        Common.UserInfo userInfo = sqlRes.result();
                        Future<JsonObject> updateRedisFuture = updateUserInfoRedis(ModelConverter.userInfo2Json(userInfo));
                        updateRedisFuture.setHandler(updateRedisRes -> {
                            if (updateRedisRes.failed()) {
                                LOGGER.info(updateRedisRes.cause());
                            }
                            // response
                            successHandler.handle(sqlRes.result());
                        });
                    } else {
                        response.fail("user query failed in both redis and mysql");
                    }
                });
            }
        });
    }

    @Override
    public void userUpdateInfo(Common.UserInfo request, Future<Db.DB.NullValue> response) {
        LOGGER.info("userUpdate:" + request);

        Future<Common.UserInfo> updateSQLFuture = updateUserInfoSQL(request);
        updateSQLFuture.setHandler(res -> {
            if (res.succeeded()) {
                Common.UserInfo userInfo = res.result();
                Future<JsonObject> redisFuture = updateUserInfoRedis(ModelConverter.userInfo2Json(userInfo));
                redisFuture.setHandler(redisRes -> {
                    if (redisRes.succeeded()) {
                        Common.RspHeader header = Common.RspHeader.newBuilder()
                            .setStatus(Common.ResultCode.RESULT_OK_VALUE)
                            .build();
                        response.complete(Db.DB.NullValue.newBuilder().build());
                    } else {
                        LOGGER.error("update user info redis failed:" + redisRes.cause());
                        response.fail("user update redis failed");
                    }
                });
            } else {
                LOGGER.error("update user info failed:" + res.cause());
                response.fail("user update failed");
            }
        });
    }

    @Override
    public void userRegister(Db.DB.UserExtendInfo request, Future<Db.DB.UserExtendInfo> response) {
        LOGGER.info("userRegister", request);

        Common.UserInfo userInfo = request.getInfo();
        if (userInfo == null) {
            response.fail("invalid user info(null)");
            LOGGER.error("invalid user info(null)");
            return;
        }

        String email = userInfo.getEmail();
        if (!Utils.isValidEmailAddress(email)) {
            response.fail("invalid email address");
            LOGGER.error("invalid email address");
            return;
        }

        String nickname = userInfo.getNickname();
        if (Utils.isEmptyString(nickname)) {
            response.fail("invalid nickname");
            LOGGER.error("invalid nickname");
            return;
        }

        String secret = request.getSecret();
        if (Utils.isEmptyString(secret) || secret.length() < AccountUtils.PASSWORD_MIN_LEN) {
            response.fail("invalid password");
            LOGGER.error("invalid password");
            return;
        }

        Db.DB.UserExtendInfo.Builder extendInfoBuilder = Db.DB.UserExtendInfo.newBuilder();
        Common.UserInfo.Builder userInfoBuilder = Common.UserInfo.newBuilder();

        // check if user with email already exist
        Db.DB.UserKey userKey = Db.DB.UserKey.newBuilder()
            .setEmail(email)
            .build();

        // get user id
        Future<Long> snowflakeFuture = generateUniqueUserId();

        Future<Common.UserInfo> sqlFuture = queryUserInfoSQL(userKey);
        // insert into mysql
        sqlFuture.setHandler(sqlRes -> {
            // email already exist
            if (sqlRes.succeeded()) {
                LOGGER.error("register failed, user with email:" + email + " already existed");
                response.fail("email has already been registered");
                return;
            } else {
                // generate user id
                snowflakeFuture.setHandler(res -> {
                    if (res.succeeded()) {
                        // user id
                        long userId = res.result();
                        byte[] saltRaw = AccountUtils.getSalt();
                        byte[] tokenRaw = AccountUtils.getToken(saltRaw);
                        String salt = AccountUtils.base64(saltRaw);
                        String token = AccountUtils.base64(tokenRaw);
                        JsonObject jsonObject = ModelConverter.userInfo2Json(userInfo);
                        jsonObject.put(ModelConverter.KEY_TOKEN, token);
                        jsonObject.put(ModelConverter.KEY_SALT, salt);
                        jsonObject.put(ModelConverter.KEY_SECRET, secret);

                        Future<JsonObject> insertFuture = insertUserInfoSQL(jsonObject);
                        // actual insert
                        insertFuture.setHandler(insertRes -> {
                            if (insertRes.succeeded()) {
                                JsonObject resultJson = insertRes.result();
                                resultJson.put(ModelConverter.KEY_TOKEN, token);
                                Future<JsonObject> redisFuture = updateUserInfoRedis(resultJson);
                                redisFuture.setHandler(redisRes -> {
                                    if (redisRes.succeeded()) {
                                        // gRPC response
                                        Common.UserInfo finalUserInfo = ModelConverter.json2UserInfo(resultJson);
                                        extendInfoBuilder.setInfo(finalUserInfo)
                                            .setToken(token)
                                            .setSecret(secret);
                                        response.complete(extendInfoBuilder.build());
                                    } else {
                                        LOGGER.error("update userInfo redis failed:" + redisRes.cause());
                                        response.fail(redisRes.cause());
                                    }
                                });
                            } else {
                                LOGGER.error(insertRes.cause());
                                response.fail(insertRes.cause());
                            }
                        });

                    } else {
                        response.fail(res.cause());
                    }
                });
            }
        });
    }

    @Override
    public void userLogin(Db.DB.UserExtendInfo request, Future<Db.DB.UserExtendInfo> response) {
        super.userLogin(request, response);
    }

    @Override
    public void userLogout(Db.DB.UserLogoutRequest request, Future<Db.DB.NullValue> response) {
        super.userLogout(request, response);
    }

    @Override
    public void userExtraInfoQuery(Db.DB.UserKey request, Future<Db.DB.UserExtendInfo> response) {
        super.userExtraInfoQuery(request, response);
    }

    /**
     * get a unique user id from snowflake service
     * a random step (1 ~ 1000) will be added to snowflake's userid
     * this step is guarantee to be even, so the userid will always be odd
     *
     * @return
     */
    private Future<Long> generateUniqueUserId() {
        Future<Long> future = Future.future();

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
                future.complete(value.getValue());
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }


    /**
     * Query user info in redis
     *
     * @param usn user serial number
     * @return Future
     */
    private Future<Common.UserInfo> queryUserInfoRedis(long usn) {
        Future<Common.UserInfo> future = Future.future();

        if (usn == 0L) {
            future.fail("invalid usn");
        } else {
            redisClient.hgetall(Utils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    future.complete(ModelConverter.json2UserInfo(res.result()));
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
     * @param userKey User key with usn/uid/email
     * @return Future
     */
    private Future<Common.UserInfo> queryUserInfoSQL(Db.DB.UserKey userKey) {
        Future<Common.UserInfo> future = Future.future();

        long usn = userKey.getUsn();
        long uid = userKey.getUid();
        String email = userKey.getEmail();

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
                    } else if (!Utils.isEmptyString(email)) {
                        params.add(email);
                        query = QUERY_USERINFO_EMAIL;
                    }

                    connection.queryWithParams(query, params, queryRes -> {
                        if (queryRes.succeeded()) {
                            List<JsonObject> results = queryRes.result().getRows();
                            if (results != null && results.size() > 0) {
                                future.complete(ModelConverter.json2UserInfo(results.get(0)));
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
    private Future<JsonObject> updateUserInfoRedis(JsonObject userInfoJson) {
        Future<JsonObject> future = Future.future();
        if (userInfoJson == null) {
            future.fail("invalid userinfo(null)");
        } else {
            long usn = userInfoJson.getLong(ModelConverter.KEY_USN);
            redisClient.hmset(Utils.getRedisKey(usn), userInfoJson, res -> {
                if (res.succeeded()) {
                    future.complete(userInfoJson);
                } else {
                    future.fail(res.cause());
                }
            });
        }

        return future;
    }

    /**
     * Update user info in MySQL
     *
     * @param userInfo
     * @return Future
     */
    private Future<Common.UserInfo> updateUserInfoSQL(Common.UserInfo userInfo) {
        Future<Common.UserInfo> future = Future.future();

        if (userInfo.getUsn() == 0L) {
            future.fail("invalid usn");
            return future;
        }

        String SQL_UPDATE = "UPDATE user SET ";

        List<String> columnNameList = new ArrayList<>();
        JsonArray params = new JsonArray();

        if (!Utils.isEmptyString(userInfo.getAvatar())) {
            params.add(userInfo.getAvatar());
            columnNameList.add(ModelConverter.KEY_AVATAR + "=?");
        }
        if (userInfo.getBirthday() != 0L) {
            params.add(userInfo.getBirthday());
            columnNameList.add(ModelConverter.KEY_BIRTHDAY + "=?");
        }
        if (!Utils.isEmptyString(userInfo.getCountry())) {
            params.add(userInfo.getCountry());
            columnNameList.add(ModelConverter.KEY_COUNTRY + "=?");
        }
        if (Utils.isValidEmailAddress(userInfo.getEmail())) {
            params.add(userInfo.getEmail());
            columnNameList.add(ModelConverter.KEY_EMAIL + "=?");
        }
        if (userInfo.getGender() != Common.Gender.GENDER_UNKNOWN) {
            params.add(userInfo.getGenderValue());
            columnNameList.add(ModelConverter.KEY_GENDER + "=?");
        }
        if (!Utils.isEmptyString(userInfo.getNickname())) {
            params.add(userInfo.getNickname());
            columnNameList.add(ModelConverter.KEY_NICKNAME + "=?");
        }

        if (columnNameList.size() == 0) {
            future.fail("update user info failed, invalid user info");
        }

        SQL_UPDATE += String.join(",", columnNameList);
        SQL_UPDATE += " WHERE usn=?";

        params.add(userInfo.getUsn());

        // update
        String finalSQL_UPDATE = SQL_UPDATE;
        sqlClient.getConnection(res -> {
            if (res.succeeded()) {
                SQLConnection connection = res.result();
                connection.updateWithParams(finalSQL_UPDATE, params, sqlRes -> {
                    if (sqlRes.succeeded()) {
                        connection.queryWithParams(QUERY_USERINFO_USN, new JsonArray().add(userInfo.getUsn()), queryRes -> {
                            if (queryRes.succeeded()) {
                                List<JsonObject> rows = queryRes.result().getRows();
                                if (rows.size() > 0) {
                                    future.complete(ModelConverter.json2UserInfo(rows.get(0)));
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
    private Future<JsonObject> insertUserInfoSQL(JsonObject userJson) {
        Future<JsonObject> future = Future.future();

        String SQL_INSERT = "INSERT INTO user (";
        String SQL_VALUES = "";

        List<String> VALUE_KEYS = new ArrayList<>();
        VALUE_KEYS.add(ModelConverter.KEY_UID);
        VALUE_KEYS.add(ModelConverter.KEY_AVATAR);
        VALUE_KEYS.add(ModelConverter.KEY_BIRTHDAY);
        VALUE_KEYS.add(ModelConverter.KEY_COUNTRY);
        VALUE_KEYS.add(ModelConverter.KEY_EMAIL);
        VALUE_KEYS.add(ModelConverter.KEY_GENDER);
        VALUE_KEYS.add(ModelConverter.KEY_LAST_LOGIN);
        VALUE_KEYS.add(ModelConverter.KEY_LOGIN_COUNT);
        VALUE_KEYS.add(ModelConverter.KEY_NICKNAME);
        VALUE_KEYS.add(ModelConverter.KEY_SALT);
        VALUE_KEYS.add(ModelConverter.KEY_SECRET);
        VALUE_KEYS.add(ModelConverter.KEY_SINCE);

        SQL_INSERT += String.join(",", VALUE_KEYS) + ") VALUES (";

        final long uid = userJson.getLong(ModelConverter.KEY_UID);
        SQL_INSERT += userJson.getLong(ModelConverter.KEY_UID) + ",";
        SQL_INSERT += userJson.getString(ModelConverter.KEY_AVATAR) + ",";
        SQL_INSERT += userJson.getLong(ModelConverter.KEY_BIRTHDAY) + ",";
        SQL_INSERT += userJson.getString(ModelConverter.KEY_COUNTRY) + ",";
        SQL_INSERT += userJson.getString(ModelConverter.KEY_EMAIL) + ",";
        SQL_INSERT += userJson.getInteger(ModelConverter.KEY_GENDER) + ",";
        SQL_INSERT += userJson.getLong(ModelConverter.KEY_LAST_LOGIN) + ",";
        SQL_INSERT += userJson.getLong(ModelConverter.KEY_LOGIN_COUNT) + ",";
        SQL_INSERT += userJson.getString(ModelConverter.KEY_NICKNAME) + ",";
        SQL_INSERT += userJson.getString(ModelConverter.KEY_SALT) + ",";
        SQL_INSERT += userJson.getString(ModelConverter.KEY_SECRET) + ",";
        SQL_INSERT += userJson.getLong(ModelConverter.KEY_SINCE) + ")";

        // insert
        String finalSQL_INSERT = SQL_INSERT;
        sqlClient.getConnection(res -> {
            if (res.succeeded()) {
                SQLConnection connection = res.result();
                connection.execute(finalSQL_INSERT, insertRes -> {
                    if (insertRes.succeeded()) {
                        connection.query("SELECT * FROM user WHERE uid=" + uid, queryRes -> {
                            if (queryRes.succeeded()) {
                                if (queryRes.result() != null && queryRes.result().getRows().size() > 0) {
                                    future.complete(queryRes.result().getRows().get(0));
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
}
