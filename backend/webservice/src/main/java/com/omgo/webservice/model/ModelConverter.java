package com.omgo.webservice.model;

import com.omgo.utils.ModelKeys;
import io.vertx.core.json.JsonObject;
import proto.Db;

public class ModelConverter {

    public static Db.DB.UserEntry json2UserEntry(JsonObject jsonObject) {
        return Db.DB.UserEntry.newBuilder()
            .setUsn(jsonObject.getLong(ModelKeys.USN, 0L))
            .setUid(jsonObject.getLong(ModelKeys.UID, 0L))
            .setAppLanguage(jsonObject.getString(ModelKeys.APP_LANGUAGE, ""))
            .setAppVersion(jsonObject.getString(ModelKeys.APP_VERSION, ""))
            .setAvatar(jsonObject.getString(ModelKeys.AVATAR, ""))
            .setBirthday(jsonObject.getLong(ModelKeys.BIRTHDAY, 0L))
            .setCountry(jsonObject.getString(ModelKeys.COUNTRY, ""))
            .setDeviceType(jsonObject.getInteger(ModelKeys.DEVICE_TYPE, 0))
            .setEmail(jsonObject.getString(ModelKeys.EMAIL, ""))
            .setEmailVerified(jsonObject.getBoolean(ModelKeys.EMAIL_VERIFIED, false))
            .setGender(jsonObject.getInteger(ModelKeys.GENDER, 0))
            .setIsOfficial(jsonObject.getBoolean(ModelKeys.IS_OFFICIAL, false))
            .setIsRobot(jsonObject.getBoolean(ModelKeys.IS_ROBOT, false))
            .setLastIp(jsonObject.getString(ModelKeys.LAST_IP, ""))
            .setLastLogin(jsonObject.getLong(ModelKeys.LAST_LOGIN, 0L))
            .setLoginCount(jsonObject.getLong(ModelKeys.LOGIN_COUNT, 0L))
            .setMcc(jsonObject.getInteger(ModelKeys.MCC, 0))
            .setNickname(jsonObject.getString(ModelKeys.NICKNAME, ""))
            .setOs(jsonObject.getString(ModelKeys.OS, ""))
            .setOsLocale(jsonObject.getString(ModelKeys.OS_LOCALE, ""))
            .setPhone(jsonObject.getString(ModelKeys.PHONE, ""))
            .setPhoneVerified(jsonObject.getBoolean(ModelKeys.PHONE_VERIFIED, false))
            .setPremiumEnd(jsonObject.getLong(ModelKeys.PREMIUM_END, 0L))
            .setPremiumExp(jsonObject.getLong(ModelKeys.PREMIUM_EXP, 0L))
            .setPremiumLevel(jsonObject.getInteger(ModelKeys.PREMIUM_LEVEL, 0))
            .setSecret(jsonObject.getString(ModelKeys.SECRET, ""))
            .setSince(jsonObject.getLong(ModelKeys.SINCE, 0L))
            .setSocialId(jsonObject.getString(ModelKeys.SOCIAL_ID, ""))
            .setSocialName(jsonObject.getString(ModelKeys.SOCIAL_NAME, ""))
            .setSocialVerified(jsonObject.getBoolean(ModelKeys.SOCIAL_VERIFIED, false))
            .setStatus(jsonObject.getInteger(ModelKeys.STATUS, 0))
            .setTimezone(jsonObject.getInteger(ModelKeys.TIMEZONE, 0))
            .setToken(jsonObject.getString(ModelKeys.TOKEN, ""))
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
}
