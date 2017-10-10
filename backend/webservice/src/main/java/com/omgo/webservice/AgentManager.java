package com.omgo.webservice;

import com.omgo.webservice.etcd.Services;
import io.vertx.core.Vertx;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class AgentManager {
    private static AgentManager instance;

    private AgentManager() {
    }

    public static AgentManager getInstance() {
        if (instance == null) {
            synchronized (AgentManager.class) {
                if (instance == null) {
                    instance = new AgentManager();
                }
            }
        }
        return instance;
    }

    private static final Logger LOGGER = LoggerFactory.getLogger(AgentManager.class);

    private Set<String> agentSet = new HashSet<>();

    public List<String> getAgentList(String root, String agentServiceType) {
        if (agentSet.isEmpty()) {
            List<String> agents = Services.getInstance().getAllValues(Services.generatePath(root, agentServiceType));
            agentSet.addAll(agents);
        }
        return new ArrayList<>(agentSet);
    }

    public void startWatch(Vertx vertx) {
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
