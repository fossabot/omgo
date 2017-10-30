package com.omgo.webservice.handler;

import com.omgo.webservice.model.HttpStatus;
import com.omgo.webservice.model.ModelConverter;
import com.omgo.webservice.service.Services;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;
import proto.Db;

import java.util.HashSet;
import java.util.Set;

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

                    // filter key-values
                    Set<String> keySet = new HashSet<>();
                    keySet.add(ModelConverter.KEY_USN);
                    keySet.add(ModelConverter.KEY_UID);
                    keySet.add(ModelConverter.KEY_AVATAR);
                    keySet.add(ModelConverter.KEY_BIRTHDAY);
                    keySet.add(ModelConverter.KEY_COUNTRY);
                    keySet.add(ModelConverter.KEY_EMAIL_VERIFIED);
                    keySet.add(ModelConverter.KEY_GENDER);
                    keySet.add(ModelConverter.KEY_IS_OFFICIAL);
                    keySet.add(ModelConverter.KEY_IS_ROBOT);
                    keySet.add(ModelConverter.KEY_LAST_IP);
                    keySet.add(ModelConverter.KEY_LAST_LOGIN);
                    keySet.add(ModelConverter.KEY_MCC);
                    keySet.add(ModelConverter.KEY_NICKNAME);
                    keySet.add(ModelConverter.KEY_PHONE_VERIFIED);
                    keySet.add(ModelConverter.KEY_PREMIUM_END);
                    keySet.add(ModelConverter.KEY_PREMIUM_EXP);
                    keySet.add(ModelConverter.KEY_PREMIUM_LEVEL);
                    keySet.add(ModelConverter.KEY_SINCE);
                    keySet.add(ModelConverter.KEY_SOCIAL_VERIFIED);
                    keySet.add(ModelConverter.KEY_STATUS);
                    keySet.add(ModelConverter.KEY_TIMEZONE);

                    JsonObject filteredResult = new JsonObject();
                    resultJson.getMap().forEach((k, v) -> {
                        if (keySet.contains(k)) {
                            filteredResult.put(k, v);
                        }
                    });

                    JsonObject rspJson = getResponseJson();
                    rspJson.put(ModelConverter.KEY_USER_INFO, filteredResult);
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
