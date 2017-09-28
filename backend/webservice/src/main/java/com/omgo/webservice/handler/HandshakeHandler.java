package com.omgo.webservice.handler;

import com.omgo.webservice.Utils;
import com.omgo.webservice.model.ModelConverter;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import org.whispersystems.curve25519.Curve25519;
import org.whispersystems.curve25519.Curve25519KeyPair;

public class HandshakeHandler extends BaseHandler {
    public HandshakeHandler(Vertx vertx) {
        super(vertx);
    }

    @Override
    public void setRoute(Router router, String path) {
        super.setRoute(router, path);

        route.handler(routingContext -> {
            HttpServerRequest request = super.getRequest(routingContext);
            HttpServerResponse response = super.getResponse(routingContext);

            if (!isSessionValid(routingContext)) {
                routingContext.fail(401);
                return;
            }

            JsonObject headerJson = getHeaderJson(request);

            String clientSendSeed = headerJson.getString(ModelConverter.KEY_SEND_SEED);
            String clientRecvSeed = headerJson.getString(ModelConverter.KEY_RECV_SEED);

            if (Utils.isEmptyString(clientSendSeed) || Utils.isEmptyString(clientRecvSeed)) {
                routingContext.fail(403);
                return;
            }

            Curve25519 cipher = Curve25519.getInstance(Curve25519.BEST);
            Curve25519KeyPair keyPair = cipher.generateKeyPair();

            JsonObject rspJson = new JsonObject();
            rspJson.put(ModelConverter.KEY_SEND_SEED, Utils.encodeBase64(keyPair.getPrivateKey()));
            rspJson.put(ModelConverter.KEY_RECV_SEED, Utils.encodeBase64(keyPair.getPublicKey()));
            response.write(rspJson.encode()).end();
        });
    }
}
