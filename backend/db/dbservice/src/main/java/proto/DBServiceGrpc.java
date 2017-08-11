package proto;

import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;

/**
 * <pre>
 * Database service definition
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.3.0)",
    comments = "Source: db.proto")
public final class DBServiceGrpc {

  private DBServiceGrpc() {}

  private static <T> io.grpc.stub.StreamObserver<T> toObserver(final io.vertx.core.Handler<io.vertx.core.AsyncResult<T>> handler) {
    return new io.grpc.stub.StreamObserver<T>() {
      private volatile boolean resolved = false;
      @Override
      public void onNext(T value) {
        if (!resolved) {
          resolved = true;
          handler.handle(io.vertx.core.Future.succeededFuture(value));
        }
      }

      @Override
      public void onError(Throwable t) {
        if (!resolved) {
          resolved = true;
          handler.handle(io.vertx.core.Future.failedFuture(t));
        }
      }

      @Override
      public void onCompleted() {
        if (!resolved) {
          resolved = true;
          handler.handle(io.vertx.core.Future.succeededFuture());
        }
      }
    };
  }

  public static final String SERVICE_NAME = "proto.DBService";

  // Static method descriptors that strictly reflect the proto.
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.Db.DB.UserKey,
      proto.common.Common.UserInfo> METHOD_USER_QUERY =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.DBService", "UserQuery"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserKey.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.common.Common.UserInfo.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.common.Common.UserInfo,
      proto.Db.DB.NullValue> METHOD_USER_UPDATE_INFO =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.DBService", "UserUpdateInfo"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.common.Common.UserInfo.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.NullValue.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.Db.DB.UserExtendInfo,
      proto.Db.DB.UserExtendInfo> METHOD_USER_REGISTER =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.DBService", "UserRegister"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserExtendInfo.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserExtendInfo.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.Db.DB.UserExtendInfo,
      proto.Db.DB.UserExtendInfo> METHOD_USER_LOGIN =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.DBService", "UserLogin"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserExtendInfo.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserExtendInfo.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.Db.DB.UserLogoutRequest,
      proto.Db.DB.NullValue> METHOD_USER_LOGOUT =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.DBService", "UserLogout"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserLogoutRequest.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.NullValue.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.Db.DB.UserKey,
      proto.Db.DB.UserExtendInfo> METHOD_USER_EXTRA_INFO_QUERY =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.DBService", "UserExtraInfoQuery"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserKey.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.Db.DB.UserExtendInfo.getDefaultInstance()));

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static DBServiceStub newStub(io.grpc.Channel channel) {
    return new DBServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static DBServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new DBServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary and streaming output calls on the service
   */
  public static DBServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new DBServiceFutureStub(channel);
  }

  /**
   * Creates a new vertx stub that supports all call types for the service
   */
  public static DBServiceVertxStub newVertxStub(io.grpc.Channel channel) {
    return new DBServiceVertxStub(channel);
  }

  /**
   * <pre>
   * Database service definition
   * </pre>
   */
  public static abstract class DBServiceImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * query user info
     * </pre>
     */
    public void userQuery(proto.Db.DB.UserKey request,
        io.grpc.stub.StreamObserver<proto.common.Common.UserInfo> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_USER_QUERY, responseObserver);
    }

    /**
     * <pre>
     * update user info
     * </pre>
     */
    public void userUpdateInfo(proto.common.Common.UserInfo request,
        io.grpc.stub.StreamObserver<proto.Db.DB.NullValue> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_USER_UPDATE_INFO, responseObserver);
    }

    /**
     * <pre>
     * register
     * </pre>
     */
    public void userRegister(proto.Db.DB.UserExtendInfo request,
        io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_USER_REGISTER, responseObserver);
    }

    /**
     * <pre>
     * login
     * </pre>
     */
    public void userLogin(proto.Db.DB.UserExtendInfo request,
        io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_USER_LOGIN, responseObserver);
    }

    /**
     * <pre>
     * logout
     * </pre>
     */
    public void userLogout(proto.Db.DB.UserLogoutRequest request,
        io.grpc.stub.StreamObserver<proto.Db.DB.NullValue> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_USER_LOGOUT, responseObserver);
    }

    /**
     * <pre>
     * verify token
     * </pre>
     */
    public void userExtraInfoQuery(proto.Db.DB.UserKey request,
        io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_USER_EXTRA_INFO_QUERY, responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_USER_QUERY,
            asyncUnaryCall(
              new MethodHandlers<
                proto.Db.DB.UserKey,
                proto.common.Common.UserInfo>(
                  this, METHODID_USER_QUERY)))
          .addMethod(
            METHOD_USER_UPDATE_INFO,
            asyncUnaryCall(
              new MethodHandlers<
                proto.common.Common.UserInfo,
                proto.Db.DB.NullValue>(
                  this, METHODID_USER_UPDATE_INFO)))
          .addMethod(
            METHOD_USER_REGISTER,
            asyncUnaryCall(
              new MethodHandlers<
                proto.Db.DB.UserExtendInfo,
                proto.Db.DB.UserExtendInfo>(
                  this, METHODID_USER_REGISTER)))
          .addMethod(
            METHOD_USER_LOGIN,
            asyncUnaryCall(
              new MethodHandlers<
                proto.Db.DB.UserExtendInfo,
                proto.Db.DB.UserExtendInfo>(
                  this, METHODID_USER_LOGIN)))
          .addMethod(
            METHOD_USER_LOGOUT,
            asyncUnaryCall(
              new MethodHandlers<
                proto.Db.DB.UserLogoutRequest,
                proto.Db.DB.NullValue>(
                  this, METHODID_USER_LOGOUT)))
          .addMethod(
            METHOD_USER_EXTRA_INFO_QUERY,
            asyncUnaryCall(
              new MethodHandlers<
                proto.Db.DB.UserKey,
                proto.Db.DB.UserExtendInfo>(
                  this, METHODID_USER_EXTRA_INFO_QUERY)))
          .build();
    }
  }

  /**
   * <pre>
   * Database service definition
   * </pre>
   */
  public static final class DBServiceStub extends io.grpc.stub.AbstractStub<DBServiceStub> {
    private DBServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private DBServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DBServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new DBServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * query user info
     * </pre>
     */
    public void userQuery(proto.Db.DB.UserKey request,
        io.grpc.stub.StreamObserver<proto.common.Common.UserInfo> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_QUERY, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * update user info
     * </pre>
     */
    public void userUpdateInfo(proto.common.Common.UserInfo request,
        io.grpc.stub.StreamObserver<proto.Db.DB.NullValue> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_UPDATE_INFO, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * register
     * </pre>
     */
    public void userRegister(proto.Db.DB.UserExtendInfo request,
        io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_REGISTER, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * login
     * </pre>
     */
    public void userLogin(proto.Db.DB.UserExtendInfo request,
        io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_LOGIN, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * logout
     * </pre>
     */
    public void userLogout(proto.Db.DB.UserLogoutRequest request,
        io.grpc.stub.StreamObserver<proto.Db.DB.NullValue> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_LOGOUT, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * verify token
     * </pre>
     */
    public void userExtraInfoQuery(proto.Db.DB.UserKey request,
        io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_EXTRA_INFO_QUERY, getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Database service definition
   * </pre>
   */
  public static final class DBServiceBlockingStub extends io.grpc.stub.AbstractStub<DBServiceBlockingStub> {
    private DBServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private DBServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DBServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new DBServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * query user info
     * </pre>
     */
    public proto.common.Common.UserInfo userQuery(proto.Db.DB.UserKey request) {
      return blockingUnaryCall(
          getChannel(), METHOD_USER_QUERY, getCallOptions(), request);
    }

    /**
     * <pre>
     * update user info
     * </pre>
     */
    public proto.Db.DB.NullValue userUpdateInfo(proto.common.Common.UserInfo request) {
      return blockingUnaryCall(
          getChannel(), METHOD_USER_UPDATE_INFO, getCallOptions(), request);
    }

    /**
     * <pre>
     * register
     * </pre>
     */
    public proto.Db.DB.UserExtendInfo userRegister(proto.Db.DB.UserExtendInfo request) {
      return blockingUnaryCall(
          getChannel(), METHOD_USER_REGISTER, getCallOptions(), request);
    }

    /**
     * <pre>
     * login
     * </pre>
     */
    public proto.Db.DB.UserExtendInfo userLogin(proto.Db.DB.UserExtendInfo request) {
      return blockingUnaryCall(
          getChannel(), METHOD_USER_LOGIN, getCallOptions(), request);
    }

    /**
     * <pre>
     * logout
     * </pre>
     */
    public proto.Db.DB.NullValue userLogout(proto.Db.DB.UserLogoutRequest request) {
      return blockingUnaryCall(
          getChannel(), METHOD_USER_LOGOUT, getCallOptions(), request);
    }

    /**
     * <pre>
     * verify token
     * </pre>
     */
    public proto.Db.DB.UserExtendInfo userExtraInfoQuery(proto.Db.DB.UserKey request) {
      return blockingUnaryCall(
          getChannel(), METHOD_USER_EXTRA_INFO_QUERY, getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Database service definition
   * </pre>
   */
  public static final class DBServiceFutureStub extends io.grpc.stub.AbstractStub<DBServiceFutureStub> {
    private DBServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private DBServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DBServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new DBServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * query user info
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.common.Common.UserInfo> userQuery(
        proto.Db.DB.UserKey request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_USER_QUERY, getCallOptions()), request);
    }

    /**
     * <pre>
     * update user info
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.Db.DB.NullValue> userUpdateInfo(
        proto.common.Common.UserInfo request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_USER_UPDATE_INFO, getCallOptions()), request);
    }

    /**
     * <pre>
     * register
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.Db.DB.UserExtendInfo> userRegister(
        proto.Db.DB.UserExtendInfo request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_USER_REGISTER, getCallOptions()), request);
    }

    /**
     * <pre>
     * login
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.Db.DB.UserExtendInfo> userLogin(
        proto.Db.DB.UserExtendInfo request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_USER_LOGIN, getCallOptions()), request);
    }

    /**
     * <pre>
     * logout
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.Db.DB.NullValue> userLogout(
        proto.Db.DB.UserLogoutRequest request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_USER_LOGOUT, getCallOptions()), request);
    }

    /**
     * <pre>
     * verify token
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.Db.DB.UserExtendInfo> userExtraInfoQuery(
        proto.Db.DB.UserKey request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_USER_EXTRA_INFO_QUERY, getCallOptions()), request);
    }
  }

  /**
   * <pre>
   * Database service definition
   * </pre>
   */
  public static abstract class DBServiceVertxImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * query user info
     * </pre>
     */
    public void userQuery(proto.Db.DB.UserKey request,
        io.vertx.core.Future<proto.common.Common.UserInfo> response) {
      asyncUnimplementedUnaryCall(METHOD_USER_QUERY, DBServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * update user info
     * </pre>
     */
    public void userUpdateInfo(proto.common.Common.UserInfo request,
        io.vertx.core.Future<proto.Db.DB.NullValue> response) {
      asyncUnimplementedUnaryCall(METHOD_USER_UPDATE_INFO, DBServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * register
     * </pre>
     */
    public void userRegister(proto.Db.DB.UserExtendInfo request,
        io.vertx.core.Future<proto.Db.DB.UserExtendInfo> response) {
      asyncUnimplementedUnaryCall(METHOD_USER_REGISTER, DBServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * login
     * </pre>
     */
    public void userLogin(proto.Db.DB.UserExtendInfo request,
        io.vertx.core.Future<proto.Db.DB.UserExtendInfo> response) {
      asyncUnimplementedUnaryCall(METHOD_USER_LOGIN, DBServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * logout
     * </pre>
     */
    public void userLogout(proto.Db.DB.UserLogoutRequest request,
        io.vertx.core.Future<proto.Db.DB.NullValue> response) {
      asyncUnimplementedUnaryCall(METHOD_USER_LOGOUT, DBServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * verify token
     * </pre>
     */
    public void userExtraInfoQuery(proto.Db.DB.UserKey request,
        io.vertx.core.Future<proto.Db.DB.UserExtendInfo> response) {
      asyncUnimplementedUnaryCall(METHOD_USER_EXTRA_INFO_QUERY, DBServiceGrpc.toObserver(response.completer()));
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_USER_QUERY,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.Db.DB.UserKey,
                proto.common.Common.UserInfo>(
                  this, METHODID_USER_QUERY)))
          .addMethod(
            METHOD_USER_UPDATE_INFO,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.common.Common.UserInfo,
                proto.Db.DB.NullValue>(
                  this, METHODID_USER_UPDATE_INFO)))
          .addMethod(
            METHOD_USER_REGISTER,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.Db.DB.UserExtendInfo,
                proto.Db.DB.UserExtendInfo>(
                  this, METHODID_USER_REGISTER)))
          .addMethod(
            METHOD_USER_LOGIN,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.Db.DB.UserExtendInfo,
                proto.Db.DB.UserExtendInfo>(
                  this, METHODID_USER_LOGIN)))
          .addMethod(
            METHOD_USER_LOGOUT,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.Db.DB.UserLogoutRequest,
                proto.Db.DB.NullValue>(
                  this, METHODID_USER_LOGOUT)))
          .addMethod(
            METHOD_USER_EXTRA_INFO_QUERY,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.Db.DB.UserKey,
                proto.Db.DB.UserExtendInfo>(
                  this, METHODID_USER_EXTRA_INFO_QUERY)))
          .build();
    }
  }

  /**
   * <pre>
   * Database service definition
   * </pre>
   */
  public static final class DBServiceVertxStub extends io.grpc.stub.AbstractStub<DBServiceVertxStub> {
    private DBServiceVertxStub(io.grpc.Channel channel) {
      super(channel);
    }

    private DBServiceVertxStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DBServiceVertxStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new DBServiceVertxStub(channel, callOptions);
    }

    /**
     * <pre>
     * query user info
     * </pre>
     */
    public void userQuery(proto.Db.DB.UserKey request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.common.Common.UserInfo>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_QUERY, getCallOptions()), request, DBServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * update user info
     * </pre>
     */
    public void userUpdateInfo(proto.common.Common.UserInfo request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.Db.DB.NullValue>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_UPDATE_INFO, getCallOptions()), request, DBServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * register
     * </pre>
     */
    public void userRegister(proto.Db.DB.UserExtendInfo request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.Db.DB.UserExtendInfo>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_REGISTER, getCallOptions()), request, DBServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * login
     * </pre>
     */
    public void userLogin(proto.Db.DB.UserExtendInfo request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.Db.DB.UserExtendInfo>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_LOGIN, getCallOptions()), request, DBServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * logout
     * </pre>
     */
    public void userLogout(proto.Db.DB.UserLogoutRequest request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.Db.DB.NullValue>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_LOGOUT, getCallOptions()), request, DBServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * verify token
     * </pre>
     */
    public void userExtraInfoQuery(proto.Db.DB.UserKey request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.Db.DB.UserExtendInfo>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_USER_EXTRA_INFO_QUERY, getCallOptions()), request, DBServiceGrpc.toObserver(response));
    }
  }

  private static final int METHODID_USER_QUERY = 0;
  private static final int METHODID_USER_UPDATE_INFO = 1;
  private static final int METHODID_USER_REGISTER = 2;
  private static final int METHODID_USER_LOGIN = 3;
  private static final int METHODID_USER_LOGOUT = 4;
  private static final int METHODID_USER_EXTRA_INFO_QUERY = 5;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final DBServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(DBServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_USER_QUERY:
          serviceImpl.userQuery((proto.Db.DB.UserKey) request,
              (io.grpc.stub.StreamObserver<proto.common.Common.UserInfo>) responseObserver);
          break;
        case METHODID_USER_UPDATE_INFO:
          serviceImpl.userUpdateInfo((proto.common.Common.UserInfo) request,
              (io.grpc.stub.StreamObserver<proto.Db.DB.NullValue>) responseObserver);
          break;
        case METHODID_USER_REGISTER:
          serviceImpl.userRegister((proto.Db.DB.UserExtendInfo) request,
              (io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo>) responseObserver);
          break;
        case METHODID_USER_LOGIN:
          serviceImpl.userLogin((proto.Db.DB.UserExtendInfo) request,
              (io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo>) responseObserver);
          break;
        case METHODID_USER_LOGOUT:
          serviceImpl.userLogout((proto.Db.DB.UserLogoutRequest) request,
              (io.grpc.stub.StreamObserver<proto.Db.DB.NullValue>) responseObserver);
          break;
        case METHODID_USER_EXTRA_INFO_QUERY:
          serviceImpl.userExtraInfoQuery((proto.Db.DB.UserKey) request,
              (io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static final class VertxMethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final DBServiceVertxImplBase serviceImpl;
    private final int methodId;

    VertxMethodHandlers(DBServiceVertxImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_USER_QUERY:
          serviceImpl.userQuery((proto.Db.DB.UserKey) request,
              (io.vertx.core.Future<proto.common.Common.UserInfo>) io.vertx.core.Future.<proto.common.Common.UserInfo>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.common.Common.UserInfo>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_USER_UPDATE_INFO:
          serviceImpl.userUpdateInfo((proto.common.Common.UserInfo) request,
              (io.vertx.core.Future<proto.Db.DB.NullValue>) io.vertx.core.Future.<proto.Db.DB.NullValue>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.Db.DB.NullValue>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_USER_REGISTER:
          serviceImpl.userRegister((proto.Db.DB.UserExtendInfo) request,
              (io.vertx.core.Future<proto.Db.DB.UserExtendInfo>) io.vertx.core.Future.<proto.Db.DB.UserExtendInfo>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_USER_LOGIN:
          serviceImpl.userLogin((proto.Db.DB.UserExtendInfo) request,
              (io.vertx.core.Future<proto.Db.DB.UserExtendInfo>) io.vertx.core.Future.<proto.Db.DB.UserExtendInfo>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_USER_LOGOUT:
          serviceImpl.userLogout((proto.Db.DB.UserLogoutRequest) request,
              (io.vertx.core.Future<proto.Db.DB.NullValue>) io.vertx.core.Future.<proto.Db.DB.NullValue>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.Db.DB.NullValue>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_USER_EXTRA_INFO_QUERY:
          serviceImpl.userExtraInfoQuery((proto.Db.DB.UserKey) request,
              (io.vertx.core.Future<proto.Db.DB.UserExtendInfo>) io.vertx.core.Future.<proto.Db.DB.UserExtendInfo>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.Db.DB.UserExtendInfo>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static final class DBServiceDescriptorSupplier implements io.grpc.protobuf.ProtoFileDescriptorSupplier {
    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return proto.Db.getDescriptor();
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (DBServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new DBServiceDescriptorSupplier())
              .addMethod(METHOD_USER_QUERY)
              .addMethod(METHOD_USER_UPDATE_INFO)
              .addMethod(METHOD_USER_REGISTER)
              .addMethod(METHOD_USER_LOGIN)
              .addMethod(METHOD_USER_LOGOUT)
              .addMethod(METHOD_USER_EXTRA_INFO_QUERY)
              .build();
        }
      }
    }
    return result;
  }
}
