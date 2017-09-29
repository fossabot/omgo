package com.omgo.webservice;

import java.security.MessageDigest;
import java.util.Base64;

public class Utils {
    public static boolean DEBUG = false;

    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }

    public static String encodeBase64(byte[] raw) {
        return Base64.getEncoder().encodeToString(raw);
    }

    public static byte[] decodeBase64(String encoded) {
        return Base64.getDecoder().decode(encoded);
    }

    public static byte[] sha1(String raw) {
        try {
            MessageDigest digestSHA1 = MessageDigest.getInstance("SHA-1");
            digestSHA1.reset();
            return digestSHA1.digest(raw.getBytes("UTF-8"));
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }
}
