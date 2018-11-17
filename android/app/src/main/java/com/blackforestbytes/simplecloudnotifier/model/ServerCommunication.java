package com.blackforestbytes.simplecloudnotifier.model;

import android.util.Log;
import android.view.View;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.service.FBMService;

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
    public static final String BASE_URL = /*SCNApp.LOCAL_DEBUG ? "http://localhost:1010/" : */"https://scn.blackforestbytes.com/api/";

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
                    Log.e("SC:register", e.toString());
                    SCNApp.showToast("Communication with server failed", 4000);
                    SCNApp.runOnUiThread(() -> { if (loader!=null)loader.setVisibility(View.GONE); });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful()) throw new IOException("Unexpected code " + response);
                        if (responseBody ==  null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
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
                    }
                    catch (Exception e)
                    {
                        Log.e("SC:register", e.toString());
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
            Log.e("SC:register", e.toString());
            SCNApp.showToast("Communication with server failed", 4000);
        }
    }

    public static void updateFCMToken(int id, String key, String token, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "updateFCMToken.php?user_id="+id+"&user_key="+key+"&fcm_token="+token)
                    .build();

            client.newCall(request).enqueue(new Callback()
            {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    Log.e("SC:update_1", e.toString());
                    SCNApp.showToast("Communication with server failed", 4000);
                    SCNApp.runOnUiThread(() -> { if (loader!=null)loader.setVisibility(View.GONE); });
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    try (ResponseBody responseBody = response.body())
                    {
                        if (!response.isSuccessful()) throw new IOException("Unexpected code " + response);
                        if (responseBody ==  null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
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
                    }
                    catch (Exception e)
                    {
                        Log.e("SC:update_1", e.toString());
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
            Log.e("SC:update_1", e.toString());
            SCNApp.showToast("Communication with server failed", 4000);
        }
    }

    public static void resetSecret(int id, String key, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "updateFCMToken.php?user_id=" + id + "&user_key=" + key)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    Log.e("SC:update_2", e.toString());
                    SCNApp.showToast("Communication with server failed", 4000);
                }

                @Override
                public void onResponse(Call call, Response response) {
                    try (ResponseBody responseBody = response.body()) {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success")) {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            return;
                        }

                        SCNSettings.inst().user_id          = json_int(json, "user_id");
                        SCNSettings.inst().user_key         = json_str(json, "user_key");
                        SCNSettings.inst().quota_curr       = json_int(json, "quota");
                        SCNSettings.inst().quota_max        = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server   = json_bool(json, "is_pro");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();
                    } catch (Exception e) {
                        Log.e("SC:update_2", e.toString());
                        SCNApp.showToast("Communication with server failed", 4000);
                    } finally {
                        SCNApp.runOnUiThread(() -> {
                            if (loader != null) loader.setVisibility(View.GONE);
                        });
                    }
                }
            });
        }
        catch (Exception e)
        {
            Log.e("SC:update_2", e.toString());
            SCNApp.showToast("Communication with server failed", 4000);
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
                    Log.e("SC:info", e.toString());
                    SCNApp.showToast("Communication with server failed", 4000);
                    SCNApp.runOnUiThread(() -> {
                        if (loader != null) loader.setVisibility(View.GONE);
                    });
                }

                @Override
                public void onResponse(Call call, Response response) {
                    try (ResponseBody responseBody = response.body()) {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);

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

                        if (json_int(json, "unack_count")>0)
                        {
                            ServerCommunication.requery(id, key, loader);
                        }

                    } catch (Exception e) {
                        Log.e("SC:info", e.toString());
                        SCNApp.showToast("Communication with server failed", 4000);
                    } finally {
                        SCNApp.runOnUiThread(() -> {
                            if (loader != null) loader.setVisibility(View.GONE);
                        });
                    }
                }
            });
        }
        catch (Exception e)
        {
            Log.e("SC:info", e.toString());
            SCNApp.showToast("Communication with server failed", 4000);
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
                    Log.e("SC:requery", e.toString());
                    SCNApp.showToast("Communication with server failed", 4000);
                    SCNApp.runOnUiThread(() -> {
                        if (loader != null) loader.setVisibility(View.GONE);
                    });
                }

                @Override
                public void onResponse(Call call, Response response) {
                    try (ResponseBody responseBody = response.body()) {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success"))
                        {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            return;
                        }

                        int count = json_int(json, "count");
                        JSONArray arr = json.getJSONArray("data");
                        for (int i = 0; i < count; i++)
                        {
                            JSONObject o = arr.getJSONObject(0);

                            long time         = json_lng(o, "timestamp");
                            String title      = json_str(o, "title");
                            String content    = json_str(o, "body");
                            PriorityEnum prio = PriorityEnum.parseAPI(json_int(o, "priority"));
                            long scn_id       = json_lng(o, "scn_msg_id");

                            FBMService.recieveData(time, title, content, prio, scn_id, true);
                        }

                    } catch (Exception e) {
                        Log.e("SC:info", e.toString());
                        SCNApp.showToast("Communication with server failed", 4000);
                    } finally {
                        SCNApp.runOnUiThread(() -> {
                            if (loader != null) loader.setVisibility(View.GONE);
                        });
                    }
                }
            });
        }
        catch (Exception e)
        {
            Log.e("SC:requery", e.toString());
            SCNApp.showToast("Communication with server failed", 4000);
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

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    Log.e("SC:upgrade", e.toString());
                    SCNApp.showToast("Communication with server failed", 4000);
                }

                @Override
                public void onResponse(Call call, Response response) {
                    try (ResponseBody responseBody = response.body()) {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success")) {
                            SCNApp.showToast(json_str(json, "message"), 4000);
                            return;
                        }

                        SCNSettings.inst().user_id          = json_int(json, "user_id");
                        SCNSettings.inst().quota_curr       = json_int(json, "quota");
                        SCNSettings.inst().quota_max        = json_int(json, "quota_max");
                        SCNSettings.inst().promode_server   = json_bool(json, "is_pro");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();
                    } catch (Exception e) {
                        Log.e("SC:upgrade", e.toString());
                        SCNApp.showToast("Communication with server failed", 4000);
                    } finally {
                        SCNApp.runOnUiThread(() -> { if (loader != null) loader.setVisibility(View.GONE); });
                    }
                }
            });
        }
        catch (Exception e)
        {
            e.printStackTrace();
            SCNApp.showToast("Communication with server failed", 4000);
        }
    }

    public static void ack(int id, String key, CMessage msg)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "ack.php?user_id=" + id + "&user_key=" + key + "&scn_msg_id=" + msg.SCN_ID)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    Log.e("SC:ack", e.toString());
                }

                @Override
                public void onResponse(Call call, Response response)
                {
                    try (ResponseBody responseBody = response.body()) {
                        if (!response.isSuccessful())
                            throw new IOException("Unexpected code " + response);
                        if (responseBody == null) throw new IOException("No response");

                        String r = responseBody.string();
                        Log.d("Server::Response", r);

                        JSONObject json = (JSONObject) new JSONTokener(r).nextValue();

                        if (!json_bool(json, "success")) SCNApp.showToast(json_str(json, "message"), 4000);

                    } catch (Exception e) {
                        Log.e("SC:ack", e.toString());
                        SCNApp.showToast("Communication with server failed", 4000);
                    }
                }
            });
        }
        catch (Exception e)
        {
            Log.e("SC:ack", e.toString());
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
}
