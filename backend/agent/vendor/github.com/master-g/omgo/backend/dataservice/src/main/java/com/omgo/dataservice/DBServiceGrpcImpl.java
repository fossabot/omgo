package com.omgo.dataservice;

import com.omgo.dataservice.model.ModelConverter;
import com.omgo.utils.ModelKeys;
import com.omgo.utils.Utils;
import io.vertx.core.Future;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.mongo.MongoClient;
import io.vertx.redis.RedisClient;
import proto.DBServiceGrpc;
import proto.Db.DB;

/**
 * Database gRPC service implementation
 * <p>
 * Created by mg on 17/07/2017.
 */
public class DBServiceGrpcImpl extends DBServiceGrpc.DBServiceVertxImplBase {
    // gRPC responses
    private static final DB.Result dbOkResult = DbProtoUtils.makeOkResult();
    private static final DB.Result userNotFoundResult = DbProtoUtils.makeResult(DB.StatusCode.STATUS_USER_NOT_FOUND);

    private static final Logger LOGGER = LoggerFactory.getLogger(DBServiceGrpcImpl.class);

    private DbOperator dbOperator;

    //
    public DBServiceGrpcImpl(MongoClient mongoClient, RedisClient redisClient) {
        dbOperator = new DbOperator(mongoClient, redisClient);
    }

    @Override
    public void userQuery(DB.UserEntry request, Future<DB.UserOpResult> response) {
        LOGGER.info("userQuery: " + request);

        Future<JsonObject> responseFuture = Future.future();
        responseFuture.setHandler(res -> {
            if (res.succeeded()) {
                response.complete(DbProtoUtils.makeUserOpOkResult(res.result()));
            } else {
                LOGGER.warn("user query failed in both redis and mongodb");
                response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_NOT_FOUND, ""));
            }
        });

        // query in redis then in mongodb
        Future<JsonObject> redisFuture = dbOperator.queryUserInfoRedis(request.getUsn());
        redisFuture.setHandler(queryRedis -> {
            if (queryRedis.succeeded()) {
                // found in redis
                JsonObject result = queryRedis.result();
                LOGGER.info(String.format("redis hit for user:%d", result.getLong(ModelKeys.USN)));
                responseFuture.complete(result);
            } else {
                // query in mongodb
                Future<JsonObject> dbFuture = dbOperator.queryUserInfoDB(ModelConverter.userEntry2Json(request));
                dbFuture.compose(queryDB -> {
                    // update in redis
                    Future<JsonObject> updateFuture = dbOperator.updateUserInfoRedis(queryDB);
                    updateFuture.setHandler(responseFuture.completer());
                }, responseFuture);
            }
        });
    }

    @Override
    public void userUpdateInfo(DB.UserEntry request, Future<DB.Result> response) {
        LOGGER.info("userUpdate: " + request);

        Future<JsonObject> updateDBFuture = dbOperator.updateUserInfoDB(ModelConverter.userEntry2Json(request));
        updateDBFuture.setHandler(res -> {
            if (res.succeeded()) {
                Future<JsonObject> redisFuture = dbOperator.updateUserInfoRedis(res.result());
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
    public void userRegister(DB.UserEntry request, Future<DB.UserOpResult> response) {
        LOGGER.info("userRegister: " + request);

        String email = request.getEmail();
        if (!Utils.isValidEmailAddress(email)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_EMAIL));
            LOGGER.error("invalid email address");
            return;
        }

        String nickname = request.getNickname();
        if (Utils.isEmptyString(nickname)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            LOGGER.error("invalid nickname");
            return;
        }

        String secret = request.getSecret();
        if (!AccountUtils.isValidSecret(secret)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_SECRET));
            LOGGER.error("invalid password");
            return;
        }

        // TODO: 14/09/2017 check more parameters

        // response future
        Future<JsonObject> responseFuture = Future.future();
        responseFuture.setHandler(res -> {
            if (res.succeeded()) {
                JsonObject resultJson = res.result();
                ModelConverter.removeKeysForLoginResponse(resultJson);
                response.complete(DbProtoUtils.makeUserOpOkResult(resultJson));
            } else {
                LOGGER.error(res.cause());
                response.complete(DbProtoUtils.makeUserOpInternalFailedResult(res.cause().toString()));
            }
        });

        // check if user with email already exist
        Future<JsonObject> dbFuture = dbOperator.queryUserInfoDB(new JsonObject().put(ModelKeys.EMAIL, email));

        dbFuture.setHandler(dbRes -> {
            // email already exist
            if (dbRes.succeeded()) {
                LOGGER.error("register failed, user with email:" + email + " already existed");
                response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_ALREADY_EXIST));
            } else {
                // generate user id
                Future<JsonObject> usnFuture = dbOperator.mongodbGenerateUsn();
                usnFuture.compose(usnRes -> {
                    long now = System.currentTimeMillis();
                    long usn = usnRes.getLong(ModelKeys.USN);
                    long uid = usnRes.getLong(ModelKeys.UID);
                    byte[] tokenRaw = AccountUtils.getToken();
                    String token = Utils.encodeBase64(tokenRaw);
                    String saltedSecret = AccountUtils.saltedSecret(secret, now);
                    JsonObject jsonObject = ModelConverter.userEntry2Json(request);
                    jsonObject.put(ModelKeys.USN, usn);
                    jsonObject.put(ModelKeys.UID, uid);
                    jsonObject.put(ModelKeys.TOKEN, token);
                    jsonObject.put(ModelKeys.SECRET, saltedSecret);
                    jsonObject.put(ModelKeys.SINCE, now);
                    jsonObject.put(ModelKeys.LAST_LOGIN, now);
                    jsonObject.put(ModelKeys.LOGIN_COUNT, 1);

                    return dbOperator.insertUserInfoDB(jsonObject);
                }).compose(insertRes -> {
                    Future<JsonObject> updateRedisFuture = dbOperator.updateUserInfoRedis(insertRes);
                    updateRedisFuture.setHandler(responseFuture.completer());
                }, responseFuture);
            }
        });
    }

    @Override
    public void userLogin(DB.UserEntry request, Future<DB.UserOpResult> response) {
        LOGGER.info("userLogin: " + request);

        long queryUsn = request.getUsn();
        String queryEmail = request.getEmail();
        String querySecret = request.getSecret();
        String queryToken = request.getToken();

        if (Utils.isEmptyString(queryEmail)
            && Utils.isEmptyString(querySecret)
            && queryUsn == 0L
            && Utils.isEmptyString(queryToken)) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            return;
        }

        // response
        Future<JsonObject> responseFuture = Future.future();
        responseFuture.setHandler(res -> {
            if (res.succeeded()) {
                response.complete(DbProtoUtils.makeUserOpOkResult(ModelConverter.removeKeysForLoginResponse(res.result())));
            } else {
                LOGGER.warn(res.cause());
                DB.StatusCode code;
                try {
                    code = DB.StatusCode.valueOf(res.cause().toString());
                } catch (IllegalArgumentException e) {
                    code = DB.StatusCode.STATUS_USER_NOT_FOUND;
                }
                response.complete(DbProtoUtils.makeUserOpResult(code));
            }
        });

        // entry
        Future<JsonObject> entry = Future.future();
        entry
            .compose(v -> {
                if (Utils.isNotEmptyString(queryToken) && queryUsn != 0L) {
                    // query user info in redis
                    return dbOperator.queryUserInfoRedis(queryUsn);
                } else {
                    // no token info, pass
                    JsonObject emptyJson = new JsonObject();
                    Future<JsonObject> future = Future.future();
                    future.complete(emptyJson);
                    return future;
                }
            })
            .compose(v -> {
                String token = v.getString(ModelKeys.TOKEN);
                if (Utils.isNotEmptyString(queryToken) && queryToken.equals(token)) {
                    // login with token success
                    Future<JsonObject> tokenLoginFuture = Future.future();
                    tokenLoginFuture.complete(v);
                    return tokenLoginFuture;
                } else if (Utils.isEmptyString(queryEmail) && Utils.isEmptyString(querySecret)) {
                    // login with token failed, and no other info
                    return Future.failedFuture(DB.StatusCode.STATUS_INVALID_PARAM.toString());
                } else {
                    // login with email and secret
                    Future<JsonObject> future = Future.future();
                    future.complete(new JsonObject());
                    return future;
                }
            })
            .compose(v -> {
                // query in mongodb and update
                Future<JsonObject> queryFuture = dbOperator.queryUserInfoDB(ModelConverter.userEntry2Json(request));
                return queryFuture.compose(queryRes -> {
                    long salt = queryRes.getLong(ModelKeys.SINCE);

                    boolean authed = v.containsKey(ModelKeys.TOKEN);
                    if (!authed) {
                        String saltedQuerySecret = AccountUtils.saltedSecret(querySecret, salt);
                        String secret = queryRes.getString(ModelKeys.SECRET);
                        authed = Utils.isNotEmptyString(saltedQuerySecret) && saltedQuerySecret.equals(secret);
                    }
                    if (authed) {
                        byte[] tokenRaw = AccountUtils.getToken();
                        int loginCount = queryRes.getInteger(ModelKeys.LOGIN_COUNT);
                        String token = Utils.encodeBase64(tokenRaw);
                        queryRes.put(ModelKeys.TOKEN, token);
                        queryRes.put(ModelKeys.LAST_LOGIN, System.currentTimeMillis());
                        queryRes.put(ModelKeys.LOGIN_COUNT, loginCount + 1);
                        queryRes.put(ModelKeys.LAST_IP, request.getLastIp());
                        queryRes.put(ModelKeys.APP_LANGUAGE, request.getAppLanguage());
                        queryRes.put(ModelKeys.APP_VERSION, request.getAppVersion());
                        queryRes.put(ModelKeys.DEVICE_TYPE, request.getDeviceType());
                        queryRes.put(ModelKeys.MCC, request.getMcc());
                        queryRes.put(ModelKeys.OS, request.getOs());
                        queryRes.put(ModelKeys.OS_LOCALE, request.getOsLocale());
                        queryRes.put(ModelKeys.TIMEZONE, request.getTimezone());

                        return dbOperator.updateUserInfoDB(queryRes);
                    } else {
                        return Future.failedFuture(DB.StatusCode.STATUS_INVALID_SECRET.toString());
                    }
                });
            })
            .compose(v -> {
                // update in redis
                Future<JsonObject> updateRedisFuture = dbOperator.updateUserInfoRedis(v);
                updateRedisFuture.setHandler(responseFuture.completer());
            }, responseFuture);

        entry.complete(new JsonObject());
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

        Future<JsonObject> redisFuture = dbOperator.queryUserInfoRedis(usn);
        redisFuture.setHandler(res -> {
            if (res.failed()) {
                response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_USER_NOT_FOUND));
            } else {
                JsonObject jsonObject = res.result();
                if (!token.equals(jsonObject.getString(ModelKeys.TOKEN))) {
                    response.complete(DbProtoUtils.makeResult(DB.StatusCode.STATUS_INVALID_TOKEN));
                } else {
                    Future<JsonObject> delFuture = dbOperator.removeUserInfoRedis(usn);
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
    public void userExtraInfoQuery(DB.UserEntry request, Future<DB.UserOpResult> response) {
        LOGGER.info("userExtra query:" + request);

        long usn = request.getUsn();
        long uid = request.getUid();
        String email = request.getEmail();

        if (usn == 0L && (uid == 0L || Utils.isEmptyString(email))) {
            response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_INVALID_PARAM));
            return;
        }

        Future<JsonObject> redisQueryFuture = dbOperator.queryUserInfoRedis(usn);

        redisQueryFuture.setHandler(redisQueryRes -> {
            if (redisQueryRes.succeeded()) {
                response.complete(DbProtoUtils.makeUserOpOkResult(redisQueryRes.result()));
            } else {
                response.complete(DbProtoUtils.makeUserOpResult(DB.StatusCode.STATUS_USER_NOT_FOUND));
            }
        });
    }
}
