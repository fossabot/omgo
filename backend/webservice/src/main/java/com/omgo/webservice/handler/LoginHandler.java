package com.omgo.webservice.handler;

import com.omgo.webservice.GRPCAuthProvider;
import com.omgo.webservice.model.HttpStatus;
import com.omgo.webservice.model.ModelConverter;
import com.omgo.webservice.service.Services;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.auth.User;
import io.vertx.ext.web.RoutingContext;

public class LoginHandler extends BaseHandler {

    private GRPCAuthProvider authProvider;

    public LoginHandler(Vertx vertx, Services.Pool servicePool) {
        super(vertx);
        notRequireValidNonce();
        notRequireValidSession();
        notRequireValidEncryption();
        this.authProvider = new GRPCAuthProvider(vertx, servicePool);
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response) {
        HttpServerRequest request = super.getRequest(routingContext);

        JsonObject authJson = super.getHeaderJson(request);

        // TODO: 15/09/2017 check parameters here
        // thought dataservice has already check these parameters
        // but just for break fast's cause

        authJson.put(ModelConverter.KEY_LAST_IP, request.connection().remoteAddress().host());
        authProvider.authenticate(authJson, authRes -> {
            if (authRes.succeeded()) {
                User user = authRes.result();
                routingContext.setUser(user);

                setSessionToken(routingContext, user.principal().getString(ModelConverter.KEY_TOKEN));

                JsonObject rspJson = getResponseJson();
                rspJson.put(ModelConverter.KEY_USER_INFO, user.principal());
                response.write(rspJson.encode()).end();
            } else {
                routingContext.fail(HttpStatus.FORBIDDEN.code);
            }
        });
    }
}
