package com.omgo.webservice;

public class ApiConstant {
    public static final String API_ROOT = "/api/";
    public static final String API_REGISTER = "register";
    public static final String API_LOGIN = "login";

    public static String getApiPath(String api) {
        return API_ROOT + api;
    }
}
