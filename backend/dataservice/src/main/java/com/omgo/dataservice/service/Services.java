package com.omgo.dataservice.service;

import com.coreos.jetcd.KV;
import com.coreos.jetcd.data.ByteSequence;
import com.coreos.jetcd.data.KeyValue;
import com.coreos.jetcd.kv.GetResponse;
import com.coreos.jetcd.kv.PutResponse;
import com.coreos.jetcd.options.GetOption;
import com.omgo.dataservice.model.Utils;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
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

public class Services {
    private static final Logger LOGGER = LoggerFactory.getLogger(Services.class);

    private static final long DEFAULT_WATCH_INTERVAL = 5000;

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

    private Services() {
    }

    // service client
    private com.coreos.jetcd.Client client;

    /**
     * init service client with single endpoint
     *
     * @param endpoint service endpoint
     */
    public void init(String endpoint) {
        if (client != null) {
            client.close();
        }
        client = com.coreos.jetcd.Client.builder().endpoints(endpoint).build();
    }

    /**
     * init service client with one or more endpoints
     *
     * @param endpoints service endpoint list
     */
    public void init(List<String> endpoints) {
        if (client != null) {
            client.close();
        }
        client = com.coreos.jetcd.Client.builder().endpoints(endpoints).build();
    }

    /**
     * generate a range key for service get/put operation
     *
     * @param key origin key
     * @return ranged key
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
     * @param args path parameters
     * @return full service path
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
     * @param port service port
     * @return service local address in ip:port format
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
     * get service key-value client
     *
     * @return KV client instance
     */
    public KV getKVClient() {
        if (client != null) {
            return client.getKVClient();
        }

        return null;
    }

    /**
     * get dir part from a path
     * getDir(root/service/name) = root/service
     *
     * @param path key
     * @return key had its last path component removed
     */
    public static String getDir(String path) {
        Path fullPath = Paths.get(path);
        return fullPath.getParent().toString();
    }

    /**
     * get filename part from a path
     * getName(root/service/name) = name
     *
     * @param path
     * @return
     */
    public static String getName(String path) {
        Path fullPath = Paths.get(path);
        return fullPath.getFileName().toString();
    }

    /**
     * get all values under path
     *
     * @param path
     * @return
     */
    public Map<String, String> getAllValues(String path) {
        Map<String, String> keyValues = new HashMap<>();
        KV client = getKVClient();
        if (client != null) {
            try {
                ByteSequence endKey = getRangeKey(path);
                ByteSequence key = ByteSequence.fromString(path);

                CompletableFuture<GetResponse> getFuture = client.get(key, GetOption.newBuilder().withRange(endKey).build());
                GetResponse response = getFuture.get();
                List<KeyValue> results = response.getKvs();
                for (KeyValue kv : results) {
                    String snKey = kv.getKey().toStringUtf8();
                    String snAddress = kv.getValue().toStringUtf8();
                    keyValues.put(snKey, snAddress);
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
            client.close();
        }
        return keyValues;
    }

    /**
     * register service to ETCD
     *
     * @param fullPath full path of the service e.g. 'backends/agent/agent-asia-01'
     * @param address  service address and port e.g. '192.168.0.1:8888'
     */
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
            LOGGER.error("error while setRoute service for interrupt");
        } catch (ExecutionException e) {
            e.printStackTrace();
            LOGGER.error("error while setRoute service for exception");
        }
        kvClient.close();
    }

    // service-type to service-pool map
    private Map<String, Pool> servicePoolMap = new HashMap<>();

    public Pool getServicePool(Vertx vertx, String root, String type) {
        if (servicePoolMap.containsKey(type)) {
            return servicePoolMap.get(type);
        }

        Pool pool = new Pool(vertx, root, type);
        pool.updateServices();
        pool.startWatch(DEFAULT_WATCH_INTERVAL);

        servicePoolMap.put(type, pool);

        return pool;
    }

    /**
     * Service client
     */
    private static class Client {
        // service name
        String name;

        // service full path
        String fullPath;

        // service address:port
        String address;

        // managed channel for creating grpc stub
        ManagedChannel channel;

        private Client() {
        }

        protected static Client create(Vertx vertx, String fullPath, String address) {
            if (vertx == null || Utils.isEmptyString(fullPath) || Utils.isEmptyString(address)) {
                return null;
            }

            String[] comps = address.split(":");
            if (comps.length < 2) {
                return null;
            }

            String host = comps[0];
            int port = Integer.parseInt(comps[1]);

            try {
                Client client = new Client();
                client.fullPath = fullPath;
                client.address = address;
                client.name = getName(fullPath);

                client.channel = VertxChannelBuilder
                    .forAddress(vertx, host, port)
                    .usePlaintext(true)
                    .build();

                return client;
            } catch (Exception e) {
                LOGGER.error(e);
            }

            return null;
        }
    }

    public static class Pool {
        public interface OnChangeListener {
            void onServiceAdded(Pool pool);
            void onServiceRemoved(Pool pool);
        }

        // vertx instance
        Vertx vertx;

        // service root
        String root;

        // service type
        String type;

        // service change listener
        List<OnChangeListener> onChangeListeners;

        // fullpath to serviceClient map
        Map<String, Client> serviceMap;

        // serviceClient array
        List<Client> serviceArray;

        // atomic index for round-robin
        AtomicInteger idx;

        // watcher timer id
        long timerId;

        protected Pool(Vertx vertx, String root, String type) {
            this.vertx = vertx;
            this.root = root;
            this.type = type;
            this.idx = new AtomicInteger(0);
            this.serviceMap = new HashMap<>();
            this.serviceArray = new ArrayList<>();
        }

        protected void startWatch(long period) {
            if (timerId != 0L) {
                vertx.cancelTimer(timerId);
            }
            timerId = vertx.setPeriodic(period, id -> {
                updateServices();
            });
        }

        public void stopWatch() {
            if (timerId != 0L) {
                vertx.cancelTimer(timerId);
            }
        }

        public void addOnChangeListener(OnChangeListener listener) {
            if (this.onChangeListeners == null) {
                this.onChangeListeners = new ArrayList<>();
            }
            this.onChangeListeners.add(listener);
        }

        public void removeOnChangeListener(OnChangeListener listener) {
            if (this.onChangeListeners != null) {
                this.onChangeListeners.remove(listener);
            }
        }

        public void clearOnChangeListener() {
            if (this.onChangeListeners != null) {
                this.onChangeListeners.clear();
            }
        }

        public void clear() {
            stopWatch();
            serviceMap.clear();
            for (Client client : serviceArray) {
                if (client.channel != null && !client.channel.isShutdown()) {
                    client.channel.shutdown();
                }
            }
            serviceArray.clear();
            clearOnChangeListener();
        }

        protected void updateServices() {
            Map<String, String> onlineKeyValues = getInstance().getAllValues(generatePath(root, type));
            // check for service offline
            List<String> serviceToRemoveName = new ArrayList<>();
            List<Client> serviceToRemoveClient = new ArrayList<>();
            serviceMap.forEach((fullPath, serviceClient) -> {
                if (!onlineKeyValues.containsKey(fullPath)) {
                    serviceToRemoveName.add(fullPath);
                    serviceToRemoveClient.add(serviceClient);
                }
            });

            final boolean[] serviceAddRemove = {false, false};
            // check for online services
            onlineKeyValues.forEach((fullPath, address) -> {
                if (!serviceMap.containsKey(fullPath)) {
                    // new service online
                    Client newClient = Client.create(vertx, fullPath, address);
                    if (newClient != null) {
                        serviceMap.put(fullPath, newClient);
                        serviceArray.add(newClient);
                        serviceAddRemove[0] = true;
                    }
                } else {
                    Client oldClient = serviceMap.get(fullPath);
                    if (!oldClient.address.equals(address)) {
                        // service change its address, remove it
                        // and it will be re-add by watch process
                        serviceToRemoveName.add(fullPath);
                        serviceToRemoveClient.add(oldClient);
                    }
                }
            });

            for (String path : serviceToRemoveName) {
                serviceMap.remove(path);
            }
            if (!serviceToRemoveClient.isEmpty()) {
                for (Client serviceClient : serviceToRemoveClient) {
                    if (serviceClient.channel != null && !serviceClient.channel.isShutdown()) {
                        serviceClient.channel.shutdown();
                    }
                }
                serviceArray.removeAll(serviceToRemoveClient);
                serviceAddRemove[1] = true;
            }

            // notify
            if (serviceAddRemove[0]) {
                dispatchOnServiceAdded();
            }
            if (serviceAddRemove[1]) {
                dispatchOnServiceRemoved();
            }
        }

        /**
         * get a gRPC channel from this service pool in Round-Robin style
         *
         * @return ManagedChannel
         */
        public ManagedChannel getClient() {
            if (!serviceArray.isEmpty()) {
                int index = idx.addAndGet(1) % serviceArray.size();
                return serviceArray.get(index).channel;
            }
            return null;
        }

        private void dispatchOnServiceAdded() {
            dispatchOnServiceListeners(true);
        }

        private void dispatchOnServiceRemoved() {
            dispatchOnServiceListeners(false);
        }

        private void dispatchOnServiceListeners(boolean isAdd) {
            if (onChangeListeners != null) {
                for (int i = 0, z = onChangeListeners.size(); i < z; i++) {
                    OnChangeListener listener = onChangeListeners.get(i);
                    if (listener != null) {
                        if (isAdd) {
                            listener.onServiceAdded(this);
                        } else {
                            listener.onServiceRemoved(this);
                        }
                    }
                }
            }
        }
    }
}
