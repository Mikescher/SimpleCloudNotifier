package com.blackforestbytes.simplecloudnotifier.model;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;
import android.view.View;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.google.firebase.iid.FirebaseInstanceId;

public class SCNSettings
{
    private final static Object _lock = new Object();
    private static SCNSettings _inst = null;
    public static SCNSettings inst()
    {
        synchronized (_lock)
        {
            if (_inst != null) return _inst;
            return _inst = new SCNSettings();
        }
    }

    public int    quota_curr;
    public int    quota_max;
    public int    user_id;
    public String user_key;

    public String fcm_token_local;
    public String fcm_token_server;

    public SCNSettings()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("Config", Context.MODE_PRIVATE);

        quota_curr       = sharedPref.getInt("quota_curr", 0);
        quota_max        = sharedPref.getInt("quota_max", 0);
        user_id          = sharedPref.getInt("user_id", -1);
        user_key         = sharedPref.getString("user_key", "");
        fcm_token_local  = sharedPref.getString("fcm_token_local", "");
        fcm_token_server = sharedPref.getString("fcm_token_server", "");
    }

    public void save()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("Config", Context.MODE_PRIVATE);
        SharedPreferences.Editor e = sharedPref.edit();

        e.putInt("quota_curr", quota_curr);
        e.putInt("quota_max", quota_max);
        e.putInt("user_id", user_id);
        e.putString("user_key", user_key);
        e.putString("fcm_token_local", fcm_token_local);
        e.putString("fcm_token_server", fcm_token_server);

        e.apply();
    }

    public boolean isConnected()
    {
        return user_id>=0 && user_key != null && !user_key.isEmpty();
    }

    public String createOnlineURL()
    {
        if (!isConnected()) return ServerCommunication.BASE_URL + "index.php";
        return ServerCommunication.BASE_URL + "index.php?preset_user_id="+user_id+"&preset_user_key="+user_key;
    }

    public void setServerToken(String token, View loader)
    {
        if (isConnected())
        {
            fcm_token_local = token;
            save();
            if (!fcm_token_local.equals(fcm_token_server)) ServerCommunication.update(user_id, user_key, fcm_token_local, loader);
        }
        else
        {
            fcm_token_local = token;
            save();
            ServerCommunication.register(fcm_token_local, loader);
        }
    }

    public void work(Activity a)
    {
        FirebaseInstanceId.getInstance().getInstanceId().addOnSuccessListener(a, instanceIdResult ->
        {
            String newToken = instanceIdResult.getToken();
            Log.e("FB::GetInstanceId", newToken);
            SCNSettings.inst().setServerToken(newToken, null);
        }).addOnCompleteListener(r ->
        {
            if (isConnected()) ServerCommunication.info(user_id, user_key, null);
        });
    }

    public void reset(View loader)
    {
        if (!isConnected()) return;

        ServerCommunication.update(user_id, user_key, loader);
    }

    public void refresh(View loader, Activity a)
    {
        if (isConnected())
        {
            ServerCommunication.info(user_id, user_key, loader);
        }
        else
        {
            FirebaseInstanceId.getInstance().getInstanceId().addOnSuccessListener(a, instanceIdResult ->
            {
                String newToken = instanceIdResult.getToken();
                Log.e("FB::GetInstanceId", newToken);
                SCNSettings.inst().setServerToken(newToken, loader);
            }).addOnCompleteListener(r ->
            {
                if (isConnected()) ServerCommunication.info(user_id, user_key, null);
            });
        }
    }
}
