package com.omgo.dbservice;

import io.vertx.core.AbstractVerticle;
import io.vertx.grpc.VertxServer;
import io.vertx.grpc.VertxServerBuilder;

import java.io.IOException;

public class MainVerticle extends AbstractVerticle {

    @Override
    public void start() {
        VertxServer rpcServer = VertxServerBuilder
            .forAddress(vertx, "localhost", 8080)
            .addService(new DBServiceGrpcImpl())
            .build();

        // Start is asynchronous
        try {
            rpcServer.start();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

}
