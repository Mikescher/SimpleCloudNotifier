package com.blackforestbytes.simplecloudnotifier.model;

import android.util.Log;
import android.view.View;

import com.blackforestbytes.simplecloudnotifier.SCNApp;

import org.json.JSONObject;
import org.json.JSONTokener;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.ResponseBody;

public class ServerCommunication
{
    public static final String BASE_URL = "https://scn.blackforestbytes.com/";

    private static final OkHttpClient client = new OkHttpClient();

    private ServerCommunication(){ throw new Error("no."); }

    public static void register(String token, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "register.php?fcm_token="+token)
                    .build();

            client.newCall(request).enqueue(new Callback()
            {
                @Override
                public void onFailure(Call call, IOException e)
                {
                    e.printStackTrace();
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

                        if (!json.getBoolean("success"))
                        {
                            SCNApp.showToast(json.getString("message"), 4000);
                            return;
                        }

                        SCNSettings.inst().user_id          = json.getInt("user_id");
                        SCNSettings.inst().user_key         = json.getString("user_key");
                        SCNSettings.inst().fcm_token_server = token;
                        SCNSettings.inst().quota_curr       = json.getInt("quota");
                        SCNSettings.inst().quota_max        = json.getInt("quota_max");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();
                    }
                    catch (Exception e)
                    {
                        e.printStackTrace();
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
            e.printStackTrace();
            SCNApp.showToast("Communication with server failed", 4000);
        }
    }

    public static void update(int id, String key, String token, View loader)
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
                    e.printStackTrace();
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

                        if (!json.getBoolean("success"))
                        {
                            SCNApp.showToast(json.getString("message"), 4000);
                            return;
                        }

                        SCNSettings.inst().user_id          = json.getInt("user_id");
                        SCNSettings.inst().user_key         = json.getString("user_key");
                        SCNSettings.inst().fcm_token_server = token;
                        SCNSettings.inst().quota_curr       = json.getInt("quota");
                        SCNSettings.inst().quota_max        = json.getInt("quota_max");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();
                    }
                    catch (Exception e)
                    {
                        e.printStackTrace();
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
            e.printStackTrace();
            SCNApp.showToast("Communication with server failed", 4000);
        }
    }

    public static void update(int id, String key, View loader)
    {
        try
        {
            Request request = new Request.Builder()
                    .url(BASE_URL + "update.php?user_id=" + id + "&user_key=" + key)
                    .build();

            client.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    e.printStackTrace();
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

                        if (!json.getBoolean("success")) {
                            SCNApp.showToast(json.getString("message"), 4000);
                            return;
                        }

                        SCNSettings.inst().user_id = json.getInt("user_id");
                        SCNSettings.inst().user_key = json.getString("user_key");
                        SCNSettings.inst().quota_curr = json.getInt("quota");
                        SCNSettings.inst().quota_max = json.getInt("quota_max");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();
                    } catch (Exception e) {
                        e.printStackTrace();
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
            e.printStackTrace();
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
                    e.printStackTrace();
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

                        if (!json.getBoolean("success")) {
                            SCNApp.showToast(json.getString("message"), 4000);
                            return;
                        }

                        SCNSettings.inst().user_id = json.getInt("user_id");
                        SCNSettings.inst().quota_curr = json.getInt("quota");
                        SCNSettings.inst().quota_max = json.getInt("quota_max");
                        SCNSettings.inst().save();

                        SCNApp.refreshAccountTab();
                    } catch (Exception e) {
                        e.printStackTrace();
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
            e.printStackTrace();
            SCNApp.showToast("Communication with server failed", 4000);
        }
    }
}
