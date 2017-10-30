package com.omgo.dataservice.model;

import io.vertx.core.json.DecodeException;
import io.vertx.core.json.JsonObject;

public class Utils {

    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }

    public static boolean isNotEmptyString(String s) {
        return !isEmptyString(s);
    }

    public static JsonObject safeParseJson(String s) {
        try {
            return new JsonObject(s);
        } catch (DecodeException e) {
            e.printStackTrace();
        }

        return null;
    }
}
