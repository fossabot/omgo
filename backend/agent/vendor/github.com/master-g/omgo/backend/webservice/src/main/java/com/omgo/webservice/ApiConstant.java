package com.omgo.webservice;

public class ApiConstant {
    public static final String API_ROOT = "/api/";
    public static final String API_REGISTER = "register";
    public static final String API_LOGIN = "login";
    public static final String API_HANDSHAKE = "handshake";
    public static final String API_USERPROFILE = "userprofile";

    public static final String API_TEST = "test";

    public static String getApiPath(String api) {
        return API_ROOT + api;
    }
}
