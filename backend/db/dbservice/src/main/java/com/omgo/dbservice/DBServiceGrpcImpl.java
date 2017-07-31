package com.omgo.dbservice;

import com.omgo.dbservice.model.ModelConverter;
import com.omgo.dbservice.model.Utils;
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
import proto.common.Common;

import java.util.List;

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
    public void userQuery(Db.DB.UserKey request, Future<Db.DB.UserQueryResponse> response) {
        LOGGER.info("userQuery", request);

        Handler<Common.UserInfo> successHandler = new Handler<Common.UserInfo>() {
            @Override
            public void handle(Common.UserInfo userInfo) {
                Db.DB.UserQueryResponse queryResponse = Db.DB.UserQueryResponse.newBuilder()
                    .setResult(ModelConverter.createSuccessRspHeader())
                    .setInfo(userInfo)
                    .build();

                response.complete(queryResponse);
            }
        };

        Future<Common.UserInfo> redisFuture = queryUserInfoRedis(request.getUsn());
        redisFuture.setHandler(res -> {
            if (res.succeeded()) {
                successHandler.handle(res.result());
            } else {
                Future<Common.UserInfo> mysqlFuture = queryUserInfoSQL(request);
                mysqlFuture.setHandler(sqlRes -> {
                    if (sqlRes.succeeded()) {
                        successHandler.handle(sqlRes.result());
                    } else {
                        response.fail("user query failed in both redis and mysql");
                    }
                });
            }
        });
    }

    @Override
    public void userUpdateInfo(Common.UserInfo request, Future<Common.RspHeader> response) {
        super.userUpdateInfo(request, response);
    }

    @Override
    public void userRegister(Db.DB.UserRegisterRequest request, Future<Db.DB.UserRegisterResponse> response) {
        LOGGER.info("userRegister", request);

        request.get
        String email = "";
        if (request.getInfo() != null) {
            email = request.getInfo().getEmail();
        }
        if (!Utils.isValidEmailAddress(email)) {
            response.fail("invalid email address");
        }


        super.userRegister(request, response);
    }

    @Override
    public void userLogin(Db.DB.UserLoginRequest request, Future<Db.DB.UserLoginResponse> response) {
        super.userLogin(request, response);
    }

    @Override
    public void userLogout(Db.DB.UserLogoutRequest request, Future<Common.RspHeader> response) {
        super.userLogout(request, response);
    }

    @Override
    public void userExtraInfoQuery(Db.DB.UserKey request, Future<Db.DB.UserExtraInfo> response) {
        super.userExtraInfoQuery(request, response);
    }

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
                            future.complete(ModelConverter.json2UserInfo(results.get(0)));
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
}
