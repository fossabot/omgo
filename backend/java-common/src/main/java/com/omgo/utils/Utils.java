package com.omgo.utils;

import io.vertx.core.json.DecodeException;
import io.vertx.core.json.JsonObject;
import org.apache.commons.validator.routines.EmailValidator;

import java.security.MessageDigest;
import java.util.Base64;

public class Utils {
    public static boolean DEBUG = false;

    public static boolean STANDALONE = false;

    /**
     * check if s is an empty string
     *
     * @param s
     * @return
     */
    public static boolean isEmptyString(String s) {
        return s == null || s.equals("");
    }

    /**
     * check if s is an empty string or not
     *
     * @param s
     * @return
     */
    public static boolean isNotEmptyString(String s) {
        return !isEmptyString(s);
    }

    /**
     * parse string to json object safely
     *
     * @param s
     * @return
     */
    public static JsonObject safeParseJson(String s) {
        try {
            return new JsonObject(s);
        } catch (DecodeException e) {
            e.printStackTrace();
        }

        return null;
    }

    /**
     * encode byte array to base64 string
     *
     * @param raw
     * @return
     */
    public static String encodeBase64(byte[] raw) {
        return Base64.getEncoder().encodeToString(raw);
    }

    /**
     * decode base64 encoded string to byte array
     *
     * @param encoded
     * @return
     */
    public static byte[] decodeBase64(String encoded) {
        return Base64.getDecoder().decode(encoded);
    }

    /**
     * calculate sha1 signature of string
     *
     * @param raw
     * @return
     */
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

    /**
     * check if email address is valid
     *
     * @param email
     * @return
     */
    public static boolean isValidEmailAddress(String email) {
        return EmailValidator.getInstance().isValid(email);
    }
}
