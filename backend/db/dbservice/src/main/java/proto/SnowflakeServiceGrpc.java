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
 * snowflake service definition
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.3.0)",
    comments = "Source: snowflake.proto")
public final class SnowflakeServiceGrpc {

  private SnowflakeServiceGrpc() {}

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

  public static final String SERVICE_NAME = "proto.SnowflakeService";

  // Static method descriptors that strictly reflect the proto.
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.SnowflakeOuterClass.Snowflake.Key,
      proto.SnowflakeOuterClass.Snowflake.Value> METHOD_NEXT =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.SnowflakeService", "Next"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.SnowflakeOuterClass.Snowflake.Key.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.SnowflakeOuterClass.Snowflake.Value.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.SnowflakeOuterClass.Snowflake.NullRequest,
      proto.SnowflakeOuterClass.Snowflake.UUID> METHOD_GET_UUID =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.SnowflakeService", "GetUUID"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.SnowflakeOuterClass.Snowflake.NullRequest.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.SnowflakeOuterClass.Snowflake.UUID.getDefaultInstance()));
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static final io.grpc.MethodDescriptor<proto.SnowflakeOuterClass.Snowflake.Key,
      proto.SnowflakeOuterClass.Snowflake.UUID> METHOD_GET_USER_ID =
      io.grpc.MethodDescriptor.create(
          io.grpc.MethodDescriptor.MethodType.UNARY,
          generateFullMethodName(
              "proto.SnowflakeService", "GetUserID"),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.SnowflakeOuterClass.Snowflake.Key.getDefaultInstance()),
          io.grpc.protobuf.ProtoUtils.marshaller(proto.SnowflakeOuterClass.Snowflake.UUID.getDefaultInstance()));

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static SnowflakeServiceStub newStub(io.grpc.Channel channel) {
    return new SnowflakeServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static SnowflakeServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new SnowflakeServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary and streaming output calls on the service
   */
  public static SnowflakeServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new SnowflakeServiceFutureStub(channel);
  }

  /**
   * Creates a new vertx stub that supports all call types for the service
   */
  public static SnowflakeServiceVertxStub newVertxStub(io.grpc.Channel channel) {
    return new SnowflakeServiceVertxStub(channel);
  }

  /**
   * <pre>
   * snowflake service definition
   * </pre>
   */
  public static abstract class SnowflakeServiceImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * Generate next serial number
     * </pre>
     */
    public void next(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.Value> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_NEXT, responseObserver);
    }

    /**
     * <pre>
     * UUID generator
     * </pre>
     */
    public void getUUID(proto.SnowflakeOuterClass.Snowflake.NullRequest request,
        io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_GET_UUID, responseObserver);
    }

    /**
     * <pre>
     * User ID generate
     * </pre>
     */
    public void getUserID(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID> responseObserver) {
      asyncUnimplementedUnaryCall(METHOD_GET_USER_ID, responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_NEXT,
            asyncUnaryCall(
              new MethodHandlers<
                proto.SnowflakeOuterClass.Snowflake.Key,
                proto.SnowflakeOuterClass.Snowflake.Value>(
                  this, METHODID_NEXT)))
          .addMethod(
            METHOD_GET_UUID,
            asyncUnaryCall(
              new MethodHandlers<
                proto.SnowflakeOuterClass.Snowflake.NullRequest,
                proto.SnowflakeOuterClass.Snowflake.UUID>(
                  this, METHODID_GET_UUID)))
          .addMethod(
            METHOD_GET_USER_ID,
            asyncUnaryCall(
              new MethodHandlers<
                proto.SnowflakeOuterClass.Snowflake.Key,
                proto.SnowflakeOuterClass.Snowflake.UUID>(
                  this, METHODID_GET_USER_ID)))
          .build();
    }
  }

  /**
   * <pre>
   * snowflake service definition
   * </pre>
   */
  public static final class SnowflakeServiceStub extends io.grpc.stub.AbstractStub<SnowflakeServiceStub> {
    private SnowflakeServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private SnowflakeServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SnowflakeServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new SnowflakeServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * Generate next serial number
     * </pre>
     */
    public void next(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.Value> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_NEXT, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * UUID generator
     * </pre>
     */
    public void getUUID(proto.SnowflakeOuterClass.Snowflake.NullRequest request,
        io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_GET_UUID, getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * User ID generate
     * </pre>
     */
    public void getUserID(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_GET_USER_ID, getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * snowflake service definition
   * </pre>
   */
  public static final class SnowflakeServiceBlockingStub extends io.grpc.stub.AbstractStub<SnowflakeServiceBlockingStub> {
    private SnowflakeServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private SnowflakeServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SnowflakeServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new SnowflakeServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * Generate next serial number
     * </pre>
     */
    public proto.SnowflakeOuterClass.Snowflake.Value next(proto.SnowflakeOuterClass.Snowflake.Key request) {
      return blockingUnaryCall(
          getChannel(), METHOD_NEXT, getCallOptions(), request);
    }

    /**
     * <pre>
     * UUID generator
     * </pre>
     */
    public proto.SnowflakeOuterClass.Snowflake.UUID getUUID(proto.SnowflakeOuterClass.Snowflake.NullRequest request) {
      return blockingUnaryCall(
          getChannel(), METHOD_GET_UUID, getCallOptions(), request);
    }

    /**
     * <pre>
     * User ID generate
     * </pre>
     */
    public proto.SnowflakeOuterClass.Snowflake.UUID getUserID(proto.SnowflakeOuterClass.Snowflake.Key request) {
      return blockingUnaryCall(
          getChannel(), METHOD_GET_USER_ID, getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * snowflake service definition
   * </pre>
   */
  public static final class SnowflakeServiceFutureStub extends io.grpc.stub.AbstractStub<SnowflakeServiceFutureStub> {
    private SnowflakeServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private SnowflakeServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SnowflakeServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new SnowflakeServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * Generate next serial number
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.SnowflakeOuterClass.Snowflake.Value> next(
        proto.SnowflakeOuterClass.Snowflake.Key request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_NEXT, getCallOptions()), request);
    }

    /**
     * <pre>
     * UUID generator
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.SnowflakeOuterClass.Snowflake.UUID> getUUID(
        proto.SnowflakeOuterClass.Snowflake.NullRequest request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_GET_UUID, getCallOptions()), request);
    }

    /**
     * <pre>
     * User ID generate
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<proto.SnowflakeOuterClass.Snowflake.UUID> getUserID(
        proto.SnowflakeOuterClass.Snowflake.Key request) {
      return futureUnaryCall(
          getChannel().newCall(METHOD_GET_USER_ID, getCallOptions()), request);
    }
  }

  /**
   * <pre>
   * snowflake service definition
   * </pre>
   */
  public static abstract class SnowflakeServiceVertxImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * Generate next serial number
     * </pre>
     */
    public void next(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.vertx.core.Future<proto.SnowflakeOuterClass.Snowflake.Value> response) {
      asyncUnimplementedUnaryCall(METHOD_NEXT, SnowflakeServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * UUID generator
     * </pre>
     */
    public void getUUID(proto.SnowflakeOuterClass.Snowflake.NullRequest request,
        io.vertx.core.Future<proto.SnowflakeOuterClass.Snowflake.UUID> response) {
      asyncUnimplementedUnaryCall(METHOD_GET_UUID, SnowflakeServiceGrpc.toObserver(response.completer()));
    }

    /**
     * <pre>
     * User ID generate
     * </pre>
     */
    public void getUserID(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.vertx.core.Future<proto.SnowflakeOuterClass.Snowflake.UUID> response) {
      asyncUnimplementedUnaryCall(METHOD_GET_USER_ID, SnowflakeServiceGrpc.toObserver(response.completer()));
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            METHOD_NEXT,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.SnowflakeOuterClass.Snowflake.Key,
                proto.SnowflakeOuterClass.Snowflake.Value>(
                  this, METHODID_NEXT)))
          .addMethod(
            METHOD_GET_UUID,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.SnowflakeOuterClass.Snowflake.NullRequest,
                proto.SnowflakeOuterClass.Snowflake.UUID>(
                  this, METHODID_GET_UUID)))
          .addMethod(
            METHOD_GET_USER_ID,
            asyncUnaryCall(
              new VertxMethodHandlers<
                proto.SnowflakeOuterClass.Snowflake.Key,
                proto.SnowflakeOuterClass.Snowflake.UUID>(
                  this, METHODID_GET_USER_ID)))
          .build();
    }
  }

  /**
   * <pre>
   * snowflake service definition
   * </pre>
   */
  public static final class SnowflakeServiceVertxStub extends io.grpc.stub.AbstractStub<SnowflakeServiceVertxStub> {
    private SnowflakeServiceVertxStub(io.grpc.Channel channel) {
      super(channel);
    }

    private SnowflakeServiceVertxStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SnowflakeServiceVertxStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new SnowflakeServiceVertxStub(channel, callOptions);
    }

    /**
     * <pre>
     * Generate next serial number
     * </pre>
     */
    public void next(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.SnowflakeOuterClass.Snowflake.Value>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_NEXT, getCallOptions()), request, SnowflakeServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * UUID generator
     * </pre>
     */
    public void getUUID(proto.SnowflakeOuterClass.Snowflake.NullRequest request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.SnowflakeOuterClass.Snowflake.UUID>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_GET_UUID, getCallOptions()), request, SnowflakeServiceGrpc.toObserver(response));
    }

    /**
     * <pre>
     * User ID generate
     * </pre>
     */
    public void getUserID(proto.SnowflakeOuterClass.Snowflake.Key request,
        io.vertx.core.Handler<io.vertx.core.AsyncResult<proto.SnowflakeOuterClass.Snowflake.UUID>> response) {
      asyncUnaryCall(
          getChannel().newCall(METHOD_GET_USER_ID, getCallOptions()), request, SnowflakeServiceGrpc.toObserver(response));
    }
  }

  private static final int METHODID_NEXT = 0;
  private static final int METHODID_GET_UUID = 1;
  private static final int METHODID_GET_USER_ID = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final SnowflakeServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(SnowflakeServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_NEXT:
          serviceImpl.next((proto.SnowflakeOuterClass.Snowflake.Key) request,
              (io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.Value>) responseObserver);
          break;
        case METHODID_GET_UUID:
          serviceImpl.getUUID((proto.SnowflakeOuterClass.Snowflake.NullRequest) request,
              (io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID>) responseObserver);
          break;
        case METHODID_GET_USER_ID:
          serviceImpl.getUserID((proto.SnowflakeOuterClass.Snowflake.Key) request,
              (io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID>) responseObserver);
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
    private final SnowflakeServiceVertxImplBase serviceImpl;
    private final int methodId;

    VertxMethodHandlers(SnowflakeServiceVertxImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_NEXT:
          serviceImpl.next((proto.SnowflakeOuterClass.Snowflake.Key) request,
              (io.vertx.core.Future<proto.SnowflakeOuterClass.Snowflake.Value>) io.vertx.core.Future.<proto.SnowflakeOuterClass.Snowflake.Value>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.Value>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_GET_UUID:
          serviceImpl.getUUID((proto.SnowflakeOuterClass.Snowflake.NullRequest) request,
              (io.vertx.core.Future<proto.SnowflakeOuterClass.Snowflake.UUID>) io.vertx.core.Future.<proto.SnowflakeOuterClass.Snowflake.UUID>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID>) responseObserver).onNext(ar.result());
                  responseObserver.onCompleted();
                } else {
                  responseObserver.onError(ar.cause());
                }
              }));
          break;
        case METHODID_GET_USER_ID:
          serviceImpl.getUserID((proto.SnowflakeOuterClass.Snowflake.Key) request,
              (io.vertx.core.Future<proto.SnowflakeOuterClass.Snowflake.UUID>) io.vertx.core.Future.<proto.SnowflakeOuterClass.Snowflake.UUID>future().setHandler(ar -> {
                if (ar.succeeded()) {
                  ((io.grpc.stub.StreamObserver<proto.SnowflakeOuterClass.Snowflake.UUID>) responseObserver).onNext(ar.result());
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

  private static final class SnowflakeServiceDescriptorSupplier implements io.grpc.protobuf.ProtoFileDescriptorSupplier {
    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return proto.SnowflakeOuterClass.getDescriptor();
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (SnowflakeServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new SnowflakeServiceDescriptorSupplier())
              .addMethod(METHOD_NEXT)
              .addMethod(METHOD_GET_UUID)
              .addMethod(METHOD_GET_USER_ID)
              .build();
        }
      }
    }
    return result;
  }
}
