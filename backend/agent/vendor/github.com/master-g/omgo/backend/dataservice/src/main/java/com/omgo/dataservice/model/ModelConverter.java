package com.omgo.dataservice.model;

import com.omgo.utils.ModelKeys;
import io.vertx.core.json.JsonObject;
import proto.Db;

import java.util.HashSet;
import java.util.Set;

public class ModelConverter {

    public static Db.DB.UserEntry json2UserEntry(JsonObject jo) {
        return Db.DB.UserEntry.newBuilder()
            .setUsn(jo.getLong(ModelKeys.USN, 0L))
            .setUid(jo.getLong(ModelKeys.UID, 0L))
            .setAppLanguage(jo.getString(ModelKeys.APP_LANGUAGE, ""))
            .setAppVersion(jo.getString(ModelKeys.APP_VERSION, ""))
            .setAvatar(jo.getString(ModelKeys.AVATAR, ""))
            .setBirthday(jo.getLong(ModelKeys.BIRTHDAY, 0L))
            .setCountry(jo.getString(ModelKeys.COUNTRY, ""))
            .setDeviceType(jo.getInteger(ModelKeys.DEVICE_TYPE, 0))
            .setEmail(jo.getString(ModelKeys.EMAIL, ""))
            .setEmailVerified(jo.getBoolean(ModelKeys.EMAIL_VERIFIED, false))
            .setGender(jo.getInteger(ModelKeys.GENDER, 0))
            .setIsOfficial(jo.getBoolean(ModelKeys.IS_OFFICIAL, false))
            .setIsRobot(jo.getBoolean(ModelKeys.IS_ROBOT, false))
            .setLastIp(jo.getString(ModelKeys.LAST_IP, ""))
            .setLastLogin(jo.getLong(ModelKeys.LAST_LOGIN, 0L))
            .setLoginCount(jo.getLong(ModelKeys.LOGIN_COUNT, 0L))
            .setMcc(jo.getInteger(ModelKeys.MCC, 0))
            .setNickname(jo.getString(ModelKeys.NICKNAME, ""))
            .setOs(jo.getString(ModelKeys.OS, ""))
            .setOsLocale(jo.getString(ModelKeys.OS_LOCALE, ""))
            .setPhone(jo.getString(ModelKeys.PHONE, ""))
            .setPhoneVerified(jo.getBoolean(ModelKeys.PHONE_VERIFIED, false))
            .setPremiumEnd(jo.getLong(ModelKeys.PREMIUM_END, 0L))
            .setPremiumExp(jo.getLong(ModelKeys.PREMIUM_EXP, 0L))
            .setPremiumLevel(jo.getInteger(ModelKeys.PREMIUM_LEVEL, 0))
            .setSecret(jo.getString(ModelKeys.SECRET, ""))
            .setSince(jo.getLong(ModelKeys.SINCE, 0L))
            .setSocialId(jo.getString(ModelKeys.SOCIAL_ID, ""))
            .setSocialName(jo.getString(ModelKeys.SOCIAL_NAME, ""))
            .setSocialVerified(jo.getBoolean(ModelKeys.SOCIAL_VERIFIED, false))
            .setStatus(jo.getInteger(ModelKeys.STATUS, 0))
            .setTimezone(jo.getInteger(ModelKeys.TIMEZONE, 0))
            .setToken(jo.getString(ModelKeys.TOKEN, ""))
            .build();
    }

    public static JsonObject userEntry2Json(Db.DB.UserEntry userEntry) {
        return new JsonObject()
            .put(ModelKeys.USN, userEntry.getUsn())
            .put(ModelKeys.UID, userEntry.getUid())
            .put(ModelKeys.APP_LANGUAGE, userEntry.getAppLanguage())
            .put(ModelKeys.APP_VERSION, userEntry.getAppVersion())
            .put(ModelKeys.AVATAR, userEntry.getAvatar())
            .put(ModelKeys.BIRTHDAY, userEntry.getBirthday())
            .put(ModelKeys.COUNTRY, userEntry.getCountry())
            .put(ModelKeys.DEVICE_TYPE, userEntry.getDeviceType())
            .put(ModelKeys.EMAIL, userEntry.getEmail())
            .put(ModelKeys.EMAIL_VERIFIED, userEntry.getEmailVerified())
            .put(ModelKeys.GENDER, userEntry.getGender())
            .put(ModelKeys.IS_OFFICIAL, userEntry.getIsOfficial())
            .put(ModelKeys.IS_ROBOT, userEntry.getIsRobot())
            .put(ModelKeys.LAST_IP, userEntry.getLastIp())
            .put(ModelKeys.LAST_LOGIN, userEntry.getLastLogin())
            .put(ModelKeys.LOGIN_COUNT, userEntry.getLoginCount())
            .put(ModelKeys.MCC, userEntry.getMcc())
            .put(ModelKeys.NICKNAME, userEntry.getNickname())
            .put(ModelKeys.OS, userEntry.getOs())
            .put(ModelKeys.OS_LOCALE, userEntry.getOsLocale())
            .put(ModelKeys.PHONE, userEntry.getPhone())
            .put(ModelKeys.PHONE_VERIFIED, userEntry.getPhoneVerified())
            .put(ModelKeys.PREMIUM_END, userEntry.getPremiumEnd())
            .put(ModelKeys.PREMIUM_EXP, userEntry.getPremiumExp())
            .put(ModelKeys.PREMIUM_LEVEL, userEntry.getPremiumLevel())
            .put(ModelKeys.SECRET, userEntry.getSecret())
            .put(ModelKeys.SINCE, userEntry.getSince())
            .put(ModelKeys.SOCIAL_ID, userEntry.getSocialId())
            .put(ModelKeys.SOCIAL_NAME, userEntry.getSocialName())
            .put(ModelKeys.SOCIAL_VERIFIED, userEntry.getSocialVerified())
            .put(ModelKeys.STATUS, userEntry.getStatus())
            .put(ModelKeys.TIMEZONE, userEntry.getTimezone())
            .put(ModelKeys.TOKEN, userEntry.getToken());
    }

    private static Set<String> getUserMapKeySet() {
        Set<String> keySet = new HashSet<>();
        keySet.add(ModelKeys.USN);
        keySet.add(ModelKeys.UID);
        keySet.add(ModelKeys.APP_LANGUAGE);
        keySet.add(ModelKeys.APP_VERSION);
        keySet.add(ModelKeys.AVATAR);
        keySet.add(ModelKeys.BIRTHDAY);
        keySet.add(ModelKeys.COUNTRY);
        keySet.add(ModelKeys.DEVICE_TYPE);
        keySet.add(ModelKeys.EMAIL);
        keySet.add(ModelKeys.EMAIL_VERIFIED);
        keySet.add(ModelKeys.GENDER);
        keySet.add(ModelKeys.IS_OFFICIAL);
        keySet.add(ModelKeys.IS_ROBOT);
        keySet.add(ModelKeys.LAST_IP);
        keySet.add(ModelKeys.LAST_LOGIN);
        keySet.add(ModelKeys.LOGIN_COUNT);
        keySet.add(ModelKeys.MCC);
        keySet.add(ModelKeys.NICKNAME);
        keySet.add(ModelKeys.OS);
        keySet.add(ModelKeys.OS_LOCALE);
        keySet.add(ModelKeys.PHONE);
        keySet.add(ModelKeys.PHONE_VERIFIED);
        keySet.add(ModelKeys.PREMIUM_END);
        keySet.add(ModelKeys.PREMIUM_EXP);
        keySet.add(ModelKeys.PREMIUM_LEVEL);
        keySet.add(ModelKeys.SECRET);
        keySet.add(ModelKeys.SINCE);
        keySet.add(ModelKeys.SOCIAL_ID);
        keySet.add(ModelKeys.SOCIAL_NAME);
        keySet.add(ModelKeys.SOCIAL_VERIFIED);
        keySet.add(ModelKeys.STATUS);
        keySet.add(ModelKeys.TIMEZONE);
        keySet.add(ModelKeys.TOKEN);
        return keySet;
    }

    public static JsonObject removeKeysForLoginResponse(JsonObject jsonObject) {
        jsonObject.remove(ModelKeys.APP_LANGUAGE);
        jsonObject.remove(ModelKeys.APP_VERSION);
        jsonObject.remove(ModelKeys.DEVICE_TYPE);
        jsonObject.remove(ModelKeys.MCC);
        jsonObject.remove(ModelKeys.OS);
        jsonObject.remove(ModelKeys.OS_LOCALE);
        jsonObject.remove(ModelKeys.SECRET);
        jsonObject.remove(ModelKeys.SOCIAL_ID);
        jsonObject.remove(ModelKeys.SOCIAL_NAME);
        jsonObject.remove(ModelKeys.TIMEZONE);

        return jsonObject;
    }

    public static JsonObject removeKeysForRegisterRequest(JsonObject jsonObject) {
        jsonObject.remove(ModelKeys.USN);
        jsonObject.remove(ModelKeys.UID);
        jsonObject.remove(ModelKeys.EMAIL_VERIFIED);
        jsonObject.remove(ModelKeys.IS_OFFICIAL);
        jsonObject.remove(ModelKeys.IS_ROBOT);
        jsonObject.remove(ModelKeys.PHONE_VERIFIED);
        jsonObject.remove(ModelKeys.PREMIUM_END);
        jsonObject.remove(ModelKeys.PREMIUM_EXP);
        jsonObject.remove(ModelKeys.PREMIUM_LEVEL);
        jsonObject.remove(ModelKeys.SECRET);
        jsonObject.remove(ModelKeys.SINCE);
        jsonObject.remove(ModelKeys.SOCIAL_VERIFIED);
        jsonObject.remove(ModelKeys.TIMEZONE);

        return jsonObject;
    }

    public static Set<String> getUserUpdatableMapKeySet() {
        Set<String> keySet = new HashSet<>();
        keySet.add(ModelKeys.APP_LANGUAGE);
        keySet.add(ModelKeys.APP_VERSION);
        keySet.add(ModelKeys.AVATAR);
        keySet.add(ModelKeys.BIRTHDAY);
        keySet.add(ModelKeys.COUNTRY);
        keySet.add(ModelKeys.DEVICE_TYPE);
        keySet.add(ModelKeys.EMAIL);
        keySet.add(ModelKeys.EMAIL_VERIFIED);
        keySet.add(ModelKeys.GENDER);
        keySet.add(ModelKeys.LAST_IP);
        keySet.add(ModelKeys.LAST_LOGIN);
        keySet.add(ModelKeys.LOGIN_COUNT);
        keySet.add(ModelKeys.MCC);
        keySet.add(ModelKeys.NICKNAME);
        keySet.add(ModelKeys.OS);
        keySet.add(ModelKeys.OS_LOCALE);
        keySet.add(ModelKeys.PHONE);
        keySet.add(ModelKeys.PHONE_VERIFIED);
        keySet.add(ModelKeys.PREMIUM_END);
        keySet.add(ModelKeys.PREMIUM_EXP);
        keySet.add(ModelKeys.PREMIUM_LEVEL);
        keySet.add(ModelKeys.SECRET);
        keySet.add(ModelKeys.SOCIAL_ID);
        keySet.add(ModelKeys.SOCIAL_NAME);
        keySet.add(ModelKeys.SOCIAL_VERIFIED);
        keySet.add(ModelKeys.STATUS);
        keySet.add(ModelKeys.TIMEZONE);
        keySet.add(ModelKeys.TOKEN);
        return keySet;
    }
}
