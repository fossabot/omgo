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

    public static String saltedSecret(String secret, String salt) {
        try {
            String salted = secret + salt;
            MessageDigest digestSHA1 = MessageDigest.getInstance("SHA-1");
            digestSHA1.reset();
            byte[] saltedRaw = digestSHA1.digest(salted.getBytes("UTF-8"));
            return base64(saltedRaw);
        } catch (Exception e) {
            e.printStackTrace();
        }

        return null;
    }
}
