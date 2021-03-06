package com.omgo.webservice.handler;

import com.omgo.utils.HttpStatus;
import com.omgo.utils.ModelKeys;
import com.omgo.utils.Utils;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.Session;
import org.whispersystems.curve25519.Curve25519;
import org.whispersystems.curve25519.Curve25519KeyPair;

public class HandshakeHandler extends BaseHandler {
    public HandshakeHandler(Vertx vertx) {
        super(vertx);
        supportStandaloneMode();
        notRequireValidEncryption();
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response) {
        HttpServerRequest request = super.getRequest(routingContext);
        JsonObject headerJson = getHeaderJson(request);

        String clientSeed = headerJson.getString(ModelKeys.SEED);

        if (Utils.isEmptyString(clientSeed)) {
            routingContext.fail(HttpStatus.FORBIDDEN.code);
            return;
        }

        byte[] clientSeedBytes;

        try {
            clientSeedBytes = Utils.decodeBase64(clientSeed);
        } catch (IllegalArgumentException e) {
            LOGGER.info(e);
            routingContext.fail(HttpStatus.FORBIDDEN.code);
            return;
        }

        Curve25519 cipher = Curve25519.getInstance(Curve25519.BEST);
        Curve25519KeyPair keyPair = cipher.generateKeyPair();

        byte[] sharedSecret = cipher.calculateAgreement(clientSeedBytes, keyPair.getPrivateKey());
        Session session = routingContext.session();
        session.put(ModelKeys.SEED, sharedSecret);

        JsonObject rspJson = getResponseJson();
        rspJson.put(ModelKeys.SEED, Utils.encodeBase64(keyPair.getPublicKey()));
        response.write(rspJson.encode()).end();
    }
}
