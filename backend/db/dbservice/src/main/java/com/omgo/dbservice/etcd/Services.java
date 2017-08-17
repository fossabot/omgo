package com.omgo.dbservice.etcd;

import com.coreos.jetcd.Client;
import com.coreos.jetcd.KV;
import com.coreos.jetcd.Watch;
import com.coreos.jetcd.data.ByteSequence;
import com.coreos.jetcd.data.KeyValue;
import com.coreos.jetcd.kv.GetResponse;
import com.coreos.jetcd.options.GetOption;
import com.omgo.dbservice.MainVerticle;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.grpc.VertxChannelBuilder;

import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.concurrent.CompletableFuture;

public class Services {
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

    private Client client;

    public void init(String endpoint) {
        if (client != null) {
            client.close();
        }
        client = Client.builder().endpoints(endpoint).build();
    }

    public void init(List<String> endpoints) {
        if (client != null) {
            client.close();
        }
        client = Client.builder().endpoints(endpoints).build();
    }

    public static ByteSequence getRangeKey(String key) {
        byte[] keyBytes = key.getBytes();
        byte[] endKeyBytes = Arrays.copyOf(keyBytes, keyBytes.length);
        endKeyBytes[endKeyBytes.length -1]++;

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

    private static class Service {
        List<ManagedChannel> channels;
        int idx;

        protected Service() {
            channels = new ArrayList<>();
            idx = 0;
        }
    }

    private static class ServicePool {
        String root;
        Map<String, Boolean> names;
        Map<String, Service> services;
        boolean namesProvided;

        public ServicePool() {
            root = "";
            names = new HashMap<>();
            services = new HashMap<>();
            namesProvided = false;
        }

        public static ServicePool Create(String root, List<String> services) {
            ServicePool pool = new ServicePool();

            pool.root = root;
            if (services != null && services.size() > 0) {
                pool.namesProvided = true;
            }

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
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
        }

        public void addService(String servicePath, String address) {
            LOGGER.info("adding " + servicePath + " @ " + address);

            // name check
            String serviceKind = getDir(servicePath);
            if (namesProvided && (!names.containsKey(serviceKind) || !names.get(serviceKind))) {
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
                .forAddress(Vertx.vertx(), InetSocketAddress)
        }
    }
}
