package com.omgo.dbservice.model;

public class Utils {

    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }

    public static boolean isNotEmptyString(String s) {
        return !isEmptyString(s);
    }
}
