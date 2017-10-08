package com.omgo.webservice;

import com.omgo.webservice.etcd.Services;
import io.vertx.core.Vertx;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;

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

    private static final Logger LOGGER = LoggerFactory.getLogger(AgentManager.class);

    public void startWatch() {
        EventBus eb = vertx.eventBus();
        eb.consumer(Services.EVENT_SERVICE_ADD, res -> {
            LOGGER.info("service added:");
            LOGGER.info(res.body());
        });

        eb.consumer(Services.EVENT_SERVICE_REMOVE, res -> {
            LOGGER.info("service removed:");
            LOGGER.info(res.body());
        });
    }
}
