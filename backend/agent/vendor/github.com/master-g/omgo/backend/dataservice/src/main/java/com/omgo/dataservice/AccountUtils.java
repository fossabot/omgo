package com.omgo.dataservice;


import com.omgo.utils.ModelKeys;
import com.omgo.utils.Utils;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Random;

public final class AccountUtils {
    public static final int PASSWORD_MIN_LEN = 6;

    private static Random random = new Random(System.currentTimeMillis());

    public static String getRedisKey(long usn) {
        return String.format("%s:%d", ModelKeys.USER, usn);
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

    public static String saltedSecret(String secret, long salt) {
        try {
            String salted = secret + String.valueOf(salt);
            MessageDigest digestSHA1 = MessageDigest.getInstance("SHA-1");
            digestSHA1.reset();
            byte[] saltedRaw = digestSHA1.digest(salted.getBytes("UTF-8"));
            return Utils.encodeBase64(saltedRaw);
        } catch (Exception e) {
            e.printStackTrace();
        }

        return null;
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
