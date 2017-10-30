package com.omgo.webservice;

import com.omgo.webservice.service.Services;
import io.vertx.core.Vertx;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class AgentManager {
    private static final long WATCH_INTERVAL = 30 * 1000; // 30 seconds

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

    private long timerId = 0;
    private Map<String, String> agentServices = new HashMap<>();

    public void init(Vertx vertx, String root, String agentServiceType) {
        if (timerId == 0L) {
            update(root, agentServiceType);
            startWatch(vertx, root, agentServiceType);
        }
    }

    public List<String> getHostList() {
        List<String> hosts = new ArrayList<>();
        hosts.addAll(agentServices.values());
        return hosts;
    }

    private void startWatch(Vertx vertx, String root, String type) {
        if (timerId != 0L) {
            return;
        }

        timerId = vertx.setPeriodic(WATCH_INTERVAL, id -> {
            update(root, type);
        });
    }

    private void update(String root, String type) {
        Map<String, String> agents = Services.getInstance().getAllValues(Services.generatePath(root, type));
        agentServices.clear();
        agentServices.putAll(agents);
    }
}
