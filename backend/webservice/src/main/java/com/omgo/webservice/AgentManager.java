package com.omgo.webservice;

import com.omgo.webservice.etcd.Services;
import io.vertx.core.Vertx;
import io.vertx.core.eventbus.EventBus;

public class AgentManager {
    private static AgentManager instance;

    private Vertx vertx;

    private AgentManager(Vertx vertx) {
        this.vertx = vertx;
    }

    public static AgentManager getInstance(Vertx vertx) {
        if (instance == null) {
            synchronized (AgentManager.class) {
                if (instance == null) {
                    instance = new AgentManager(vertx);
                }
            }
        }
        return instance;
    }

    public void startWatch() {
        EventBus eb = vertx.eventBus();
        eb.consumer(Services.EVENT_SERVICE_ADD, res -> {

        });

        eb.consumer(Services.EVENT_SERVICE_REMOVE, res -> {

        });
    }
}
