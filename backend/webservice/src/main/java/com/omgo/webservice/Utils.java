package com.omgo.webservice;

import java.util.Base64;

public class Utils {
    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }

    public static String encodeBase64(byte[] raw) {
        return Base64.getEncoder().encodeToString(raw);
    }

    public static byte[] decodeBase64(String encoded) {
        return Base64.getDecoder().decode(encoded);
    }
}
