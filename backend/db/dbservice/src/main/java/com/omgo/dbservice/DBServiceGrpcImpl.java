package com.omgo.dbservice;

import com.omgo.dbservice.etcd.Services;
import com.omgo.dbservice.model.ModelConverter;
import com.omgo.dbservice.model.Utils;
import com.sun.xml.internal.fastinfoset.stax.events.Util;
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
import proto.Db.DB;
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

    // SQL query constants
    private static final String QUERY_USERINFO_USN = "SELECT * FROM user WHERE usn=?";
    private static final String QUERY_USERINFO_UID = "SELECT * FROM user WHERE uid=?";
    private static final String QUERY_USERINFO_EMAIL = "SELECT * FROM user WHERE email=?";

    // gRPC responses
    private static final DB.Result dbOkResult = DbProtoUtils.makeOkResult();
    private static final DB.Result userNotFoundResult = DbProtoUtils.makeResult(DB.StatusCode.STATUS_USER_NOT_FOUND);

    private static final Logger LOGGER = LoggerFactory.getLogger(DBServiceGrpcImpl.class);

    // clients
    private SQLClient sqlClient;
    private RedisClient redisClient;

    //
    public DBServiceGrpcImpl(SQLClient sqlClient, RedisClient redisClient) {
        this.sqlClient = sqlClient;
        this.redisClient = redisClient;
    }

    @Override
    public void userQuery(DB.UserKey request, Future<DB.UserOpResult> response) {
        LOGGER.info("userQuery: " + request);

        // query success handler
        Handler<DB.UserOpResult> successHandler = response::complete;

        // query in redis then in mysql
        Future<DB.UserExtendInfo> redisFuture = queryUserInfoRedis(request.getUsn());
        redisFuture.setHandler(res -> {
            if (res.succeeded()) {
                DB.UserExtendInfo extendInfo = res.result();
                LOGGER.info(String.format("redis hit for user:%d", extendInfo.getInfo().getUsn()));
                successHandler.handle(DbProtoUtils.makeUserOpOkResult(extendInfo));
            } else {
                Future<JsonObject> mysqlFuture = queryUserInfoSQL(request);
                mysqlFuture.setHandler(sqlRes -> {
                    if (sqlRes.succeeded()) {
                        // update redis
                        JsonObject userJson = sqlRes.result();
                        Future<JsonObject> updateRedisFuture = updateUserInfoRedis(userJson);
                        updateRedisFuture.setHandler(updateRedisRes -> {
                            if (updateRedisRes.failed()) {
                                LOGGER.info(updateRedisRes.cause());
                            }
                            // response
                            DB.UserOpResult result = DbProtoUtils.makeUserOpOkResult(ModelConverter.json2UserInfo(userJson));
                            successHandler.handle(result);
                        });
                    } else {
                        // query failed
                        LOGGER.warn("user query failed in both redis and mysql");
                        response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_NOT_FOUND, ""));
                    }
                });
            }
        });
    }

    @Override
    public void userUpdateInfo(Common.UserInfo request, Future<DB.Result> response) {
        LOGGER.info("userUpdate: " + request);

        Future<Common.UserInfo> updateSQLFuture = updateUserInfoSQL(request);
        updateSQLFuture.setHandler(res -> {
            if (res.succeeded()) {
                Common.UserInfo userInfo = res.result();
                Future<JsonObject> redisFuture = updateUserInfoRedis(ModelConverter.userInfo2Json(userInfo));
                redisFuture.setHandler(redisRes -> {
                    if (redisRes.succeeded()) {
                        response.complete(dbOkResult);
                    } else {
                        LOGGER.error("update user info redis failed:" + redisRes.cause());
                        response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INTERNAL_ERROR, redisRes.cause().toString()));
                    }
                });
            } else {
                LOGGER.error("update user info failed:" + res.cause());
                response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INTERNAL_ERROR, res.cause().toString()));
            }
        });
    }

    @Override
    public void userRegister(DB.UserExtendInfo request, Future<DB.UserOpResult> response) {
        LOGGER.info("userRegister: " + request);

        Common.UserInfo userInfo = request.getInfo();
        if (userInfo == null) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            LOGGER.error("invalid user info(null)");
            return;
        }

        String email = userInfo.getEmail();
        if (!AccountUtils.isValidEmailAddress(email)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_EMAIL));
            LOGGER.error("invalid email address");
            return;
        }

        String nickname = userInfo.getNickname();
        if (Utils.isEmptyString(nickname)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            LOGGER.error("invalid nickname");
            return;
        }

        String secret = request.getSecret();
        if (AccountUtils.isValidSecret(secret)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_SECRET));
            LOGGER.error("invalid password");
            return;
        }

        DB.UserExtendInfo.Builder extendInfoBuilder = DB.UserExtendInfo.newBuilder();

        // check if user with email already exist
        DB.UserKey userKey = DB.UserKey.newBuilder()
            .setEmail(email)
            .build();

        // get user id
        Future<Long> snowflakeFuture = generateUniqueUserId();

        Future<JsonObject> sqlFuture = queryUserInfoSQL(userKey);
        // insert into mysql
        sqlFuture.setHandler(sqlRes -> {
            // email already exist
            if (sqlRes.succeeded()) {
                LOGGER.error("register failed, user with email:" + email + " already existed");
                response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_ALREADY_EXIST));
            } else {
                // generate user id
                snowflakeFuture.setHandler(res -> {
                    if (res.succeeded()) {
                        // user id
                        long userId = res.result();
                        byte[] saltRaw = AccountUtils.getSalt();
                        byte[] tokenRaw = AccountUtils.getToken(saltRaw);
                        String salt = AccountUtils.encodeBase64(saltRaw);
                        String token = AccountUtils.encodeBase64(tokenRaw);
                        String saltedSecret = AccountUtils.saltedSecret(secret, salt);
                        JsonObject jsonObject = ModelConverter.userInfo2Json(userInfo);
                        jsonObject.put(ModelConverter.KEY_UID, userId);
                        jsonObject.put(ModelConverter.KEY_TOKEN, token);
                        jsonObject.put(ModelConverter.KEY_SALT, salt);
                        jsonObject.put(ModelConverter.KEY_SECRET, saltedSecret);

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
                                        response.complete(DbProtoUtils.makeUserOpOkResult(extendInfoBuilder.build()));
                                    } else {
                                        LOGGER.error("update userInfo redis failed:" + redisRes.cause());
                                        response.complete(DbProtoUtils.makeUserOpInternalFailedResult(redisRes.cause().toString()));
                                    }
                                });
                            } else {
                                LOGGER.error(insertRes.cause());
                                response.complete(DbProtoUtils.makeUserOpInternalFailedResult(insertRes.cause().toString()));
                            }
                        });

                    } else {
                        response.complete(DbProtoUtils.makeUserOpInternalFailedResult(res.cause().toString()));
                    }
                });
            }
        });
    }

    @Override
    public void userLogin(DB.UserExtendInfo request, Future<DB.UserOpResult> response) {
        LOGGER.info("userLogin: " + request);

        Common.UserInfo userInfo = request.getInfo();
        if (userInfo == null) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            return;
        }

        long queryUsn = userInfo.getUsn();
        String queryEmail = userInfo.getEmail();
        String querySecret = request.getSecret();
        String queryToken = request.getToken();

        // 1. usn + token
        if (queryUsn != 0L && !Util.isEmptyString(queryToken)) {
            Future<DB.UserExtendInfo> redisFuture = queryUserInfoRedis(userInfo.getUsn());
            redisFuture.setHandler(res -> {
               if (res.succeeded()) {
                   DB.UserExtendInfo extendInfo = res.result();
                   // check token
                   if (!queryToken.equals(extendInfo.getToken())) {
                       response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_TOKEN));
                   } else {
                       response.complete(DbProtoUtils.makeUserOpOkResult(extendInfo));
                   }
               } else {
                   response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_NOT_FOUND));
               }
            });
        } else if (AccountUtils.isValidEmailAddress(queryEmail) && AccountUtils.isValidSecret(querySecret)) {
            // 2. email + secret
            DB.UserKey key = DB.UserKey.newBuilder().setEmail(queryEmail).build();
            Future<JsonObject> sqlFuture = queryUserInfoSQL(key);
            sqlFuture.setHandler(res -> {
                if (res.succeeded()) {
                    JsonObject userJson = res.result();
                    String salt = userJson.getString(ModelConverter.KEY_SALT);
                    String saltedQuerySecret = AccountUtils.saltedSecret(querySecret, salt);
                    if (!Utils.isEmptyString(saltedQuerySecret) && saltedQuerySecret.equals(userJson.getString(ModelConverter.KEY_SECRET))) {
                        byte[] saltRaw = AccountUtils.decodeBase64(salt);
                        byte[] tokenRaw = AccountUtils.getToken(saltRaw);
                        String token = AccountUtils.encodeBase64(tokenRaw);
                        userJson.put(ModelConverter.KEY_TOKEN, token);

                        // update redis
                        Future<JsonObject> redisUpdateFuture = updateUserInfoRedis(userJson);
                        redisUpdateFuture.setHandler(updateRedisRes -> {
                           if (updateRedisRes.succeeded()) {
                               Common.UserInfo retUserInfo = ModelConverter.json2UserInfo(userJson);
                               DB.UserExtendInfo extendInfo = DbProtoUtils.makeUserExtendInfo(retUserInfo, querySecret, queryToken);
                               response.complete(DbProtoUtils.makeUserOpOkResult(extendInfo));
                           } else {
                               response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INTERNAL_ERROR));
                           }
                        });
                    } else {
                        response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_SECRET));
                    }
                } else {
                    response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_NOT_FOUND));
                }
            });
        } else {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
        }
    }

    @Override
    public void userLogout(DB.UserLogoutRequest request, Future<DB.Result> response) {
        LOGGER.info("userLogout: " + request);

        long usn = request.getUsn();
        String token = request.getToken();
        if (usn == 0L) {
            response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INVALID_USN));
            return;
        }

        if (Utils.isEmptyString(token)) {
            response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INVALID_TOKEN));
            return;
        }

        Future<DB.UserExtendInfo> redisFuture = queryUserInfoRedis(usn);
        redisFuture.setHandler(res -> {
            if (res.failed()) {
                response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_USER_NOT_FOUND));
            } else {
                DB.UserExtendInfo userExtendInfo = res.result();
                if (!token.equals(userExtendInfo.getToken())) {
                    response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INVALID_TOKEN));
                } else {
                    Future<Void> delFuture = removeUserInfoRedis(usn);
                    delFuture.setHandler(removeRes -> {
                        if (removeRes.succeeded()) {
                            response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_OK));
                        } else {
                            response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INTERNAL_ERROR));
                        }
                    });
                }
            }
        });
    }

    @Override
    public void userExtraInfoQuery(DB.UserKey request, Future<DB.UserOpResult> response) {
        long usn = request.getUsn();
        long uid = request.getUid();
        String email = request.getEmail();

        if (usn == 0L && (uid == 0L || Utils.isEmptyString(email))) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            return;
        }

        Future<DB.UserExtendInfo> redisQueryFuture = queryUserInfoRedis(usn);
        if (usn != 0) {
            redisQueryFuture.fail("invalid usn");
        }

        redisQueryFuture.setHandler(redisQueryRes -> {
           if (redisQueryRes.succeeded()) {
               response.complete(DbProtoUtils.makeUserOpOkResult(redisQueryRes.result()));
           } else {
               Future<JsonObject> mysqlFuture = queryUserInfoSQL(request);
               mysqlFuture.setHandler(mysqlRes -> {
                   if (mysqlRes.succeeded()) {
                       JsonObject mysqlJson = mysqlRes.result();
                       Common.UserInfo userInfo = ModelConverter.json2UserInfo(mysqlJson);
                       String secret = mysqlJson.getString(ModelConverter.KEY_SECRET);
                       DB.UserExtendInfo userExtendInfo = DbProtoUtils.makeUserExtendInfo(userInfo, secret, "");
                       response.complete(DbProtoUtils.makeUserOpOkResult(userExtendInfo));
                   } else {
                       response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_NOT_FOUND));
                   }
               });
           }
        });
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
     * Remove user info from redis
     *
     * @param usn
     * @return
     */
    private Future<Void> removeUserInfoRedis(long usn) {
        Future<Void> future = Future.future();
        if (usn == 0L) {
            future.fail("invalid usn");
        } else {
            redisClient.del(AccountUtils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    future.complete();
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
    private Future<DB.UserExtendInfo> queryUserInfoRedis(long usn) {
        Future<DB.UserExtendInfo> future = Future.future();

        if (usn == 0L) {
            future.fail("invalid usn");
        } else {
            redisClient.hgetall(AccountUtils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    JsonObject jsonObject = res.result();
                    String secret = jsonObject.getString(ModelConverter.KEY_SECRET);
                    String token = jsonObject.getString(ModelConverter.KEY_TOKEN);
                    future.complete(DbProtoUtils.makeUserExtendInfo(ModelConverter.json2UserInfo(jsonObject), secret, token));
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
    private Future<JsonObject> queryUserInfoSQL(DB.UserKey userKey) {
        Future<JsonObject> future = Future.future();

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
    private Future<JsonObject> updateUserInfoRedis(JsonObject userInfoJson) {
        Future<JsonObject> future = Future.future();
        if (userInfoJson == null) {
            future.fail("invalid userinfo(null)");
        } else {
            long usn = userInfoJson.getLong(ModelConverter.KEY_USN);
            redisClient.hmset(AccountUtils.getRedisKey(usn), userInfoJson, res -> {
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
        if (AccountUtils.isValidEmailAddress(userInfo.getEmail())) {
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

        long uid = userJson.getLong(ModelConverter.KEY_UID);
        // query
        String queryQuery = ModelConverter.SQLQueryQueryUid(uid);
        // insert
        String insertQuery = ModelConverter.SQLQueryInsert(userJson);
        LOGGER.info(insertQuery);
        sqlClient.getConnection(res -> {
            if (res.succeeded()) {
                SQLConnection connection = res.result();
                connection.execute(insertQuery, insertRes -> {
                    if (insertRes.succeeded()) {
                        connection.query(queryQuery, queryRes -> {
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
