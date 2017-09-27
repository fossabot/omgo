package com.omgo.webservice.handler;

import com.omgo.webservice.GRPCAuthProvider;
import com.omgo.webservice.model.ModelConverter;
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
    public void initRoute(Router router, String path) {
        super.initRoute(router, path);

        GRPCAuthProvider authProvider = new GRPCAuthProvider(vertx, grpcChannel);
        route.handler(routingContext -> {
            HttpServerRequest request = super.getRequest(routingContext);
            HttpServerResponse response = super.getResponse(routingContext);

            // authenticate
            Session session = routingContext.session();

            JsonObject authJson = super.getHeaderJson(request);

            // TODO: 15/09/2017 check parameters here
            // thought dbservice has already check these parameters
            // but just for break fast's cause

            authJson.put(ModelConverter.KEY_LAST_IP, request.connection().remoteAddress().host());
            authProvider.authenticate(authJson, authRes -> {
                if (authRes.succeeded()) {
                    User user = authRes.result();
                    routingContext.setUser(user);

                    session.put(ModelConverter.KEY_TOKEN, user.principal().getString(ModelConverter.KEY_TOKEN));
                    session.regenerateId();

                    response.write(user.principal().encode()).end();
                } else {
                    routingContext.fail(403);
                }
            });
        });
    }
}
