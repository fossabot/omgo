package com.omgo.dbservice;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;
import java.util.Random;

public final class AccountUtils {
    public static final int PASSWORD_MIN_LEN = 6;


    private static Random random = new Random(System.currentTimeMillis());

    public static byte[] getSalt() {
        byte[] salt = new byte[32];
        random.nextBytes(salt);
        return salt;
    }

    public static byte[] getToken(byte[] salt) {
        byte[] raw = new byte[32 + salt.length];
        random.nextBytes(raw);
        System.arraycopy(salt, 0, raw, 32, salt.length);
        try {
            MessageDigest digestMD5 = MessageDigest.getInstance("MD5");
            return digestMD5.digest(raw);
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        }

        return null;
    }

    public static String base64(byte[] raw) {
        return Base64.getEncoder().encodeToString(raw);
    }
}
