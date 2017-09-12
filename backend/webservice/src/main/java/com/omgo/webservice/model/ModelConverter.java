package com.omgo.webservice.model;

import io.vertx.core.json.JsonObject;
import proto.Db;
import proto.common.Common;

import java.util.*;

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

    public static JsonObject key2Json(Db.DB.UserKey key) {
        return new JsonObject()
            .put(KEY_EMAIL, key.getEmail())
            .put(KEY_UID, key.getUid())
            .put(KEY_USN, key.getUsn());
    }

    public static Db.DB.UserKey json2UserKey(JsonObject jsonObject) {
        Db.DB.UserKey.Builder builder = Db.DB.UserKey.newBuilder();
        if (jsonObject.containsKey(KEY_USN)) {
            builder.setUsn(jsonObject.getLong(KEY_USN));
        }
        if (jsonObject.containsKey(KEY_UID)) {
            builder.setUid(jsonObject.getLong(KEY_UID));
        }
        if (jsonObject.containsKey(KEY_EMAIL)) {
            builder.setEmail(jsonObject.getString(KEY_EMAIL));
        }

        return builder.build();
    }

    public static Db.DB.UserExtendInfo json2UserExtendInfo(JsonObject jsonObject) {
        Db.DB.UserExtendInfo.Builder builder = Db.DB.UserExtendInfo.newBuilder();
        Common.UserInfo userInfo = json2UserInfo(jsonObject);
        builder.setInfo(userInfo)
            .setSecret(jsonObject.getString(KEY_SECRET))
            .setToken(jsonObject.getString(KEY_TOKEN));
        return builder.build();
    }

    public static String SQLQueryQueryUid(long uid) {
        return "SELECT * FROM user WHERE uid=" + uid;
    }

    private static Set<String> getUserMapKeySet() {
        Set<String> keySet = new HashSet<>();
        keySet.add(ModelConverter.KEY_UID);
        keySet.add(ModelConverter.KEY_AVATAR);
        keySet.add(ModelConverter.KEY_BIRTHDAY);
        keySet.add(ModelConverter.KEY_COUNTRY);
        keySet.add(ModelConverter.KEY_EMAIL);
        keySet.add(ModelConverter.KEY_GENDER);
        keySet.add(ModelConverter.KEY_LAST_LOGIN);
        keySet.add(ModelConverter.KEY_LOGIN_COUNT);
        keySet.add(ModelConverter.KEY_NICKNAME);
        keySet.add(ModelConverter.KEY_SALT);
        keySet.add(ModelConverter.KEY_SECRET);
        keySet.add(ModelConverter.KEY_SINCE);
        return keySet;
    }

    public static String SQLQueryInsert(JsonObject jsonObject) {
        String SQL_INSERT = "INSERT INTO user ";

        SQL_INSERT += toKeyValues(jsonObject, getUserMapKeySet());

        return SQL_INSERT;
    }

    public static String toKeyValues(JsonObject jsonObject, Set<String> keySet) {
        List<String> keys = new ArrayList<>();
        List<String> values = new ArrayList<>();

        Map<String, Object> map = jsonObject.getMap();
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            String key = entry.getKey();
            if (!keySet.contains(key)) {
                continue;
            }
            keys.add(key);
            Object value = entry.getValue();
            if (value instanceof String) {
                values.add(COMMA + (String)value + COMMA);
            } else {
                values.add(value.toString());
            }
        }

        return "(" + String.join(",", keys) + ") VALUES (" + String.join(",", values) + ")";
    }

    public static Set<String> getUserUpdatableMapKeySet() {
        Set<String> keySet = new HashSet<>();
        keySet.add(ModelConverter.KEY_AVATAR);
        keySet.add(ModelConverter.KEY_BIRTHDAY);
        keySet.add(ModelConverter.KEY_COUNTRY);
        keySet.add(ModelConverter.KEY_EMAIL);
        keySet.add(ModelConverter.KEY_GENDER);
        keySet.add(ModelConverter.KEY_LAST_LOGIN);
        keySet.add(ModelConverter.KEY_LOGIN_COUNT);
        keySet.add(ModelConverter.KEY_NICKNAME);
        return keySet;
    }
}
