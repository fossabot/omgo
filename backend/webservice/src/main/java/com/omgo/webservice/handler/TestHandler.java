package com.omgo.webservice.handler;

import com.omgo.webservice.Utils;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;

public class TestHandler extends BaseHandler {
    public TestHandler(Vertx vertx) {
        super(vertx);
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response) {
        HttpServerRequest request = super.getRequest(routingContext);

        JsonObject rsp = getResponseJson();
        rsp.put("foo", "bar");
        JsonObject headerJson = getHeaderJson(request);

        JsonObject test = new JsonObject();
        test.put("cdef", "sdfasdf");
        test.put("adfas", 53);
        test.put("erqe", "dfasdf");
        test.put("hdfaier", false);
        test.put("ndrwe", "dfasdf");
        test.put("djfdsd", 123456L);

        byte[] b = calculateSignature(test);
        String bb = Utils.encodeBase64(b);
        LOGGER.info(bb);

        response.write(rsp.encode()).end();
    }
}
