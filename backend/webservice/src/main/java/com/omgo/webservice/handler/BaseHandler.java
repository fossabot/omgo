package com.omgo.webservice.handler;

import com.omgo.webservice.Utils;
import com.omgo.webservice.model.ModelConverter;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpMethod;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.Session;

import java.util.Map;

public class BaseHandler {
    protected String MIME_JSON = "application/json";
    protected String CONTENT_TYPE = "content-type";

    protected Logger LOGGER;
    protected Vertx vertx;

    protected Route route;

    protected String path;

    public BaseHandler(Vertx vertx) {
        this.vertx = vertx;
        LOGGER = LoggerFactory.getLogger(this.getClass());
    }

    public void initRoute(Router router, String path) {
        LOGGER.info("initRoute handler for : " + path);
        this.path = path;

        route = router.route(httpMethod(), path)
            .consumes(consumes())
            .produces(produces());
    }

    protected HttpServerRequest getRequest(RoutingContext context) {
        HttpServerRequest request = context.request();
        LOGGER.info("handling request: " + request.uri());
        LOGGER.info("header: " + getHeaderJson(request));
        return request;
    }

    protected HttpServerResponse getResponse(RoutingContext context) {
        HttpServerResponse response = context.response();
        // enable chunked responses because we will be adding data as
        // we execute over other handlers. This is only required once and
        // only if several handlers do output.
        response.setChunked(true);
        response.putHeader(CONTENT_TYPE, MIME_JSON);
        return response;
    }

    protected JsonObject getHeaderJson(HttpServerRequest request) {
        JsonObject headerJson = new JsonObject();
        for (Map.Entry<String, String> entry : request.headers().entries()) {
            headerJson.put(entry.getKey(), entry.getValue());
        }
        return headerJson;
    }

    protected boolean isSessionValid(RoutingContext context) {
        Session session = context.session();
        if (session != null) {
            JsonObject headerJson = getHeaderJson(getRequest(context));
            String clientToken = headerJson.getString(ModelConverter.KEY_TOKEN);
            String sessionToken = session.get(ModelConverter.KEY_TOKEN);
            if (!Utils.isEmptyString(clientToken) && !Utils.isEmptyString(sessionToken)) {
                return sessionToken.equals(clientToken);
            }
        }
        return false;
    }

    protected HttpMethod httpMethod() {
        return HttpMethod.GET;
    }

    protected String consumes() {
        return MIME_JSON;
    }

    protected String produces() {
        return MIME_JSON;
    }
}
