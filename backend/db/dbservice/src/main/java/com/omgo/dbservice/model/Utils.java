package com.omgo.dbservice.model;

import com.omgo.dbservice.model.ModelConstant;

import java.util.Objects;

public class Utils {
    private static final String STRING_EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$";


    public static String getRedisKey(long usn) {
        return String.format("%s:%d", ModelConstant.KEY_USER, usn);
    }

    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }

    public static boolean isValidEmailAddress(String s) {
        if (isEmptyString(s)) {
            return false;
        } else {
            return s.matches(STRING_EMAIL_REGEX);
        }
    }
}
