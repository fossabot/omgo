package com.omgo.dbservice.etcd;

import com.coreos.jetcd.Client;
import com.coreos.jetcd.KV;
import com.coreos.jetcd.Watch;
import com.coreos.jetcd.data.ByteSequence;
import com.coreos.jetcd.data.KeyValue;
import com.coreos.jetcd.kv.GetResponse;
import com.coreos.jetcd.kv.PutResponse;
import com.coreos.jetcd.options.GetOption;
import com.coreos.jetcd.options.WatchOption;
import com.coreos.jetcd.watch.WatchEvent;
import com.coreos.jetcd.watch.WatchResponse;
import com.omgo.dbservice.MainVerticle;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.grpc.VertxChannelBuilder;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Service discovery util
 */
public class Services {

    // etcd watcher events
    public static final String EVENT_SERVICE_ADD = "service.add";
    public static final String EVENT_SERVICE_REMOVE = "service.remove";

    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class);

    private static Services instance = null;

    public static Services getInstance() {
        if (instance == null) {
            synchronized (Services.class) {
                if (instance == null) {
                    instance = new Services();
                }
            }
        }
        return instance;
    }

    protected Services() {
    }

    // etcd client
    private Client client;

    // service pool
    private ServicePool servicePool;

    /**
     * init etcd client with single endpoint
     *
     * @param endpoint
     */
    public void init(String endpoint) {
        if (client != null) {
            client.close();
        }
        client = Client.builder().endpoints(endpoint).build();
    }

    /**
     * init etcd client with one or more endpoints
     *
     * @param endpoints
     */
    public void init(List<String> endpoints) {
        if (client != null) {
            client.close();
        }
        client = Client.builder().endpoints(endpoints).build();
    }

    /**
     * generate a range key for etcd get/put operation
     *
     * @param key
     * @return
     */
    public static ByteSequence getRangeKey(String key) {
        byte[] keyBytes = key.getBytes();
        byte[] endKeyBytes = Arrays.copyOf(keyBytes, keyBytes.length);
        endKeyBytes[endKeyBytes.length - 1]++;

        return ByteSequence.fromBytes(endKeyBytes);
    }

    public KV getKVClient() {
        if (client != null) {
            return client.getKVClient();
        }

        return null;
    }

    public Watch getWatchClient() {
        if (client != null) {
            client.getWatchClient();
        }
        return null;
    }

    public static String getDir(String path) {
        Path fullPath = Paths.get(path);
        return fullPath.getParent().toString();
    }

    public ServicePool getServicePool() {
        return servicePool;
    }

    public ServicePool createServicePool(Vertx vertx, String root, List<String> services) {
        if (servicePool != null) {
            LOGGER.warn("service pool already exist");
            return servicePool;
        }

        servicePool = ServicePool.create(vertx, root, services);
        return servicePool;
    }


    // Private classes

    private static class ServiceClient {
        String key;
        ManagedChannel channel;
    }

    private static class Service {
        List<ServiceClient> clients;
        AtomicInteger idx;

        protected Service() {
            clients = new ArrayList<>();
            idx = new AtomicInteger(0);
        }
    }

    public static class ServicePool {
        Vertx vertx;
        String root;
        Set<String> names;
        Map<String, Service> services;
        boolean namesProvided;

        public ServicePool() {
            root = "";
            names = new HashSet<>();
            services = new HashMap<>();
            namesProvided = false;
        }

        public static ServicePool create(Vertx vertx, String root, List<String> services) {
            ServicePool pool = new ServicePool();
            pool.vertx = vertx;
            pool.root = root;
            if (services != null && services.size() > 0) {
                pool.namesProvided = true;
            }

            for (String serviceName : services) {
                String name = root + "/" + serviceName;
                pool.names.add(name);
            }
            LOGGER.info("all service names:", pool.names);

            pool.connectAll(root);

            return pool;
        }

        public void connectAll(String directory) {
            KV client = Services.getInstance().getKVClient();
            if (client != null) {
                try {
                    ByteSequence endKey = Services.getRangeKey(directory);

                    ByteSequence key = ByteSequence.fromString(directory);

                    CompletableFuture<GetResponse> getFuture = client.get(key, GetOption.newBuilder().withRange(endKey).build());
                    GetResponse response = getFuture.get();
                    List<KeyValue> results = response.getKvs();
                    for (KeyValue kv : results) {
                        String snHost = kv.getValue().toStringUtf8();
                        String snKey = kv.getKey().toStringUtf8();
                        addService(snKey, snHost);
                    }
                    LOGGER.info("services added");
                    startWatcher();
                } catch (Exception e) {
                    e.printStackTrace();
                }
                client.close();
            }
        }

        private void startWatcher() {
            LOGGER.info("start watching");
            vertx.executeBlocking(future -> {
                watcher();
                future.complete();
            }, res -> {
                LOGGER.info("watch complete");
            });
        }

        public void addService(String servicePath, String address) {
            LOGGER.info("adding " + servicePath + " @ " + address);

            if (address == null || address.equals("")) {
                LOGGER.error("invalid service address");
                return;
            }

            String[] comps = address.split(":");
            if (comps.length < 2) {
                LOGGER.error("invalid service address");
                return;
            }

            String host = comps[0];
            int port = Integer.parseInt(comps[1]);

            // name check
            String serviceKind = getDir(servicePath);
            if (namesProvided && !names.contains(serviceKind)) {
                return;
            }

            // create try new service kind init
            if (!services.containsKey(serviceKind)) {
                Service service = new Service();
                services.put(serviceKind, service);
            }

            // create service connections
            Service service = services.get(serviceKind);
            ManagedChannel channel = VertxChannelBuilder
                .forAddress(vertx, host, port)
                .usePlaintext(true)
                .build();
            ServiceClient serviceClient = new ServiceClient();
            serviceClient.channel = channel;
            serviceClient.key = servicePath;
            service.clients.add(serviceClient);

            EventBus eb = vertx.eventBus();
            eb.publish(EVENT_SERVICE_ADD, servicePath);
        }

        public void watcher() {
            Watch watch = Services.getInstance().getWatchClient();
            if (watch != null) {
                ByteSequence key = ByteSequence.fromString("backends/");
                ByteSequence endKey = Services.getRangeKey("backends/");
                Watch.Watcher watcher = watch.watch(key, WatchOption.newBuilder().withRange(endKey).build());

                try {
                    WatchResponse response = watcher.listen();
                    for (WatchEvent event : response.getEvents()) {
                        LOGGER.info("type={}, key={}, value={}",
                            event.getEventType(),
                            Optional.ofNullable(event.getKeyValue().getKey())
                                .map(ByteSequence::toStringUtf8)
                                .orElse(""),
                            Optional.ofNullable(event.getKeyValue().getValue())
                                .map(ByteSequence::toStringUtf8)
                                .orElse(""));

                        KeyValue kv = event.getKeyValue();
                        if (kv != null) {
                            if (kv.getValue() == null || kv.getKey() == null) {
                                continue;
                            }

                            EventBus eb = vertx.eventBus();
                            if (event.getEventType() == WatchEvent.EventType.PUT) {
                                eb.publish(EVENT_SERVICE_ADD, kv.getKey());
                            } else if (event.getEventType() == WatchEvent.EventType.DELETE) {
                                eb.publish(EVENT_SERVICE_REMOVE, kv.getKey());
                            }
                        }
                    }
                } catch (InterruptedException e) {
                    e.printStackTrace();
                    LOGGER.error(e);
                }

                LOGGER.info("closing watcher");
                watcher.close();
            }
        }

        public void removeService(String fullPath) {
            String serviceKind = getDir(fullPath);
            if (namesProvided && !names.contains(serviceKind)) {
                return;
            }

            if (!services.containsKey(serviceKind)) {
                return;
            }

            Service service = services.get(serviceKind);
            List<ServiceClient> toRemove = new ArrayList<>();
            for (ServiceClient serviceClient : service.clients) {
                if (serviceClient.key.equals(fullPath) && serviceClient.channel != null) {
                    serviceClient.channel.shutdown();
                    toRemove.add(serviceClient);
                    LOGGER.info("service removed:", fullPath);
                }
            }
            if (!toRemove.isEmpty()) {
                service.clients.removeAll(toRemove);
            }
        }

        public ManagedChannel getChannelWithId(String fullPath) {
            String serviceKind = getDir(fullPath);
            if (services.containsKey(serviceKind)) {
                Service service = services.get(serviceKind);
                for (ServiceClient client : service.clients) {
                    if (client.key.equals(fullPath)) {
                        return client.channel;
                    }
                }
            }

            return null;
        }

        public ManagedChannel getChannel(String serviceKind) {
            if (services.containsKey(serviceKind)) {
                Service service = services.get(serviceKind);
                if (service.clients.size() == 0) {
                    return null;
                }
                int idx = service.idx.addAndGet(1) % service.clients.size();
                return service.clients.get(idx).channel;
            }
            return null;
        }

        public void registerService(String fullPath, String address) {
            KV kvClient = getInstance().getKVClient();
            if (kvClient == null) {
                return;
            }

            ByteSequence key = ByteSequence.fromString(fullPath);
            ByteSequence value = ByteSequence.fromString(address);
            CompletableFuture<PutResponse> putFuture = kvClient.put(key, value);
            try {
                PutResponse putResponse = putFuture.get();
                LOGGER.info(String.format("service %s @ %s added", fullPath, address));
            } catch (InterruptedException e) {
                e.printStackTrace();
            } catch (ExecutionException e) {
                e.printStackTrace();
            }
            kvClient.close();
        }
    }
}
