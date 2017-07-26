package com.omgo.dbservice.driver;

import com.omgo.dbservice.model.ModelConstant;

public class Utils {
    public static String getRedisKey(long usn) {
        return String.format("%s:%d", ModelConstant.KEY_USER, usn);
    }
}
