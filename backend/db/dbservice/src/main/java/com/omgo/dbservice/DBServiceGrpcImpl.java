package com.omgo.dbservice;

import com.omgo.dbservice.driver.Utils;
import com.omgo.dbservice.model.ModelConverter;
import io.vertx.core.Future;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.sql.SQLClient;
import io.vertx.ext.sql.SQLConnection;
import io.vertx.redis.RedisClient;
import proto.DBServiceGrpc;
import proto.Db;
import proto.common.Common;

/**
 * Database gRPC service implementation
 * <p>
 * Created by mg on 17/07/2017.
 */
public class DBServiceGrpcImpl extends DBServiceGrpc.DBServiceVertxImplBase {

    private static final Logger LOGGER = LoggerFactory.getLogger(DBServiceGrpcImpl.class);

    private SQLClient sqlClient;
    private RedisClient redisClient;

    public DBServiceGrpcImpl(SQLClient sqlClient, RedisClient redisClient) {
        this.sqlClient = sqlClient;
        this.redisClient = redisClient;
    }

    @Override
    public void userQuery(Db.DB.UserKey request, Future<Db.DB.UserQueryResponse> response) {
        super.userQuery(request, response);

        long usn = request.getUsn();

    }

    @Override
    public void userUpdateInfo(Common.UserInfo request, Future<Common.RspHeader> response) {
        super.userUpdateInfo(request, response);
    }

    @Override
    public void userRegister(Db.DB.UserRegisterRequest request, Future<Db.DB.UserRegisterResponse> response) {
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
        sqlClient.getConnection(connRes -> {
            if (connRes.succeeded()) {
                SQLConnection connection = connRes.result();
                
            } else {
                future.fail(connRes.cause());
            }
        });

        return future;
    }
}
