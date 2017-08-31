package com.omgo.dbservice.model;

import com.omgo.dbservice.model.ModelConstant;

import java.util.Objects;

public class Utils {

    public static String getRedisKey(long usn) {
        return String.format("%s:%d", ModelConstant.KEY_USER, usn);
    }

    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }
}
