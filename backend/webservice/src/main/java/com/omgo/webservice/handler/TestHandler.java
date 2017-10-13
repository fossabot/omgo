package com.omgo.webservice.handler;

import com.omgo.webservice.AgentManager;
import com.omgo.webservice.Utils;
import com.omgo.webservice.model.ModelConverter;
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
    protected void handle(RoutingContext routingContext, HttpServerResponse response, JsonObject paramJson) {
        HttpServerRequest request = super.getRequest(routingContext);

        JsonObject rsp = getResponseJson();
        rsp.put("foo", "bar");
        rsp.put(ModelConverter.KEY_HOSTS, AgentManager.getInstance().getHostList());

        LOGGER.info(paramJson);

        response.write(rsp.encode()).end();
    }
}
