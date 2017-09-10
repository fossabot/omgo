package com.omgo.webservice;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import io.vertx.core.AsyncResult;
import io.vertx.core.Future;
import io.vertx.core.Handler;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.auth.AbstractUser;
import io.vertx.ext.auth.AuthProvider;
import proto.Db;

import java.nio.charset.StandardCharsets;

// https://github.com/vert-x3/vertx-auth/blob/master/vertx-auth-jdbc/src/main/java/io/vertx/ext/auth/jdbc/impl/JDBCUser.java
// https://github.com/vert-x3/vertx-auth/blob/master/vertx-auth-jwt/src/main/java/io/vertx/ext/auth/jwt/impl/JWTUser.java

public class GrpcAuthUser extends AbstractUser {

    private GRPCAuthProvider authProvider;
    private String email;
    private JsonObject principle;
    private Db.DB.UserExtendInfo userExtendInfo;

    public GrpcAuthUser(GRPCAuthProvider authProvider, String email, Db.DB.UserExtendInfo userExtendInfo) {
        this.authProvider = authProvider;
        this.email = email;
        this.userExtendInfo = userExtendInfo;
    }

    @Override
    protected void doIsPermitted(String s, Handler<AsyncResult<Boolean>> handler) {
        handler.handle(Future.succeededFuture(true));
    }

    @Override
    public JsonObject principal() {
        if (principle == null) {
//            principle = new JsonObject().put("email", email);
            principle = new JsonObject();
            try {
                String info = JsonFormat.printer().print(userExtendInfo);
                principle = new JsonObject(info);
                return principle;
            } catch (InvalidProtocolBufferException e) {
                e.printStackTrace();
            }
        }
        return principle;
    }

    @Override
    public void setAuthProvider(AuthProvider authProvider) {
        if (authProvider instanceof GRPCAuthProvider) {
            this.authProvider = (GRPCAuthProvider) authProvider;
        } else {
            throw new IllegalArgumentException("Not a GrpcAuthProvider");
        }
    }

    @Override
    public void writeToBuffer(Buffer buffer) {
        super.writeToBuffer(buffer);
        byte[] bytes = email.getBytes(StandardCharsets.UTF_8);
        buffer.appendInt(bytes.length);
        buffer.appendBytes(bytes);

        bytes = userExtendInfo.toByteArray();
        buffer.appendInt(bytes.length);
        buffer.appendBytes(bytes);
    }

    @Override
    public int readFromBuffer(int pos, Buffer buffer) {
        pos = super.readFromBuffer(pos, buffer);
        int len = buffer.getInt(pos);
        pos += 4;
        byte[] bytes = buffer.getBytes(pos, pos + len);
        email = new String(bytes, StandardCharsets.UTF_8);
        pos += len;

        len = buffer.getInt(pos);
        pos += 4;
        bytes = buffer.getBytes(pos, pos + len);
        pos += len;
        try {
            userExtendInfo = Db.DB.UserExtendInfo.parseFrom(bytes);
        } catch (InvalidProtocolBufferException e) {
            e.printStackTrace();
        }

        return pos;
    }
}
