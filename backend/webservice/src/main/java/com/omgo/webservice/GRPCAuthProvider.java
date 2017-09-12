package com.omgo.webservice;

import io.grpc.ManagedChannel;
import io.vertx.core.AsyncResult;
import io.vertx.core.Future;
import io.vertx.core.Handler;
import io.vertx.core.Vertx;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.auth.AuthProvider;
import io.vertx.ext.auth.User;
import proto.DBServiceGrpc;
import proto.Db;
import proto.common.Common;

public class GRPCAuthProvider implements AuthProvider {

    private static final Logger LOGGER = LoggerFactory.getLogger(GRPCAuthProvider.class);

    private static final String STRING_EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$";

    private DBServiceGrpc.DBServiceVertxStub dbServiceVertxStub;

    public GRPCAuthProvider(Vertx vertx, ManagedChannel channel) {
        dbServiceVertxStub = DBServiceGrpc.newVertxStub(channel);
    }

    @Override
    public void authenticate(JsonObject jsonObject, Handler<AsyncResult<User>> handler) {
        String email = jsonObject.getString("email");
        String password = jsonObject.getString("password");
        String token = jsonObject.getString("token");
        String strUsn = jsonObject.getString("usn");
        long usn = Utils.isEmptyString(strUsn) ? 0L : Long.parseLong(strUsn);

        // TODO: 11/09/2017

        if (Utils.isEmptyString(email) || !email.matches(STRING_EMAIL_REGEX)) {
            handler.handle(Future.failedFuture("auth info invalid email:" + email));
            return;
        }

        if (Utils.isEmptyString(password)) {
            handler.handle(Future.failedFuture("auth info invalid password"));
            return;
        }

        Db.DB.UserExtendInfo.Builder extendInfoBuilder = Db.DB.UserExtendInfo.newBuilder();
        Common.UserInfo.Builder userInfoBuilder = Common.UserInfo.newBuilder();

        userInfoBuilder.setEmail(email);
        extendInfoBuilder.setInfo(userInfoBuilder.build())
            .setSecret(password);

        dbServiceVertxStub.userLogin(extendInfoBuilder.build(), res -> {
            if (res.failed()) {
                handler.handle(Future.failedFuture(res.cause()));
            } else {
                Db.DB.UserOpResult loginResult = res.result();
                Db.DB.StatusCode status = loginResult.getResult().getStatus();
                // RPC invoked, check result
                if (status != Db.DB.StatusCode.STATUS_OK) {
                    handler.handle(Future.failedFuture(loginResult.getResult().getMsg()));
                    return;
                }
                // check user
                Db.DB.UserExtendInfo userExtendInfo = loginResult.getUserExtInfo();
                if (userExtendInfo == null) {
                    handler.handle(Future.failedFuture("invalid user extend info"));
                    return;
                }

                LOGGER.info(loginResult);
                handler.handle(Future.succeededFuture(new GrpcAuthUser(this, email, userExtendInfo)));
            }
        });
    }
}
