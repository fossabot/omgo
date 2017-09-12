package com.omgo.webservice.handler;

import com.omgo.webservice.GRPCAuthProvider;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.auth.User;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.Session;

public class LoginHandler extends BaseHandler {
    private ManagedChannel grpcChannel;

    public LoginHandler(Vertx vertx, ManagedChannel channel) {
        super(vertx);
        this.grpcChannel = channel;
    }

    @Override
    public void register(Router router, String path) {
        super.register(router, path);

        GRPCAuthProvider authProvider = new GRPCAuthProvider(vertx, grpcChannel);

        route.handler(routingContext -> {
            HttpServerRequest request = super.handle(routingContext);
            HttpServerResponse response = super.response(routingContext);

            // authenticate
            Session session = routingContext.session();

            JsonObject authJson = super.getHeaderJson(request);
            authProvider.authenticate(authJson, authRes -> {
                if (authRes.succeeded()) {
                    User user = authRes.result();
                    routingContext.setUser(user);
                    if (session != null) {
                        session.regenerateId();
                    }
                    response.write(user.principal().encode()).end();
                } else {
                    routingContext.fail(403);
                }
            });
        });
    }
}
