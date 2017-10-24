package com.omgo.webservice.handler;

import com.omgo.webservice.AgentManager;
import com.omgo.webservice.Utils;
import com.omgo.webservice.model.HttpStatus;
import com.omgo.webservice.model.ModelConverter;
import com.omgo.webservice.service.Services;
import io.grpc.ManagedChannel;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.http.HttpServerResponse;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;
import proto.DBServiceGrpc;
import proto.Db;

// request

/*
app_language:	zh-rCN
app_version:	0.0.1
avatar:	http://gravatar.com/avatar/fddd805f5195dab1951784d2a6b69448?s=200
birthday:	531262800000
country:	CN
device_type:	1
email:	test1@qq.com
gender:	2
mcc:	460
nickname:	test1
os:	android 8 oreo
os_locale:	zh-rCN
phone:	18600001111
secret:	p4ssw0rd
timezone:	8
 */

// getResponse
/*
{
"usn": 0,
"uid": 0,
"app_language": "zh-rCN",
"app_version": "0.0.1",
"avatar": "http://gravatar.com/avatar/fddd805f5195dab1951784d2a6b69448?s=200",
"birthday": 531262800000,
"country": "CN",
"device_type": 1,
"email": "dearunclexiao@qq.com",
"email_verified": false,
"gender": 2,
"is_official": false,
"is_robot": false,
"last_ip": "127.0.0.1",
"last_login": 1505406356268,
"login_count": 1,
"mcc": 460,
"nickname": "dearunclexiao",
"os": "android 8 oreo",
"os_locale": "zh-rCN",
"phone": "18600001111",
"phone_verified": false,
"premium_end": 0,
"premium_exp": 0,
"premium_level": 0,
"secret": "",
"since": 0,
"social_id": "",
"social_name": "",
"social_verified": false,
"status": 0,
"timezone": 0,
"token": "vj+XGmrXueG9SNQr03Phog=="
}
 */

public class RegisterHandler extends BaseHandler implements Services.Pool.OnChangeListener {

    private DBServiceGrpc.DBServiceVertxStub dbServiceVertxStub;
    private Services.Pool dataServicePool;
    private ManagedChannel channel;

    public RegisterHandler(Vertx vertx, Services.Pool servicePool) {
        super(vertx);
        notRequireValidNonce();
        notRequireValidSession();
        notRequireValidEncryption();

        this.dataServicePool = servicePool;
        init();
    }

    private void init() {
        channel = dataServicePool.getClient();
        if (channel != null) {
            dbServiceVertxStub = DBServiceGrpc.newVertxStub(channel);
        }
        dataServicePool.addOnChangeListener(this);
    }

    @Override
    protected void handle(RoutingContext routingContext, HttpServerResponse response) {
        HttpServerRequest request = super.getRequest(routingContext);

        JsonObject registerJson = super.getHeaderJson(request);
        String app_language = registerJson.getString(ModelConverter.KEY_APP_LANGUAGE, "");
        String app_version = registerJson.getString(ModelConverter.KEY_APP_VERSION, "");
        String avatar = registerJson.getString(ModelConverter.KEY_AVATAR, "");
        String birthday = registerJson.getString(ModelConverter.KEY_BIRTHDAY, "");
        String country = registerJson.getString(ModelConverter.KEY_COUNTRY, "");
        String device_type = registerJson.getString(ModelConverter.KEY_DEVICE_TYPE, "");
        String email = registerJson.getString(ModelConverter.KEY_EMAIL, "");
        String gender = registerJson.getString(ModelConverter.KEY_GENDER, "");
        String mcc = registerJson.getString(ModelConverter.KEY_MCC, "");
        String nickname = registerJson.getString(ModelConverter.KEY_NICKNAME, "");
        String os = registerJson.getString(ModelConverter.KEY_OS, "");
        String os_locale = registerJson.getString(ModelConverter.KEY_OS_LOCALE, "");
        String phone = registerJson.getString(ModelConverter.KEY_PHONE, "");
        String secret = registerJson.getString(ModelConverter.KEY_SECRET, "");
        String timezone = registerJson.getString(ModelConverter.KEY_TIMEZONE, "");

        long birthdayLong = Utils.isEmptyString(birthday) ? 0L : Long.parseLong(birthday);
        int genderInt = Utils.isEmptyString(gender) ? 0 : Integer.parseInt(gender);
        int deviceType = Utils.isEmptyString(device_type) ? 0 : Integer.parseInt(device_type);
        int mccInt = Utils.isEmptyString(mcc) ? 0 : Integer.parseInt(mcc);
        int timezoneInt = Utils.isEmptyString(timezone) ? 0 : Integer.parseInt(timezone);


        Db.DB.UserEntry.Builder userEntryBuilder = Db.DB.UserEntry.newBuilder();
        userEntryBuilder
            .setAppLanguage(app_language)
            .setAppVersion(app_version)
            .setAvatar(avatar)
            .setBirthday(birthdayLong)
            .setCountry(country)
            .setDeviceType(deviceType)
            .setEmail(email)
            .setGender(genderInt)
            .setLastIp(request.connection().remoteAddress().host())
            .setMcc(mccInt)
            .setNickname(nickname)
            .setOs(os)
            .setOsLocale(os_locale)
            .setPhone(phone)
            .setSecret(secret)
            .setTimezone(timezoneInt);

        dbServiceVertxStub.userRegister(userEntryBuilder.build(), res -> {
            if (res.succeeded()) {
                Db.DB.StatusCode code = res.result().getResult().getStatus();
                if (code == Db.DB.StatusCode.STATUS_OK) {
                    JsonObject resultJson = ModelConverter.userEntry2Json(res.result().getUser());

                    String token = resultJson.getString(ModelConverter.KEY_TOKEN);
                    setSessionToken(routingContext, resultJson.getString(ModelConverter.KEY_TOKEN));

                    JsonObject rspJson = getResponseJson();
                    rspJson.put(ModelConverter.KEY_USER_INFO, resultJson);
                    rspJson.put(ModelConverter.KEY_HOSTS, AgentManager.getInstance().getHostList());
                    response.write(rspJson.encode()).end();
                } else {
                    LOGGER.info(res.result().getResult());
                    routingContext.fail(HttpStatus.INTERNAL_SERVER_ERROR.code);
                }
            } else {
                LOGGER.info(res.cause());
                routingContext.fail(HttpStatus.INTERNAL_SERVER_ERROR.code);
            }
        });
    }

    @Override
    public void onServiceAdded(Services.Pool pool) {
        if (channel == null) {
            LOGGER.info("dataservice online, init...");
            init();
        }
    }

    @Override
    public void onServiceRemoved(Services.Pool pool) {
        if (channel != null && channel.isShutdown()) {
            LOGGER.info("dataservice offline, try re-init");
            channel = null;
            dbServiceVertxStub = null;
            init();
        }
    }
}
