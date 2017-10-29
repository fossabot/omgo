package com.omgo.webservice.handler;

import com.omgo.webservice.model.HttpStatus;
import com.omgo.webservice.model.ModelConverter;
import com.omgo.webservice.service.Services;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;
import proto.Db;

public class UserProfileHandler extends BaseGrpcHandler {

    public UserProfileHandler(Vertx vertx, Services.Pool servicePool) {
        super(vertx, servicePool);
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response, JsonObject paramJson) {
        super.handle(routingContext, response);
        if (dbServiceVertxStub == null) {
            return;
        }

        long usn = paramJson.getLong(ModelConverter.KEY_USN);

        Db.DB.UserEntry.Builder builder = Db.DB.UserEntry.newBuilder();
        builder.setUsn(usn);
        dbServiceVertxStub.userQuery(builder.build(), res -> {
            if (res.succeeded()) {
                int code = res.result().getResult().getStatus();
                if (code == Db.DB.StatusCode.STATUS_OK_VALUE) {
                    JsonObject resultJson = ModelConverter.userEntry2Json(res.result().getUser());

                    JsonObject rspJson = getResponseJson();
                    rspJson.put(ModelConverter.KEY_USER_INFO, resultJson);
                    response.write(rspJson.encode()).end();
                } else {
                    LOGGER.info(res.result().getResult());
                    routingContext.fail(HttpStatus.INTERNAL_SERVER_ERROR.code);
                }
            } else {
                LOGGER.info(res.cause());
                routingContext.fail(HttpStatus.INTERNAL_SERVER_ERROR.code);
            }
        });
    }
}
