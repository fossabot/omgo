package com.omgo.dbservice;

import com.coreos.jetcd.Client;
import com.coreos.jetcd.KV;

public final class EtcdUtils {
    private static Client client;

    public static void init(String url) {
        client = Client.builder().endpoints(url).build();
    }

    public static KV getKVClient() {
        if (client != null) {
            return client.getKVClient();
        } else {
            return null;
        }
    }
}
