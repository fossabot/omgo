package com.omgo.dbservice;

import com.omgo.dbservice.driver.MySQLDriver;
import com.omgo.dbservice.driver.Utils;
import io.vertx.core.Future;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
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

    private MySQLDriver mySQLDriver;
    private RedisClient redisClient;

    public DBServiceGrpcImpl(MySQLDriver sqlDriver, RedisClient redisClient) {
        this.mySQLDriver = sqlDriver;
        this.redisClient = redisClient;
    }

    @Override
    public void userQuery(Db.DB.UserKey request, Future<Db.DB.UserQueryResponse> response) {
        super.userQuery(request, response);

        long usn = request.getUsn();
        // Valid usn, find in redis first
        if (usn != 0) {
            redisClient.hgetall(Utils.getRedisKey(usn), res -> {
                if (res.succeeded()) {
                    res.result();

                } else {
                    LOGGER.error(res.cause());
                }
            });
        }
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
}
