package com.omgo.dbservice.driver;

import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.sql.SQLClient;
import io.vertx.ext.sql.SQLConnection;

public class MySQLDriver {

    private static final Logger LOGGER = LoggerFactory.getLogger(MySQLDriver.class);

    private SQLClient client;

    public MySQLDriver(SQLClient client) {
        this.client = client;
    }

    public void queryUser(long usn) {
        client.getConnection(res -> {
            if (res.succeeded()) {
                SQLConnection connection = res.result();

                connection.close();
            } else {
                LOGGER.error(res.cause());
            }
        });
    }
}
