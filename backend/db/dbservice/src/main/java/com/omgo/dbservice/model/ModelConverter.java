package com.omgo.dbservice.model;

import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import proto.common.Common;

import java.util.ArrayList;
import java.util.List;

public class ModelConverter {
    public static final String KEY_AVATAR = "avatar";
    public static final String KEY_BIRTHDAY = "birthday";
    public static final String KEY_COUNTRY = "country";
    public static final String KEY_EMAIL = "email";
    public static final String KEY_GENDER = "gender";
    public static final String KEY_LAST_LOGIN = "last_login";
    public static final String KEY_LOGIN_COUNT = "login_count";
    public static final String KEY_NICKNAME = "nickname";
    public static final String KEY_SALT = "salt";
    public static final String KEY_SECRET = "secret";
    public static final String KEY_SINCE = "since";
    public static final String KEY_TOKEN = "token";
    public static final String KEY_UID = "uid";
    public static final String KEY_USN = "usn";

    private static final String COMMA = "'";

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

    public static String SQLQueryQueryUid(long uid) {
        return "SELECT * FROM user WHERE uid=" + uid;
    }

    public static String SQLQueryInsert(JsonObject jsonObject) {
        String SQL_INSERT = "INSERT INTO user (";
        String SQL_VALUES = "";

        List<String> VALUE_KEYS = new ArrayList<>();
        VALUE_KEYS.add(ModelConverter.KEY_UID);
        VALUE_KEYS.add(ModelConverter.KEY_AVATAR);
        VALUE_KEYS.add(ModelConverter.KEY_BIRTHDAY);
        VALUE_KEYS.add(ModelConverter.KEY_COUNTRY);
        VALUE_KEYS.add(ModelConverter.KEY_EMAIL);
        VALUE_KEYS.add(ModelConverter.KEY_GENDER);
        VALUE_KEYS.add(ModelConverter.KEY_LAST_LOGIN);
        VALUE_KEYS.add(ModelConverter.KEY_LOGIN_COUNT);
        VALUE_KEYS.add(ModelConverter.KEY_NICKNAME);
        VALUE_KEYS.add(ModelConverter.KEY_SALT);
        VALUE_KEYS.add(ModelConverter.KEY_SECRET);
        VALUE_KEYS.add(ModelConverter.KEY_SINCE);

        SQL_INSERT += String.join(",", VALUE_KEYS) + ") VALUES (";

        final long uid = jsonObject.getLong(ModelConverter.KEY_UID);
        SQL_INSERT += jsonObject.getLong(ModelConverter.KEY_UID) + ",";
        SQL_INSERT += COMMA + jsonObject.getString(ModelConverter.KEY_AVATAR) + COMMA + ",";
        SQL_INSERT += jsonObject.getLong(ModelConverter.KEY_BIRTHDAY) + ",";
        SQL_INSERT += COMMA + jsonObject.getString(ModelConverter.KEY_COUNTRY) + COMMA + ",";
        SQL_INSERT += COMMA + jsonObject.getString(ModelConverter.KEY_EMAIL) + COMMA + ",";
        SQL_INSERT += jsonObject.getInteger(ModelConverter.KEY_GENDER) + ",";
        SQL_INSERT += jsonObject.getLong(ModelConverter.KEY_LAST_LOGIN) + ",";
        SQL_INSERT += jsonObject.getLong(ModelConverter.KEY_LOGIN_COUNT) + ",";
        SQL_INSERT += COMMA + jsonObject.getString(ModelConverter.KEY_NICKNAME) + COMMA + ",";
        SQL_INSERT += COMMA + jsonObject.getString(ModelConverter.KEY_SALT) + COMMA + ",";
        SQL_INSERT += COMMA + jsonObject.getString(ModelConverter.KEY_SECRET) + COMMA + ",";
        SQL_INSERT += jsonObject.getLong(ModelConverter.KEY_SINCE) + ")";

        return SQL_INSERT;
    }
}
