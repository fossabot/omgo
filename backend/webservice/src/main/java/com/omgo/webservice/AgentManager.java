package com.omgo.webservice;

import com.omgo.webservice.service.Services;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;

import java.util.*;

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
            Map<String, String> agents = Services.getInstance().getAllValues(Services.generatePath(root, agentServiceType));
            agentSet.addAll(agents.keySet());
        }

        // TODO: 13/10/2017 create a timer here to update agent services by calling Services.getAllValues
        return new ArrayList<>(agentSet);
    }
}
