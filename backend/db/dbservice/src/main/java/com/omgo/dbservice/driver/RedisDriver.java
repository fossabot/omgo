package com.omgo.dbservice.driver;

import io.vertx.redis.RedisClient;

public class RedisDriver {

    private RedisClient client;

    public RedisDriver(RedisClient client) {
        this.client = client;
    }
}
