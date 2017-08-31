package com.omgo.dbservice;

import com.omgo.dbservice.model.Utils;
import proto.Db.DB;
import proto.common.Common;

public class DbProtoUtils {

    public static DB.Result makeOkResult() {
        return makeResult(DB.StatusCode.STATUS_OK, "");
    }

    public static DB.Result makeResult(DB.StatusCode code) {
        return makeResult(code, "");
    }

    public static DB.Result makeResult(DB.StatusCode code, String msg) {
        DB.Result.Builder builder = DB.Result.newBuilder();
        builder.setStatus(code);
        if (!Utils.isEmptyString(msg)) {
            builder.setMsg(msg);
        }
        return builder.build();
    }

    public static DB.UserExtendInfo makeUserExtendInfo(Common.UserInfo userInfo) {
        return makeUserExtendInfo(userInfo, "", "");
    }

    public static DB.UserExtendInfo makeUserExtendInfo(Common.UserInfo userInfo, String secret, String token) {
        DB.UserExtendInfo.Builder builder = DB.UserExtendInfo.newBuilder();
        builder.setInfo(userInfo);
        if (!Utils.isEmptyString(secret)) {
            builder.setSecret(secret);
        }
        if (!Utils.isEmptyString(token)) {
            builder.setToken(token);
        }

        return builder.build();
    }

    public static DB.UserOpResult makeUserOpOkResult(DB.UserExtendInfo extendInfo) {
        return makeUserOpResult(makeOkResult(), extendInfo);
    }

    public static DB.UserOpResult makeUserOpOkResult(Common.UserInfo userInfo) {
        return makeUserOpResult(makeOkResult(), makeUserExtendInfo(userInfo, "", ""));
    }

    public static DB.UserOpResult makeUserOpResult(DB.StatusCode code, String msg) {
        return makeUserOpResult(makeResult(code, msg), null);
    }

    public static DB.UserOpResult makeUserOpResult(DB.Result result, DB.UserExtendInfo extendInfo) {
        DB.UserOpResult.Builder builder = DB.UserOpResult.newBuilder();
        builder.setResult(result);
        if (extendInfo != null) {
            builder.setUserExtInfo(extendInfo);
        }

        return builder.build();
    }
}
