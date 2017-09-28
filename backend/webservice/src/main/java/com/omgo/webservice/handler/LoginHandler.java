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

public class LoginHandler extends BaseHandler {
    private ManagedChannel grpcChannel;

    public LoginHandler(Vertx vertx, ManagedChannel channel) {
        super(vertx);
        this.grpcChannel = channel;
    }

    @Override
    public void setRoute(Router router, String path) {
        super.setRoute(router, path);

        GRPCAuthProvider authProvider = new GRPCAuthProvider(vertx, grpcChannel);
        route.handler(routingContext -> {
            HttpServerRequest request = super.getRequest(routingContext);
            HttpServerResponse response = super.getResponse(routingContext);

            JsonObject authJson = super.getHeaderJson(request);

            // TODO: 15/09/2017 check parameters here
            // thought dbservice has already check these parameters
            // but just for break fast's cause

            authJson.put(ModelConverter.KEY_LAST_IP, request.connection().remoteAddress().host());
            authProvider.authenticate(authJson, authRes -> {
                if (authRes.succeeded()) {
                    User user = authRes.result();
                    routingContext.setUser(user);

                    setSessionToken(routingContext, user.principal().getString(ModelConverter.KEY_TOKEN));

                    response.write(user.principal().encode()).end();
                } else {
                    routingContext.fail(403);
                }
            });
        });
    }
}
