package com.omgo.dataservice;

import com.omgo.dataservice.model.ModelConverter;
import com.omgo.dataservice.model.Utils;
import io.vertx.core.json.JsonObject;
import proto.Db.DB;

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
        if (Utils.isNotEmptyString(msg)) {
            builder.setMsg(msg);
        }
        return builder.build();
    }

    public static DB.UserOpResult makeUserOpOkResult(JsonObject jsonObject) {
        return DB.UserOpResult.newBuilder()
            .setResult(makeOkResult())
            .setUser(ModelConverter.json2UserEntry(jsonObject))
            .build();
    }

    public static DB.UserOpResult makeUserOpResult(DB.StatusCode code) {
        return DB.UserOpResult.newBuilder()
            .setResult(makeResult(code))
            .build();
    }

    public static DB.UserOpResult makeUserOpResult(DB.StatusCode code, String msg) {
        return DB.UserOpResult.newBuilder()
            .setResult(makeResult(code, msg))
            .build();
    }

    public static DB.UserOpResult makeUserOpInternalFailedResult(String msg) {
        return DB.UserOpResult.newBuilder()
            .setResult(makeResult(DB.StatusCode.STATUS_INTERNAL_ERROR, msg))
            .build();
    }
}
