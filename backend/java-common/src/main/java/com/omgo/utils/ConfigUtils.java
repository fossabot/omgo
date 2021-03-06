package com.omgo.utils;

import io.vertx.core.Vertx;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonObject;

public class ConfigUtils {

    public static final String INFO_VERSION = "info.version";
    public static final String DEBUG = "debug";
    public static final String STANDALONE = "standalone";

    public static final String METRICS_PATH = "metrics.path";

    public static final String SERVICE_ROOT = "service.root";
    public static final String SERVICE_KIND = "service.kind";
    public static final String SERVICE_SELF = "service.self";
    public static final String SERVICE_HOST = "service.host";
    public static final String SERVICE_PORT = "service.port";
    public static final String SERVICE_TYPES = "service.types";

    public static final String SESSION_MAP = "session.map";
    public static final String SESSION_EXPIRE = "session.expire";

    public static final String ETCD_HOST = "etcd.host";

    public static final String REDIS_HOST = "redis.host";
    public static final String REDIS_PORT = "redis.port";
    public static final String REDIS_ENCODING = "redis.encoding";
    public static final String REDIS_TCP_KEEPALIVE = "redis.tcpKeepAlive";
    public static final String REDIS_TCP_NODELAY = "redis.tcpNoDelay";

    public static final String MONGODB_CONFIG = "mongodb.config";


    /**
     * extract config json file path from arguments
     * -conf path/to/config.json
     *
     * @param args
     * @return
     */
    public static String extractConfigPath(String[] args) {
        String cfgPath = "";
        for (int i = 0; i < args.length - 1; i++) {
            if (args[i].equals("-conf")) {
                cfgPath = args[i + 1];
                break;
            }
        }
        return cfgPath;
    }

    /**
     * load config json file and parse into json object
     *
     * @param vertx
     * @param fullPath
     * @param defaultConfig
     * @return
     */
    public static JsonObject loadConfigFromPath(Vertx vertx, String fullPath, JsonObject defaultConfig) {
        JsonObject configObject;
        if (fullPath != null && !fullPath.equals("")) {
            Buffer fileBuf = vertx.fileSystem().readFileBlocking(fullPath);
            configObject = new JsonObject(fileBuf);
        } else {
            configObject = defaultConfig;
        }

        return configObject;
    }

}
