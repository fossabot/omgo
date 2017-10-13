package com.omgo.webservice;

import com.omgo.webservice.model.ModelConverter;
import com.omgo.webservice.service.Services;
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

public class GRPCAuthProvider implements AuthProvider, Services.Pool.OnChangeListener {

    private static final Logger LOGGER = LoggerFactory.getLogger(GRPCAuthProvider.class);

    private static final String STRING_EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$";

    private DBServiceGrpc.DBServiceVertxStub dbServiceVertxStub;
    private Services.Pool dataServicePool;
    private ManagedChannel channel;

    public GRPCAuthProvider(Vertx vertx, Services.Pool pool) {
        this.dataServicePool = pool;
        init();
    }

    private void init() {
        channel = dataServicePool.getClient();
        if (channel != null) {
            dbServiceVertxStub = DBServiceGrpc.newVertxStub(channel);
        }
        dataServicePool.addOnChangeListener(this);
    }

    @Override
    public void authenticate(JsonObject jsonObject, Handler<AsyncResult<User>> handler) {
        String email = jsonObject.getString("email", "");
        String password = jsonObject.getString("password", "");
        String token = jsonObject.getString("token", "");
        String strUsn = jsonObject.getString("usn");
        long usn = Utils.isEmptyString(strUsn) ? 0L : Long.parseLong(strUsn);
        String clientIpAddress = jsonObject.getString(ModelConverter.KEY_LAST_IP, "");

        // TODO: 11/09/2017 this regex is kinda invalid for xxx@xxx

        if (Utils.isEmptyString(email) || !email.matches(STRING_EMAIL_REGEX)) {
            handler.handle(Future.failedFuture("auth info invalid email:" + email));
            return;
        }

        if (Utils.isEmptyString(password)) {
            handler.handle(Future.failedFuture("auth info invalid password"));
            return;
        }

        Db.DB.UserEntry.Builder entryBuilder = Db.DB.UserEntry.newBuilder();
        entryBuilder
            .setEmail(email)
            .setSecret(password)
            .setUsn(usn)
            .setLastIp(clientIpAddress)
            .setToken(token);

        if (dbServiceVertxStub == null) {
            handler.handle(Future.failedFuture("dataservice not ready yet"));
            return;
        }

        dbServiceVertxStub.userLogin(entryBuilder.build(), res -> {
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
                Db.DB.UserEntry userEntry = loginResult.getUser();
                if (userEntry == null) {
                    handler.handle(Future.failedFuture("invalid user extend info"));
                    return;
                }

                LOGGER.info(loginResult);
                handler.handle(Future.succeededFuture(new GrpcAuthUser(this, email, userEntry)));
            }
        });
    }

    @Override
    public void onServiceAdded(Services.Pool pool) {
        if (channel == null) {
            LOGGER.info("dataservice online, init...");
            init();
        }
    }

    @Override
    public void onServiceRemoved(Services.Pool pool) {
        if (channel != null && channel.isShutdown()) {
            LOGGER.info("dataservice offline, halt...");
            channel = null;
            dbServiceVertxStub = null;
        }
    }
}
