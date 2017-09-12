package com.omgo.webservice.handler;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import com.omgo.webservice.Utils;
import com.omgo.webservice.model.ModelConverter;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import proto.DBServiceGrpc;
import proto.Db;
import proto.common.Common;

public class RegisterHandler extends BaseHandler {

    private DBServiceGrpc.DBServiceVertxStub dbServiceVertxStub;

    public RegisterHandler(Vertx vertx, ManagedChannel channel) {
        super(vertx);
        dbServiceVertxStub = DBServiceGrpc.newVertxStub(channel);
    }

    @Override
    public void register(Router router, String path) {
        super.register(router, path);

        route.handler(routingContext -> {
            HttpServerRequest request = super.handle(routingContext);
            HttpServerResponse response = super.response(routingContext);

            JsonObject registerJson = super.getHeaderJson(request);
            String avatar = registerJson.getString(ModelConverter.KEY_AVATAR);
            String birthday = registerJson.getString(ModelConverter.KEY_BIRTHDAY);
            String country = registerJson.getString(ModelConverter.KEY_COUNTRY);
            String email = registerJson.getString(ModelConverter.KEY_EMAIL);
            String gender = registerJson.getString(ModelConverter.KEY_GENDER);
            String nickname = registerJson.getString(ModelConverter.KEY_NICKNAME);
            String secret = registerJson.getString(ModelConverter.KEY_SECRET);

            long birthdayLong = Utils.isEmptyString(birthday) ? 0L : Long.parseLong(birthday);
            int genderInt = Utils.isEmptyString(gender) ? 0 : Integer.parseInt(gender);

            Common.UserInfo.Builder userInfoBuilder = Common.UserInfo.newBuilder();
            userInfoBuilder
                .setAvatar(avatar)
                .setBirthday(birthdayLong)
                .setCountry(country)
                .setEmail(email)
                .setGender(Common.Gender.forNumber(genderInt))
                .setNickname(nickname);

            Db.DB.UserExtendInfo.Builder builder = Db.DB.UserExtendInfo.newBuilder();
            Db.DB.UserExtendInfo extendInfo = builder
                .setInfo(userInfoBuilder.build())
                .setSecret(secret)
                .build();

            dbServiceVertxStub.userRegister(extendInfo, res -> {
                if (res.succeeded()) {
                    JsonObject resultJson = new JsonObject();
                    try {
                        String result = JsonFormat.printer().print(res.result());
                        resultJson = new JsonObject(result);
                    } catch (InvalidProtocolBufferException e) {
                        e.printStackTrace();
                    }
                    response.write(resultJson.encode()).end();
                } else {
                    routingContext.fail(500);
                }
            });
        });
    }
}
