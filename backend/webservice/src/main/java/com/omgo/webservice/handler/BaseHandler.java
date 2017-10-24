package com.omgo.webservice.handler;

import com.omgo.webservice.Utils;
import com.omgo.webservice.model.HttpStatus;
import com.omgo.webservice.model.ModelConverter;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpMethod;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.DecodeException;
import io.vertx.core.json.JsonObject;
import io.vertx.core.logging.Logger;
import io.vertx.core.logging.LoggerFactory;
import io.vertx.ext.web.Route;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.Session;

import java.nio.charset.Charset;
import java.util.Map;
import java.util.TreeMap;

public class BaseHandler {
    protected String MIME_JSON = "application/json";
    protected String CONTENT_TYPE = "content-type";

    protected Logger LOGGER;
    protected Vertx vertx;

    protected Route route;

    protected String path;

    // security
    protected boolean requireValidSession;
    protected boolean requireValidNonce;
    protected boolean requireValidEncryption;


    public BaseHandler(Vertx vertx) {
        this.vertx = vertx;
        this.requireValidSession = true;
        this.requireValidNonce = true;
        this.requireValidEncryption = true;
        LOGGER = LoggerFactory.getLogger(this.getClass());
    }

    public void notRequireValidSession() {
        requireValidSession = false;
    }

    public void notRequireValidNonce() {
        requireValidSession = false;
    }

    public void notRequireValidEncryption() {
        requireValidEncryption = false;
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


        route.handler(routingContext -> {
            if (requireValidSession) {
                if (!isSessionValid(routingContext)) {
                    LOGGER.info("invalid session");
                    routingContext.fail(HttpStatus.UNAUTHORIZED.code);
                    return;
                }
            }

            String nonce = getValidNonce(routingContext);
            if (requireValidNonce && !Utils.DEBUG) {
                if (nonce == null) {
                    LOGGER.info("invalid nonce");
                    routingContext.fail(HttpStatus.UNAUTHORIZED.code);
                    return;
                }
            }

            if (nonce != null) {
                setSessionNonce(routingContext, nonce);
            }

            HttpServerResponse response = getResponse(routingContext);

            if (requireValidEncryption) {
                JsonObject paramJson = getRequestParam(routingContext);
                if (paramJson.isEmpty()) {
                    LOGGER.info("invalid param");
                    routingContext.fail(HttpStatus.BAD_REQUEST.code);
                    return;
                }
                handle(routingContext, response, paramJson);
            } else {
                handle(routingContext, response);
            }
        });
    }

    protected void handle(RoutingContext routingContext, HttpServerResponse response, JsonObject paramJson) {
    }

    protected void handle(RoutingContext routingContext, HttpServerResponse response) {

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
            headerJson.put(entry.getKey().toLowerCase(), entry.getValue());
        }
        return headerJson;
    }

    protected JsonObject getResponseJson() {
        JsonObject rspJson = new JsonObject();
        rspJson.put(ModelConverter.KEY_TIMESTAMP, System.currentTimeMillis());

        return rspJson;
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
     * set nonce
     *
     * @param routingContext
     * @param nonce
     * @return
     */
    protected Session setSessionNonce(RoutingContext routingContext, String nonce) {
        Session session = routingContext.session();
        if (session != null) {
            session.put(ModelConverter.KEY_NONCE, nonce);
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
            if (Utils.DEBUG) {
                return true;
            }
            JsonObject headerJson = getHeaderJson(getRequest(context));
            String clientToken = headerJson.getString(ModelConverter.KEY_TOKEN);
            String sessionToken = session.get(ModelConverter.KEY_TOKEN);
            if (!Utils.isEmptyString(clientToken) && !Utils.isEmptyString(sessionToken)) {
                return sessionToken.equals(clientToken);
            }
        }
        return false;
    }

    protected String getValidNonce(RoutingContext context) {
        Session session = context.session();
        if (session != null) {
            JsonObject headerJson = getHeaderJson(getRequest(context));
            String requestNonce = headerJson.getString(ModelConverter.KEY_NONCE);
            String sessionNonce = session.get(ModelConverter.KEY_NONCE);
            if (Utils.isEmptyString(requestNonce)) {
                return null;
            }

            if (Utils.isEmptyString(sessionNonce)) {
                return requestNonce;
            }

            try {
                long reqNonce = Long.parseLong(requestNonce);
                long sesNonce = Long.parseLong(sessionNonce);
                if (reqNonce > sesNonce) {
                    return requestNonce;
                } else {
                    return null;
                }
            } catch (Exception e) {
                LOGGER.info(e);
            }
        }
        return null;
    }

    protected JsonObject getRequestParam(RoutingContext context) {
        Session session = context.session();
        HttpServerRequest request = context.request();
        String paramStr = request.headers().get(ModelConverter.KEY_PARAM);
        String paramSignature = request.headers().get(ModelConverter.KEY_SIGNATURE);

        if (Utils.DEBUG) {
            return safeParseJsonString(paramStr);
        } else if (session != null) {
            byte[] key = session.get(ModelConverter.KEY_SEED);
            if (key != null) {
                // decrypt
                String decryptString = cryptionXOR(paramStr, key);
                JsonObject paramJson = safeParseJsonString(decryptString);
                if (!paramJson.isEmpty()) {
                    // check signature
                    byte[] signature = calculateSignature(paramJson);
                    String signatureBase64 = Utils.encodeBase64(signature);
                    if (signatureBase64.equals(paramSignature)) {
                        return paramJson;
                    }
                }
            }
        }

        return new JsonObject();
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

    protected JsonObject safeParseJsonString(String jsonStr) {
        try {
            JsonObject parsedJson = new JsonObject(jsonStr);
            return parsedJson;
        } catch (DecodeException e) {
            LOGGER.info(e);
        }
        return new JsonObject();
    }

    protected String cryptionXOR(String input, byte[] key) {
        Charset charSet = Charset.forName("UTF-8");
        byte[] inputBytes = input.getBytes(charSet);
        for (int i = 0; i < inputBytes.length; i++) {
            inputBytes[i] = (byte) (inputBytes[i] ^ key[i % key.length]);
        }
        return new String(inputBytes, charSet);
    }

    protected byte[] calculateSignature(JsonObject jsonObject) {
        if (jsonObject == null || jsonObject.isEmpty()) {
            return null;
        }

        TreeMap<String, Object> treeMap = new TreeMap<>(jsonObject.getMap());
        StringBuilder sb = new StringBuilder();
        for (String key : treeMap.keySet()) {
            sb.append(key);
            sb.append(treeMap.get(key));
        }
        String raw = sb.toString();
        return Utils.sha1(raw);
    }
}
