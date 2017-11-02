package com.omgo.webservice.handler;

import com.omgo.utils.HttpStatus;
import com.omgo.utils.ModelKeys;
import com.omgo.utils.Services;
import com.omgo.webservice.model.ModelConverter;
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

        long usn = paramJson.getLong(ModelKeys.USN);

        Db.DB.UserEntry.Builder builder = Db.DB.UserEntry.newBuilder();
        builder.setUsn(usn);
        dbServiceVertxStub.userQuery(builder.build(), res -> {
            if (res.succeeded()) {
                int code = res.result().getResult().getStatus();
                if (code == Db.DB.StatusCode.STATUS_OK_VALUE) {
                    JsonObject resultJson = ModelConverter.userEntry2Json(res.result().getUser());

                    // filter key-values
                    Set<String> keySet = new HashSet<>();
                    keySet.add(ModelKeys.USN);
                    keySet.add(ModelKeys.UID);
                    keySet.add(ModelKeys.AVATAR);
                    keySet.add(ModelKeys.BIRTHDAY);
                    keySet.add(ModelKeys.COUNTRY);
                    keySet.add(ModelKeys.EMAIL_VERIFIED);
                    keySet.add(ModelKeys.GENDER);
                    keySet.add(ModelKeys.IS_OFFICIAL);
                    keySet.add(ModelKeys.IS_ROBOT);
                    keySet.add(ModelKeys.LAST_IP);
                    keySet.add(ModelKeys.LAST_LOGIN);
                    keySet.add(ModelKeys.MCC);
                    keySet.add(ModelKeys.NICKNAME);
                    keySet.add(ModelKeys.PHONE_VERIFIED);
                    keySet.add(ModelKeys.PREMIUM_END);
                    keySet.add(ModelKeys.PREMIUM_EXP);
                    keySet.add(ModelKeys.PREMIUM_LEVEL);
                    keySet.add(ModelKeys.SINCE);
                    keySet.add(ModelKeys.SOCIAL_VERIFIED);
                    keySet.add(ModelKeys.STATUS);
                    keySet.add(ModelKeys.TIMEZONE);

                    JsonObject filteredResult = new JsonObject();
                    resultJson.getMap().forEach((k, v) -> {
                        if (keySet.contains(k)) {
                            filteredResult.put(k, v);
                        }
                    });

                    JsonObject rspJson = getResponseJson();
                    rspJson.put(ModelKeys.USER_INFO, filteredResult);
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
