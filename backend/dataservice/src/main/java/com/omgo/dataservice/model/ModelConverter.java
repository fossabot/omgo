package com.omgo.dataservice.model;

import io.vertx.core.json.JsonObject;
import proto.Db;

import java.util.HashSet;
import java.util.Set;

public class ModelConverter {
    public static final String KEY_USER = "user";

    public static final String KEY_USN = "usn";
    public static final String KEY_UID = "uid";
    public static final String KEY_APP_LANGUAGE = "app_language";
    public static final String KEY_APP_VERSION = "app_version";
    public static final String KEY_AVATAR = "avatar";
    public static final String KEY_BIRTHDAY = "birthday";
    public static final String KEY_COUNTRY = "country";
    public static final String KEY_DEVICE_TYPE = "device_type";
    public static final String KEY_EMAIL = "email";
    public static final String KEY_EMAIL_VERIFIED = "email_verified";
    public static final String KEY_GENDER = "gender";
    public static final String KEY_IS_OFFICIAL = "is_official";
    public static final String KEY_IS_ROBOT = "is_robot";
    public static final String KEY_LAST_IP = "last_ip";
    public static final String KEY_LAST_LOGIN = "last_login";
    public static final String KEY_LOGIN_COUNT = "login_count";
    public static final String KEY_MCC = "mcc";
    public static final String KEY_NICKNAME = "nickname";
    public static final String KEY_OS = "os";
    public static final String KEY_OS_LOCALE = "os_locale";
    public static final String KEY_PHONE = "phone";
    public static final String KEY_PHONE_VERIFIED = "phone_verified";
    public static final String KEY_PREMIUM_END = "premium_end";
    public static final String KEY_PREMIUM_EXP = "premium_exp";
    public static final String KEY_PREMIUM_LEVEL = "premium_level";
    public static final String KEY_SECRET = "secret";
    public static final String KEY_SINCE = "since";
    public static final String KEY_SOCIAL_ID = "social_id";
    public static final String KEY_SOCIAL_NAME = "social_name";
    public static final String KEY_SOCIAL_VERIFIED = "social_verified";
    public static final String KEY_STATUS = "status";
    public static final String KEY_TIMEZONE = "timezone";
    public static final String KEY_TOKEN = "token";

    public static final String KEY_PARAM = "param";
    public static final String KEY_NONCE = "nonce";
    public static final String KEY_SEED = "seed";
    public static final String KEY_SIGNATURE = "signature";
    public static final String KEY_TIMESTAMP = "timestamp";
    public static final String KEY_USER_INFO = "user_info";

    public static final String KEY_HOSTS = "hosts";

    public static Db.DB.UserEntry json2UserEntry(JsonObject jsonObject) {
        return Db.DB.UserEntry.newBuilder()
            .setUsn(jsonObject.getLong(KEY_USN, 0L))
            .setUid(jsonObject.getLong(KEY_UID, 0L))
            .setAppLanguage(jsonObject.getString(KEY_APP_LANGUAGE, ""))
            .setAppVersion(jsonObject.getString(KEY_APP_VERSION, ""))
            .setAvatar(jsonObject.getString(KEY_AVATAR, ""))
            .setBirthday(jsonObject.getLong(KEY_BIRTHDAY, 0L))
            .setCountry(jsonObject.getString(KEY_COUNTRY, ""))
            .setDeviceType(jsonObject.getInteger(KEY_DEVICE_TYPE, 0))
            .setEmail(jsonObject.getString(KEY_EMAIL, ""))
            .setEmailVerified(jsonObject.getBoolean(KEY_EMAIL_VERIFIED, false))
            .setGender(jsonObject.getInteger(KEY_GENDER, 0))
            .setIsOfficial(jsonObject.getBoolean(KEY_IS_OFFICIAL, false))
            .setIsRobot(jsonObject.getBoolean(KEY_IS_ROBOT, false))
            .setLastIp(jsonObject.getString(KEY_LAST_IP, ""))
            .setLastLogin(jsonObject.getLong(KEY_LAST_LOGIN, 0L))
            .setLoginCount(jsonObject.getLong(KEY_LOGIN_COUNT, 0L))
            .setMcc(jsonObject.getInteger(KEY_MCC, 0))
            .setNickname(jsonObject.getString(KEY_NICKNAME, ""))
            .setOs(jsonObject.getString(KEY_OS, ""))
            .setOsLocale(jsonObject.getString(KEY_OS_LOCALE, ""))
            .setPhone(jsonObject.getString(KEY_PHONE, ""))
            .setPhoneVerified(jsonObject.getBoolean(KEY_PHONE_VERIFIED, false))
            .setPremiumEnd(jsonObject.getLong(KEY_PREMIUM_END, 0L))
            .setPremiumExp(jsonObject.getLong(KEY_PREMIUM_EXP, 0L))
            .setPremiumLevel(jsonObject.getInteger(KEY_PREMIUM_LEVEL, 0))
            .setSecret(jsonObject.getString(KEY_SECRET, ""))
            .setSince(jsonObject.getLong(KEY_SINCE, 0L))
            .setSocialId(jsonObject.getString(KEY_SOCIAL_ID, ""))
            .setSocialName(jsonObject.getString(KEY_SOCIAL_NAME, ""))
            .setSocialVerified(jsonObject.getBoolean(KEY_SOCIAL_VERIFIED, false))
            .setStatus(jsonObject.getInteger(KEY_STATUS, 0))
            .setTimezone(jsonObject.getInteger(KEY_TIMEZONE, 0))
            .setToken(jsonObject.getString(KEY_TOKEN, ""))
            .build();
    }

    public static JsonObject userEntry2Json(Db.DB.UserEntry userEntry) {
        return new JsonObject()
            .put(KEY_USN, userEntry.getUsn())
            .put(KEY_UID, userEntry.getUid())
            .put(KEY_APP_LANGUAGE, userEntry.getAppLanguage())
            .put(KEY_APP_VERSION, userEntry.getAppVersion())
            .put(KEY_AVATAR, userEntry.getAvatar())
            .put(KEY_BIRTHDAY, userEntry.getBirthday())
            .put(KEY_COUNTRY, userEntry.getCountry())
            .put(KEY_DEVICE_TYPE, userEntry.getDeviceType())
            .put(KEY_EMAIL, userEntry.getEmail())
            .put(KEY_EMAIL_VERIFIED, userEntry.getEmailVerified())
            .put(KEY_GENDER, userEntry.getGender())
            .put(KEY_IS_OFFICIAL, userEntry.getIsOfficial())
            .put(KEY_IS_ROBOT, userEntry.getIsRobot())
            .put(KEY_LAST_IP, userEntry.getLastIp())
            .put(KEY_LAST_LOGIN, userEntry.getLastLogin())
            .put(KEY_LOGIN_COUNT, userEntry.getLoginCount())
            .put(KEY_MCC, userEntry.getMcc())
            .put(KEY_NICKNAME, userEntry.getNickname())
            .put(KEY_OS, userEntry.getOs())
            .put(KEY_OS_LOCALE, userEntry.getOsLocale())
            .put(KEY_PHONE, userEntry.getPhone())
            .put(KEY_PHONE_VERIFIED, userEntry.getPhoneVerified())
            .put(KEY_PREMIUM_END, userEntry.getPremiumEnd())
            .put(KEY_PREMIUM_EXP, userEntry.getPremiumExp())
            .put(KEY_PREMIUM_LEVEL, userEntry.getPremiumLevel())
            .put(KEY_SECRET, userEntry.getSecret())
            .put(KEY_SINCE, userEntry.getSince())
            .put(KEY_SOCIAL_ID, userEntry.getSocialId())
            .put(KEY_SOCIAL_NAME, userEntry.getSocialName())
            .put(KEY_SOCIAL_VERIFIED, userEntry.getSocialVerified())
            .put(KEY_STATUS, userEntry.getStatus())
            .put(KEY_TIMEZONE, userEntry.getTimezone())
            .put(KEY_TOKEN, userEntry.getToken());
    }

    private static Set<String> getUserMapKeySet() {
        Set<String> keySet = new HashSet<>();
        keySet.add(KEY_USN);
        keySet.add(KEY_UID);
        keySet.add(KEY_APP_LANGUAGE);
        keySet.add(KEY_APP_VERSION);
        keySet.add(KEY_AVATAR);
        keySet.add(KEY_BIRTHDAY);
        keySet.add(KEY_COUNTRY);
        keySet.add(KEY_DEVICE_TYPE);
        keySet.add(KEY_EMAIL);
        keySet.add(KEY_EMAIL_VERIFIED);
        keySet.add(KEY_GENDER);
        keySet.add(KEY_IS_OFFICIAL);
        keySet.add(KEY_IS_ROBOT);
        keySet.add(KEY_LAST_IP);
        keySet.add(KEY_LAST_LOGIN);
        keySet.add(KEY_LOGIN_COUNT);
        keySet.add(KEY_MCC);
        keySet.add(KEY_NICKNAME);
        keySet.add(KEY_OS);
        keySet.add(KEY_OS_LOCALE);
        keySet.add(KEY_PHONE);
        keySet.add(KEY_PHONE_VERIFIED);
        keySet.add(KEY_PREMIUM_END);
        keySet.add(KEY_PREMIUM_EXP);
        keySet.add(KEY_PREMIUM_LEVEL);
        keySet.add(KEY_SECRET);
        keySet.add(KEY_SINCE);
        keySet.add(KEY_SOCIAL_ID);
        keySet.add(KEY_SOCIAL_NAME);
        keySet.add(KEY_SOCIAL_VERIFIED);
        keySet.add(KEY_STATUS);
        keySet.add(KEY_TIMEZONE);
        keySet.add(KEY_TOKEN);
        return keySet;
    }

    public static JsonObject removeKeysForLoginResponse(JsonObject jsonObject) {
        jsonObject.remove(KEY_APP_LANGUAGE);
        jsonObject.remove(KEY_APP_VERSION);
        jsonObject.remove(KEY_DEVICE_TYPE);
        jsonObject.remove(KEY_MCC);
        jsonObject.remove(KEY_OS);
        jsonObject.remove(KEY_OS_LOCALE);
        jsonObject.remove(KEY_SECRET);
        jsonObject.remove(KEY_SOCIAL_ID);
        jsonObject.remove(KEY_SOCIAL_NAME);
        jsonObject.remove(KEY_TIMEZONE);

        return jsonObject;
    }

    public static JsonObject removeKeysForRegisterRequest(JsonObject jsonObject) {
        jsonObject.remove(KEY_USN);
        jsonObject.remove(KEY_UID);
        jsonObject.remove(KEY_EMAIL_VERIFIED);
        jsonObject.remove(KEY_IS_OFFICIAL);
        jsonObject.remove(KEY_IS_ROBOT);
        jsonObject.remove(KEY_PHONE_VERIFIED);
        jsonObject.remove(KEY_PREMIUM_END);
        jsonObject.remove(KEY_PREMIUM_EXP);
        jsonObject.remove(KEY_PREMIUM_LEVEL);
        jsonObject.remove(KEY_SECRET);
        jsonObject.remove(KEY_SINCE);
        jsonObject.remove(KEY_SOCIAL_VERIFIED);
        jsonObject.remove(KEY_TIMEZONE);

        return jsonObject;
    }

    public static Set<String> getUserUpdatableMapKeySet() {
        Set<String> keySet = new HashSet<>();
        keySet.add(KEY_APP_LANGUAGE);
        keySet.add(KEY_APP_VERSION);
        keySet.add(KEY_AVATAR);
        keySet.add(KEY_BIRTHDAY);
        keySet.add(KEY_COUNTRY);
        keySet.add(KEY_DEVICE_TYPE);
        keySet.add(KEY_EMAIL);
        keySet.add(KEY_EMAIL_VERIFIED);
        keySet.add(KEY_GENDER);
        keySet.add(KEY_LAST_IP);
        keySet.add(KEY_LAST_LOGIN);
        keySet.add(KEY_LOGIN_COUNT);
        keySet.add(KEY_MCC);
        keySet.add(KEY_NICKNAME);
        keySet.add(KEY_OS);
        keySet.add(KEY_OS_LOCALE);
        keySet.add(KEY_PHONE);
        keySet.add(KEY_PHONE_VERIFIED);
        keySet.add(KEY_PREMIUM_END);
        keySet.add(KEY_PREMIUM_EXP);
        keySet.add(KEY_PREMIUM_LEVEL);
        keySet.add(KEY_SECRET);
        keySet.add(KEY_SOCIAL_ID);
        keySet.add(KEY_SOCIAL_NAME);
        keySet.add(KEY_SOCIAL_VERIFIED);
        keySet.add(KEY_STATUS);
        keySet.add(KEY_TIMEZONE);
        keySet.add(KEY_TOKEN);
        return keySet;
    }
}
