package com.omgo.dbservice.model;

public class SQLConstant {
    public static final String QUERY_WITH_USN = "SELECT * FROM user WHERE usn=?";
    public static final String QUERY_WITH_UID = "SELECT * FROM user WHERE uid=?";
    public static final String QUERY_WITH_EMAIL = "SELECT * FROM user WHERE email=?";
}
