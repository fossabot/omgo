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

    /**
     * setup a route for path
     *
     * @param router
     * @param path
     */
    public void setRoute(Router router, String path) {
        LOGGER.info("setRoute handler for : " + path);
        this.path = path;

        route = router.route(httpMethod(), path)
            .consumes(consumes())
            .produces(produces());
    }

    /**
     * get request object from routing context, and log its info
     *
     * @param context
     * @return
     */
    protected HttpServerRequest getRequest(RoutingContext context) {
        HttpServerRequest request = context.request();
        LOGGER.info("handling request: " + request.uri());
        LOGGER.info("header: " + getHeaderJson(request));
        return request;
    }

    /**
     * get response object from routing context, and setup chunk, content-type etc.
     *
     * @param context
     * @return
     */
    protected HttpServerResponse getResponse(RoutingContext context) {
        HttpServerResponse response = context.response();
        // enable chunked responses because we will be adding data as
        // we execute over other handlers. This is only required once and
        // only if several handlers do output.
        response.setChunked(true);
        response.putHeader(CONTENT_TYPE, MIME_JSON);
        return response;
    }

    /**
     * get request header json object
     *
     * @param request
     * @return
     */
    protected JsonObject getHeaderJson(HttpServerRequest request) {
        JsonObject headerJson = new JsonObject();
        for (Map.Entry<String, String> entry : request.headers().entries()) {
            headerJson.put(entry.getKey(), entry.getValue());
        }
        return headerJson;
    }

    /**
     * set token
     *
     * @param routingContext
     * @param token
     * @return
     */
    protected Session setSessionToken(RoutingContext routingContext, String token) {
        Session session = routingContext.session();
        if (session != null) {
            session.regenerateId();
            session.put(ModelConverter.KEY_TOKEN, token);
        }
        return session;
    }

    /**
     * verify session via token
     *
     * @param context
     * @return
     */
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
