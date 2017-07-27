package com.omgo.dbservice.model;

import io.vertx.core.json.JsonObject;
import proto.common.Common;

public class ModelConverter {
    public static final String KEY_AVATAR = "avatar";
    public static final String KEY_BIRTHDAY = "birthday";
    public static final String KEY_COUNTRY = "country";
    public static final String KEY_EMAIL = "email";
    public static final String KEY_GENDER = "gender";
    public static final String KEY_LAST_LOGIN = "lastLogin";
    public static final String KEY_LOGIN_COUNT = "loginCount";
    public static final String KEY_NICKNAME = "nickname";
    public static final String KEY_SINCE = "since";
    public static final String KEY_UID = "uid";
    public static final String KEY_USN = "usn";

    public static Common.RspHeader createRspHeader(String msg, int status, long timestamp) {
        return Common.RspHeader.newBuilder()
            .setMsg(msg)
            .setStatus(status)
            .setTimestamp(timestamp)
            .build();
    }

    public static Common.RspHeader createSuccessRspHeader() {
        return createRspHeader("", Common.ResultCode.RESULT_OK_VALUE, System.currentTimeMillis());
    }

    public static Common.UserInfo json2UserInfo(JsonObject jsonObject) {
        return Common.UserInfo.newBuilder()
            .setAvatar(jsonObject.getString(KEY_AVATAR, ""))
            .setBirthday(jsonObject.getLong(KEY_BIRTHDAY, 0L))
            .setCountry(jsonObject.getString(KEY_COUNTRY, ""))
            .setEmail(jsonObject.getString(KEY_EMAIL, ""))
            .setGender(Common.Gender.forNumber(jsonObject.getInteger(KEY_GENDER, 0)))
            .setLastLogin(jsonObject.getLong(KEY_LAST_LOGIN, 0L))
            .setLoginCount(jsonObject.getInteger(KEY_LOGIN_COUNT, 0))
            .setNickname(jsonObject.getString(KEY_NICKNAME, ""))
            .setSince(jsonObject.getLong(KEY_SINCE, 0L))
            .setUid(jsonObject.getLong(KEY_UID, 0L))
            .setUsn(jsonObject.getLong(KEY_USN, 0L))
            .build();
    }

    public static JsonObject userInfo2Json(Common.UserInfo userInfo) {
        return new JsonObject()
            .put(KEY_AVATAR, userInfo.getAvatar())
            .put(KEY_BIRTHDAY, userInfo.getBirthday())
            .put(KEY_COUNTRY, userInfo.getCountry())
            .put(KEY_EMAIL, userInfo.getEmail())
            .put(KEY_GENDER, userInfo.getGenderValue())
            .put(KEY_LAST_LOGIN, userInfo.getLastLogin())
            .put(KEY_LOGIN_COUNT, userInfo.getLoginCount())
            .put(KEY_NICKNAME, userInfo.getNickname())
            .put(KEY_SINCE, userInfo.getSince())
            .put(KEY_UID, userInfo.getUid())
            .put(KEY_USN, userInfo.getUsn());
    }
}
