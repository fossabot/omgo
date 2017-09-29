package com.omgo.webservice.handler;

import com.omgo.webservice.Utils;
import com.omgo.webservice.model.ModelConverter;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.Session;
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

            Session session = routingContext.session();

            JsonObject headerJson = getHeaderJson(request);

            String clientSeed = headerJson.getString(ModelConverter.KEY_SEED);

            if (Utils.isEmptyString(clientSeed)) {
                routingContext.fail(403);
                return;
            }

            byte[] clientSeedBytes;

            try {
                clientSeedBytes = Utils.decodeBase64(clientSeed);
            } catch (IllegalArgumentException e) {
                LOGGER.info(e);
                routingContext.fail(403);
                return;
            }

            Curve25519 cipher = Curve25519.getInstance(Curve25519.BEST);
            Curve25519KeyPair keyPair = cipher.generateKeyPair();

            byte[] sharedSecret = cipher.calculateAgreement(clientSeedBytes, keyPair.getPrivateKey());
            session.put(ModelConverter.KEY_SEED, sharedSecret);

            JsonObject rspJson = getResponseJson();
            rspJson.put(ModelConverter.KEY_SEED, Utils.encodeBase64(keyPair.getPublicKey()));
            response.write(rspJson.encode()).end();
        });
    }
}
