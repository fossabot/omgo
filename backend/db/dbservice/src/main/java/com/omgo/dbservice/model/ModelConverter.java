package com.omgo.dbservice.model;

import io.vertx.core.json.JsonObject;
import proto.common.Common;

public class ModelConverter {
    public static Common.UserInfo toUserInfo(JsonObject jsonObject) {
        return Common.UserInfo.newBuilder()
            .setAvatar(jsonObject.getString("avatar", ""))
            .setBirthday(jsonObject.getLong("birthday", 0L))
            .setCountry(jsonObject.getString("country", ""))
            .setEmail(jsonObject.getString("email", ""))
            .setGender(Common.Gender.forNumber(jsonObject.getInteger("gender", 0)))
            .setLastLogin(jsonObject.getLong("lastLogin", 0L))
            .setLoginCount(jsonObject.getInteger("loginCount", 0))
            .setNickname(jsonObject.getString("nickname", ""))
            .setSince(jsonObject.getLong("since", 0L))
            .setUid(jsonObject.getLong("uid", 0L))
            .setUsn(jsonObject.getLong("usn", 0L))
            .build();
    }
}
