package com.omgo.dataservice;

import com.omgo.dataservice.model.ModelConverter;
import com.omgo.utils.ModelKeys;
import com.omgo.utils.Utils;
import io.vertx.core.AsyncResult;
import io.vertx.core.Future;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.mongo.FindOptions;
import io.vertx.ext.mongo.MongoClient;
import io.vertx.ext.mongo.UpdateOptions;
import io.vertx.redis.RedisClient;
import io.vertx.redis.RedisTransaction;
import proto.Db;

import java.util.List;
import java.util.Map;
import java.util.Set;

public class DbOperator {

    // redis config
    private static final int REDIS_USER_EXPIRE_DURATION = 24 * 60 * 60;

    // clients
    private MongoClient mongoClient;
    private RedisClient redisClient;

    public DbOperator(MongoClient mongoClient, RedisClient redisClient) {
        this.mongoClient = mongoClient;
        this.redisClient = redisClient;
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
            redisClient.get(AccountUtils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    String strJson = res.result();
                    JsonObject result = Utils.safeParseJson(strJson);
                    if (result != null) {
                        future.complete(result);
                    } else {
                        future.fail(String.format("invalid value in redis (usn:%d) value:%s", usn, strJson));
                    }
                } else {
                    future.fail(res.cause());
                }
            });
        }

        return future;
    }

    private void mongodbFindUser(JsonObject queryObject, Future<JsonObject> future) {
        mongoClient.find(ModelKeys.USER, queryObject, res -> {
            if (res.succeeded()) {
                List<JsonObject> results = res.result();
                if (results != null && results.size() > 0) {
                    future.complete(results.get(0));
                } else {
                    future.fail("query success with no results");
                }
            } else {
                future.fail(res.cause());
            }
        });
    }

    public Future<JsonObject> mongodbGenerateUsn() {
        JsonObject queryObject = new JsonObject();
        queryObject.put("_id", ModelKeys.USER);

        JsonObject updateObject = new JsonObject();
        updateObject.put(ModelKeys.USN, AccountUtils.nextUsnIncrement());
        updateObject.put(ModelKeys.UID, AccountUtils.nextUidIncrement());
        JsonObject incObject = new JsonObject();
        incObject.put("$inc", updateObject);

        UpdateOptions updateOptions = new UpdateOptions();
        updateOptions.setMulti(false);
        updateOptions.setUpsert(false);
        updateOptions.setReturningNewDocument(true);

        Future<JsonObject> future = Future.future();

        mongoClient.findOneAndUpdateWithOptions("status", queryObject, incObject, new FindOptions(), updateOptions, res -> {
            if (res.succeeded()) {
                future.complete(res.result());
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }

    /**
     * Query user info in Mongodb
     *
     * @param json
     * @return Future
     */
    public Future<JsonObject> queryUserInfoDB(JsonObject json) {
        Future<JsonObject> future = Future.future();

        Db.DB.UserEntry userEntry = ModelConverter.json2UserEntry(json);

        long usn = userEntry.getUsn();
        long uid = userEntry.getUid();
        String email = userEntry.getEmail();

        if (usn == 0L && uid == 0L && Utils.isEmptyString(email)) {
            future.fail("invalid query key");
        } else {
            String query = "";
            JsonObject queryObject = new JsonObject();
            if (usn != 0L) {
                queryObject.put(ModelKeys.USN, usn);
            } else if (uid != 0L) {
                queryObject.put(ModelKeys.UID, uid);
            } else if (Utils.isNotEmptyString(email)) {
                queryObject.put(ModelKeys.EMAIL, email);
            }
            mongodbFindUser(queryObject, future);
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
            long usn = userInfoJson.getLong(ModelKeys.USN);
            String key = AccountUtils.getRedisKey(usn);
            RedisTransaction transaction = redisClient.transaction();
            transaction.multi(event -> {
                transaction.set(key, userInfoJson.encode(), setEvent -> {
//                transaction.hmset(key, userInfoJson, setEvent -> {
                    if (isRedisTransactionSucceed(setEvent)) {
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
                            future.fail(setEvent.cause());
                        });
                    }
                });
            });
        }

        return future;
    }

    /**
     * Update user info in Mongodb
     *
     * @param userJson
     * @return Future
     */
    public Future<JsonObject> updateUserInfoDB(JsonObject userJson) {
        Future<JsonObject> future = Future.future();

        long usn = userJson.getLong(ModelKeys.USN);
        if (usn == 0L) {
            future.fail("invalid usn");
            return future;
        }

        Set<String> updatableKeys = ModelConverter.getUserUpdatableMapKeySet();
        JsonObject userObject = new JsonObject();

        Map<String, Object> map = userJson.getMap();
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            String key = entry.getKey();
            if (!updatableKeys.contains(key)) {
                continue;
            }
            userObject.put(key, entry.getValue());
        }

        if (userObject.isEmpty()) {
            future.fail("update user info failed, invalid user info");
            return future;
        }

        JsonObject queryObject = new JsonObject();
        queryObject.put(ModelKeys.USN, usn);

        JsonObject updateObject = new JsonObject();
        updateObject.put("$set", userObject);

        // update
        mongoClient.findOneAndUpdateWithOptions(ModelKeys.USER, queryObject, updateObject, new FindOptions(), new UpdateOptions(true), res -> {
            if (res.succeeded()) {
                mongodbFindUser(queryObject, future);
            } else {
                future.fail(res.cause());
            }
        });

        return future;
    }

    /**
     * Insert userInfo to Mongodb
     *
     * @param userJson
     * @return
     */
    public Future<JsonObject> insertUserInfoDB(JsonObject userJson) {
        Future<JsonObject> future = Future.future();

        JsonObject insertObject = userJson.copy();
        insertObject.put("_id", userJson.getLong(ModelKeys.USN));
        mongoClient.insert(ModelKeys.USER, insertObject, res -> {
            if (res.succeeded()) {
                future.complete(userJson);
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
