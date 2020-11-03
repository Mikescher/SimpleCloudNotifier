package com.blackforestbytes.simplecloudnotifier.model;

import android.util.Log;
import android.view.View;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.lambda.Func5to0;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.service.FBMService;

import org.joda.time.Instant;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.json.JSONTokener;

import java.io.IOException;
import java.net.URLEncoder;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.ResponseBody;

public class ServerCommunication
{
    public static final String PAGE_URL_LONG  = "https://simplecloudnotifier.blackforestbytes.com/";
    public static final String PAGE_URL_SHORT = "https://scn.blackforestbytes.com/";
    public static final String BASE_URL = "https://scn.blackforestbytes.com/api/";

    private static final OkHttpClient client = new OkHttpClient();

    private ServerCommunication(){ throw new Error("no."); }

    public static void register(String token, View loader, boolean pro, String pro_token)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "register.php?fcm_token=" + token + "&pro=" + pro + "&pro_token=" + URLEncoder.encode(pro_token, "utf-8"))
                    .build();

            client.newCall(request).enqueue(new Callback()
            {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    handleError("register", call, null, Str.Empty, true, e);
                    SCNApp.runOnUiThread(() -> { if (loader!=null)loader.setVisibility(View.GONE); });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful()) throw new IOException("Unexpected code " + response);
                        if (responseBody ==  null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("register", call, response, r);
                            return;
                        }

                        SCNSettings.inst().user_id          = json_int(json, "user_id");
                        SCNSettings.inst().user_key         = json_str(json, "user_key");
                        SCNSettings.inst().fcm_token_server = token;
                        SCNSettings.inst().quota_curr       = json_int(json, "quota");
                        SCNSettings.inst().quota_max        = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server   = json_bool(json, "is_pro");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();

                        handleSuccess("register", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("register", call, response, r, false, e);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> { if (loader!=null)loader.setVisibility(View.GONE); });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("register", null, null, Str.Empty, false, e);
        }
    }

    public static void updateFCMToken(int id, String key, String token, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "update.php?user_id="+id+"&user_key="+key+"&fcm_token="+token)
                    .build();

            client.newCall(request).enqueue(new Callback()
            {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    handleError("update<1>", call, null, Str.Empty, true, e);
                    SCNApp.runOnUiThread(() -> { if (loader!=null)loader.setVisibility(View.GONE); });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful()) throw new IOException("Unexpected code " + response);
                        if (responseBody ==  null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("update<1>", call, response, r);
                            return;
                        }

                        SCNSettings.inst().user_id          = json_int(json, "user_id");
                        SCNSettings.inst().user_key         = json_str(json, "user_key");
                        SCNSettings.inst().fcm_token_server = token;
                        SCNSettings.inst().quota_curr       = json_int(json, "quota");
                        SCNSettings.inst().quota_max        = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server   = json_bool(json, "is_pro");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();

                        handleSuccess("update<1>", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("update<1>", call, response, r, false, e);
                        SCNApp.showToast("Communication with server failed", 4000);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> { if (loader!=null)loader.setVisibility(View.GONE); });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("update<1>", null, null, Str.Empty, false, e);
        }
    }

    public static void resetSecret(int id, String key, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "update.php?user_id=" + id + "&user_key=" + key)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    handleError("update<1>", call, null, Str.Empty, true, e);
                    SCNApp.showToast("Communication with server failed", 4000);
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success")) {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("update<2>", call, response, r);
                            return;
                        }

                        SCNSettings.inst().user_id          = json_int(json, "user_id");
                        SCNSettings.inst().user_key         = json_str(json, "user_key");
                        SCNSettings.inst().quota_curr       = json_int(json, "quota");
                        SCNSettings.inst().quota_max        = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server   = json_bool(json, "is_pro");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();

                        handleSuccess("update<2>", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("update<2>", call, response, r, false, e);
                        SCNApp.showToast("Communication with server failed", 4000);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> {
                            if (loader != null) loader.setVisibility(View.GONE);
                        });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("update<2>", null, null, Str.Empty, false, e);
        }
    }

    public static void info(int id, String key, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "info.php?user_id=" + id + "&user_key=" + key)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    handleError("info", call, null, Str.Empty, true, e);
                    SCNApp.runOnUiThread(() -> {
                        if (loader != null) loader.setVisibility(View.GONE);
                    });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("info", call, response, r);

                            int errid = json.optInt("errid", 0);

                            if (errid == 201 || errid == 202 || errid == 203 || errid == 204)
                            {
                                // user not found or auth failed

                                SCNSettings.inst().user_id          = -1;
                                SCNSettings.inst().user_key         = "";
                                SCNSettings.inst().quota_curr       = 0;
                                SCNSettings.inst().quota_max        = 0;
                                SCNSettings.inst().promode_server   = false;
                                SCNSettings.inst().fcm_token_server = "";
                                SCNSettings.inst().save();

                                SCNApp.refreshAccountTab();
                            }

                            return;
                        }

                        SCNSettings.inst().user_id        = json_int(json, "user_id");
                        SCNSettings.inst().quota_curr     = json_int(json, "quota");
                        SCNSettings.inst().quota_max      = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server = json_bool(json, "is_pro");
                        if (!json_bool(json, "fcm_token_set")) SCNSettings.inst().fcm_token_server = "";
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();

                        if (json_int(json, "unack_count")>0) ServerCommunication.requery(id, key, loader);

                        handleSuccess("info", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("info", call, response, r, false, e);
                        SCNApp.showToast("Communication with server failed", 4000);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> { if (loader != null) loader.setVisibility(View.GONE); });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("info", null, null, Str.Empty, false, e);
        }
    }

    public static void requery(int id, String key, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "requery.php?user_id=" + id + "&user_key=" + key)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    handleError("requery", call, null, Str.Empty, true, e);
                    SCNApp.runOnUiThread(() -> {
                        if (loader != null) loader.setVisibility(View.GONE);
                    });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("requery", call, response, r);
                            return;
                        }

                        int count = json_int(json, "count");
                        JSONArray arr = json.getJSONArray("data");
                        for (int i = 0; i < count; i++)
                        {
                            JSONObject o = arr.getJSONObject(i);

                            long time         = json_lng(o, "timestamp");
                            String title      = json_str(o, "title");
                            String content    = json_str(o, "body");
                            PriorityEnum prio = PriorityEnum.parseAPI(json_int(o, "priority"));
                            long scn_id       = json_lng(o, "scn_msg_id");

                            FBMService.recieveData(time, title, content, prio, scn_id, true);
                        }

                        handleSuccess("requery", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("requery", call, response, r, false, e);
                        SCNApp.showToast("Communication with server failed", 4000);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> {
                            if (loader != null) loader.setVisibility(View.GONE);
                        });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("requery", null, null, Str.Empty, false, e);
        }
    }

    public static void upgrade(int id, String key, View loader, boolean pro, String pro_token)
    {
        try
        {
            SCNApp.runOnUiThread(() -> { if (loader != null) loader.setVisibility(View.GONE); });

            Request request = new Request.Builder()
                    .url(BASE_URL + "upgrade.php?user_id=" + id + "&user_key=" + key + "&pro=" + pro + "&pro_token=" + URLEncoder.encode(pro_token, "utf-8"))
                    .build();

            client.newCall(request).enqueue(new Callback()
            {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    handleError("upgrade", call, null, Str.Empty, true, e);
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful()) throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success")) {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("upgrade", call, response, r);
                            return;
                        }

                        SCNSettings.inst().user_id          = json_int(json, "user_id");
                        SCNSettings.inst().quota_curr       = json_int(json, "quota");
                        SCNSettings.inst().quota_max        = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server   = json_bool(json, "is_pro");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();

                        handleSuccess("upgrade", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("upgrade", call, response, r, false, e);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> { if (loader != null) loader.setVisibility(View.GONE); });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("upgrade", null, null, Str.Empty, false, e);
        }
    }

    public static void ack(int id, String key, long msg_scn_id)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "ack.php?user_id=" + id + "&user_key=" + key + "&scn_msg_id=" + msg_scn_id)
                    .build();

            client.newCall(request).enqueue(new Callback()
            {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    handleError("ack", call, null, Str.Empty, true, e);
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("ack", call, response, r);
                        }

                        handleSuccess("ack", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("ack", call, response, r, false, e);
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("ack", null, null, Str.Empty, false, e);
        }
    }

    public static void expand(int id, String key, long scn_msg_id, View loader, Func5to0<String, String, PriorityEnum, Long, Long> okResult)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "expand.php?user_id=" + id + "&user_key=" + key + "&scn_msg_id=" + scn_msg_id)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    handleError("expand", call, null, Str.Empty, true, e);
                    SCNApp.runOnUiThread(() -> { if (loader != null) loader.setVisibility(View.GONE); });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    String r = Str.Empty;
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        r = responseBody.string();
                        Log.d("Server::Response", request.url().toString()+"\n"+r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            handleNonSuccess("expand", call, response, r);
                            return;
                        }

                        JSONObject o = json.getJSONObject("data");

                        long time         = json_lng(o, "timestamp");
                        String title      = json_str(o, "title");
                        String content    = json_str(o, "body");
                        PriorityEnum prio = PriorityEnum.parseAPI(json_int(o, "priority"));
                        long scn_id       = json_lng(o, "scn_msg_id");

                        okResult.invoke(title, content, prio, time, scn_id);

                        handleSuccess("expand", call, response, r);
                    }
                    catch (Exception e)
                    {
                        handleError("expand", call, response, r, false, e);
                    }
                    finally
                    {
                        SCNApp.runOnUiThread(() -> { if (loader != null) loader.setVisibility(View.GONE); });
                    }
                }
            });
        }
        catch (Exception e)
        {
            handleError("expand", null, null, Str.Empty, false, e);
        }
    }

    private static boolean json_bool(JSONObject o, String key) throws JSONException
    {
        Object v = o.get(key);
        if (v instanceof Integer) return ((int)v) != 0;
        if (v instanceof Boolean) return ((boolean)v);
        if (v instanceof String) return !Str.equals(((String)v), "0") && !Str.equals(((String)v), "false");

        return o.getBoolean(key);
    }

    private static int json_int(JSONObject o, String key) throws JSONException
    {
        return o.getInt(key);
    }

    private static long json_lng(JSONObject o, String key) throws JSONException
    {
        return o.getLong(key);
    }

    private static String json_str(JSONObject o, String key) throws JSONException
    {
        return o.getString(key);
    }

    private static void handleSuccess(String source, Call call, Response resp, String respBody)
    {
        Log.d("SC:"+source, respBody);

        try
        {
            Instant i  = Instant.now();
            String s   = source;
            String u   = call.request().url().toString();
            int rc     = resp.code();
            String r   = respBody;
            LogLevel l = LogLevel.INFO;

            SingleQuery q = new SingleQuery(l, i, s, u, r, rc, "SUCCESS");
            QueryLog.inst().add(q);
        }
        catch (Exception e2)
        {
            Log.e("SC:HandleSuccess", e2.toString());
        }
    }

    private static void handleNonSuccess(String source, Call call, Response resp, String respBody)
    {
        Log.d("SC:"+source, respBody);

        try
        {
            Instant i  = Instant.now();
            String s   = source;
            String u   = call.request().url().toString();
            int rc     = resp.code();
            String r   = respBody;
            LogLevel l = LogLevel.WARN;

            SingleQuery q = new SingleQuery(l, i, s, u, r, rc, "NON-SUCCESS");
            QueryLog.inst().add(q);
        }
        catch (Exception e2)
        {
            Log.e("SC:HandleSuccess", e2.toString());
        }
    }

    private static void handleError(String source, Call call, Response resp, String respBody, boolean isio, Exception e)
    {
        Log.e("SC:"+source, e.toString());

        if (isio)
        {
            SCNApp.showToast("Can't connect to server", 3000);
        }
        else
        {
            SCNApp.showToast("Communication with server failed", 4000);
        }

        try
        {
            Instant i  = Instant.now();
            String s   = source;
            String u   = (call==null)?Str.Empty:call.request().url().toString();
            int rc     = (resp==null)?-1:resp.code();
            String r   = respBody;
            LogLevel l = isio?LogLevel.WARN:LogLevel.ERROR;

            SingleQuery q = new SingleQuery(l, i, s, u, r, rc, e.toString());
            QueryLog.inst().add(q);
        }
        catch (Exception e2)
        {
            Log.e("SC:HandleError", e2.toString());
        }
    }
}
