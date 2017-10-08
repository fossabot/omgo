package com.omgo.dataservice;


import com.omgo.dataservice.model.ModelConverter;
import com.omgo.dataservice.model.Utils;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;
import java.util.Random;

public final class AccountUtils {
    public static final int PASSWORD_MIN_LEN = 6;

    private static final String STRING_EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$";

    private static Random random = new Random(System.currentTimeMillis());

    public static String getRedisKey(long usn) {
        return String.format("%s:%d", ModelConverter.KEY_USER, usn);
    }

    public static byte[] getToken() {
        byte[] raw = new byte[32];
        random.nextBytes(raw);
        try {
            MessageDigest digestMD5 = MessageDigest.getInstance("MD5");
            return digestMD5.digest(raw);
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        }

        return null;
    }

    public static String encodeBase64(byte[] raw) {
        return Base64.getEncoder().encodeToString(raw);
    }

    public static byte[] decodeBase64(String encoded) {
        return Base64.getDecoder().decode(encoded);
    }

    public static String saltedSecret(String secret, long salt) {
        try {
            String salted = secret + String.valueOf(salt);
            MessageDigest digestSHA1 = MessageDigest.getInstance("SHA-1");
            digestSHA1.reset();
            byte[] saltedRaw = digestSHA1.digest(salted.getBytes("UTF-8"));
            return encodeBase64(saltedRaw);
        } catch (Exception e) {
            e.printStackTrace();
        }

        return null;
    }

    // FIXME: 29/09/2017 invalid email address like xxx@xxx can pass this test
    public static boolean isValidEmailAddress(String s) {
        if (Utils.isEmptyString(s)) {
            return false;
        } else {
            return s.matches(STRING_EMAIL_REGEX);
        }
    }

    public static boolean isValidSecret(String s) {
        if (Utils.isEmptyString(s)) {
            return false;
        } else {
            return s.length() >= PASSWORD_MIN_LEN;
        }
    }

    public static int nextUsnIncrement() {
        return random.nextInt(20000) * 2 + 1;
    }

    public static int nextUidIncrement() {
        return random.nextInt(10000);
    }
}
