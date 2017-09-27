package com.omgo.webservice.etcd;

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
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.eventbus.EventBus;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.grpc.VertxChannelBuilder;

import java.net.DatagramSocket;
import java.net.InetAddress;
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

    // service name
    public static final String SERVICE_SNOWFLAKE = "snowflake";

    private static final Logger LOGGER = LoggerFactory.getLogger(Services.class);

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

    /**
     * generate service full path by concat them with '/'
     *
     * @param args
     * @return
     */
    public static String generatePath(Object... args) {
        List<String> comps = new ArrayList<>();
        for (int i = 0; i < args.length; i++) {
            if (args[i] != null) {
                comps.add(args[i].toString());
            }
        }

        return String.join("/", comps);
    }

    /**
     * get local ip address and concat with port
     *
     * @param port
     * @return
     */
    public static String getLocalAddress(int port) {
        String ip = "";
        try {
            final DatagramSocket socket = new DatagramSocket();
            socket.connect(InetAddress.getByName("8.8.8.8"), 80);
            ip = socket.getLocalAddress().getHostAddress();
            ip += ":" + String.valueOf(port);
        } catch (Exception e) {
            e.printStackTrace();
        }

        return ip;
    }

    /**
     * get etcd key-value client
     *
     * @return
     */
    public KV getKVClient() {
        if (client != null) {
            return client.getKVClient();
        }

        return null;
    }

    /**
     * get etcd watch client
     *
     * @return
     */
    public Watch getWatchClient() {
        if (client != null) {
            client.getWatchClient();
        }
        return null;
    }

    /**
     * get dir part from a path
     * getDir(root/service/name) = root/service
     *
     * @param path
     * @return
     */
    public static String getDir(String path) {
        Path fullPath = Paths.get(path);
        return fullPath.getParent().toString();
    }

    /**
     * get default service pool
     *
     * @return
     */
    public ServicePool getServicePool() {
        return servicePool;
    }

    /**
     * create service pool
     * after creation, service pool will connect to all service under 'root' given by names
     *
     * @param vertx    Vertx instance
     * @param root     service root ('roots', 'backends', etc.)
     * @param services service names ('snowflake', 'agent', 'game', etc.)
     * @return
     */
    public ServicePool createServicePool(Vertx vertx, String root, List<String> services) {
        if (servicePool != null) {
            LOGGER.warn("service pool already exist");
            return servicePool;
        }

        servicePool = ServicePool.newBuilder().setRoot(root).setVertx(vertx).addServices(services).build();
        return servicePool;
    }


    // Private classes

    /**
     * Service client
     */
    private static class ServiceClient {
        // service key, root + '/' + service name, for example 'root/snowflake/sn-0'
        String key;
        // managed channel for creating grpc stub
        ManagedChannel channel;
    }

    /**
     * Service
     */
    private static class Service {
        List<ServiceClient> clients;
        // atomic index for round-robin
        AtomicInteger idx;

        protected Service() {
            clients = new ArrayList<>();
            idx = new AtomicInteger(0);
        }
    }

    /**
     * Service pool
     */
    public static class ServicePool {
        // vertx instance
        Vertx vertx;

        // etcd root ('root', 'backends', etc.)
        String root;

        // service names that will be connect and watched
        // ('root/snowflake', 'backends/agent', etc.)
        Set<String> names;

        // services
        Map<String, Service> services;

        boolean namesProvided;

        public ServicePool() {
            root = "";
            names = new HashSet<>();
            services = new HashMap<>();
            namesProvided = false;
        }

        public String getServicePath(String serviceName) {
            return generatePath(root, serviceName);
        }

        public static Builder newBuilder() {
            return new Builder();
        }

        /**
         * Service pool builder
         */
        public static final class Builder {
            protected Vertx vertx;
            protected String root;
            protected List<String> names = new ArrayList<>();

            public Builder setVertx(Vertx vertx) {
                this.vertx = vertx;
                return this;
            }

            public Builder setRoot(String root) {
                this.root = root;
                return this;
            }

            public Builder addService(String name) {
                this.names.add(name);
                return this;
            }

            public Builder addServices(List<String> nameList) {
                this.names.addAll(nameList);
                return this;
            }

            public ServicePool build() {
                ServicePool pool = new ServicePool();
                pool.vertx = vertx;
                pool.root = root;
                if (names != null && names.size() > 0) {
                    pool.namesProvided = true;

                    for (String serviceName : names) {
                        String name = root + "/" + serviceName;
                        pool.names.add(name);
                    }
                    LOGGER.info("all service names:" + pool.names);
                } else {
                    LOGGER.info("no service name provided");
                }

                pool.connectAll(root);

                return pool;
            }
        }

        /**
         * connect all services under root
         * after adding the services, a watcher will be created to watch the root
         *
         * @param root
         */
        public void connectAll(String root) {
            KV client = Services.getInstance().getKVClient();
            if (client != null) {
                try {
                    ByteSequence endKey = Services.getRangeKey(root);

                    ByteSequence key = ByteSequence.fromString(root);

                    CompletableFuture<GetResponse> getFuture = client.get(key, GetOption.newBuilder().withRange(endKey).build());
                    GetResponse response = getFuture.get();
                    List<KeyValue> results = response.getKvs();
                    for (KeyValue kv : results) {
                        String snHost = kv.getValue().toStringUtf8();
                        String snKey = kv.getKey().toStringUtf8();
                        addService(snKey, snHost);
                    }
                    LOGGER.info("services added");
                    startWatcher(root);
                } catch (Exception e) {
                    e.printStackTrace();
                }
                client.close();
            }
        }

        /**
         * start etcd watcher in vertx blocking style
         *
         * @param root
         */
        private void startWatcher(String root) {
            LOGGER.info("start watching");
            vertx.executeBlocking(future -> {
                watcher(root);
                future.complete();
            }, res -> {
                LOGGER.info("watch complete");
            });
        }

        /**
         * add a service to pool
         *
         * @param servicePath full path of the service ('roots/snowflake/sn-0', 'backends/agent/a-0', etc.)
         * @param address     grpc address(ip:port) of the service ('127.0.0.1:8888', 'localhost:443', etc.)
         */
        public void addService(String servicePath, String address) {
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

            LOGGER.info("adding " + servicePath + " @ " + address);

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

        /**
         * etcd watcher
         *
         * @param root
         */
        public void watcher(String root) {
            Watch watch = Services.getInstance().getWatchClient();
            if (watch != null) {
                ByteSequence key = ByteSequence.fromString(root);
                ByteSequence endKey = Services.getRangeKey(root);
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

        /**
         * remove a kind of service from pool
         *
         * @param fullPath
         */
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
                LOGGER.info(String.format("service %s @ %s registered", fullPath, address));
            } catch (InterruptedException e) {
                e.printStackTrace();
                LOGGER.error("error while initRoute service for interrupt");
            } catch (ExecutionException e) {
                e.printStackTrace();
                LOGGER.error("error while initRoute service for exception");
            }
            kvClient.close();
        }
    }
}
