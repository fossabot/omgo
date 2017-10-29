package com.omgo.webservice.handler;

import com.omgo.webservice.model.HttpStatus;
import com.omgo.webservice.service.Services;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.ext.web.RoutingContext;
import proto.DBServiceGrpc;

public class BaseGrpcHandler extends BaseHandler implements Services.Pool.OnChangeListener {

    protected DBServiceGrpc.DBServiceVertxStub dbServiceVertxStub;
    protected Services.Pool dataServicePool;
    protected ManagedChannel channel;

    public BaseGrpcHandler(Vertx vertx, Services.Pool servicePool) {
        super(vertx);
        this.dataServicePool = servicePool;
        initServicePool();
    }

    protected void initServicePool() {
        channel = dataServicePool.getClient();
        if (channel != null) {
            dbServiceVertxStub = DBServiceGrpc.newVertxStub(channel);
        }
        dataServicePool.addOnChangeListener(this);
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response) {
        if (dbServiceVertxStub == null) {
            LOGGER.info("dataservice not ready yet");
            routingContext.fail(HttpStatus.INTERNAL_SERVER_ERROR.code);
        }
    }

    @Override
    public void onServiceAdded(Services.Pool pool) {
        if (channel == null) {
            LOGGER.info("dataservice online, init...");
            initServicePool();
        }
    }

    @Override
    public void onServiceRemoved(Services.Pool pool) {
        if (channel != null && channel.isShutdown()) {
            LOGGER.info("dataservice offline, try re-init");
            channel = null;
            dbServiceVertxStub = null;
            initServicePool();
        }
    }
}
