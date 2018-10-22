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

    // ------------------------------------------------------------

    public final static Integer[] CHOOSABLE_CACHE_SIZES = new Integer[]{20, 50, 100, 200, 500, 1000, 2000, 5000};

    // ------------------------------------------------------------

    public int    quota_curr;
    public int    quota_max;
    public int    user_id;
    public String user_key;

    public String fcm_token_local;
    public String fcm_token_server;

    // ------------------------------------------------------------

    public boolean Enabled = true;
    public int LocalCacheSize = 500;

    public final NotificationSettings PriorityLow  = new NotificationSettings();
    public final NotificationSettings PriorityNorm = new NotificationSettings();
    public final NotificationSettings PriorityHigh = new NotificationSettings();

    // ------------------------------------------------------------

    public SCNSettings()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("Config", Context.MODE_PRIVATE);

        quota_curr       = sharedPref.getInt(   "quota_curr",        0);
        quota_max        = sharedPref.getInt(   "quota_max",         0);
        user_id          = sharedPref.getInt(   "user_id",          -1);
        user_key         = sharedPref.getString("user_key",         "");
        fcm_token_local  = sharedPref.getString("fcm_token_local",  "");
        fcm_token_server = sharedPref.getString("fcm_token_server", "");

        Enabled                     = sharedPref.getBoolean("app_enabled",  Enabled);
        LocalCacheSize              = sharedPref.getInt("local_cache_size", LocalCacheSize);

        PriorityLow.EnableLED        = sharedPref.getBoolean("priority_low:enabled_led",       PriorityLow.EnableLED);
        PriorityLow.EnableSound      = sharedPref.getBoolean("priority_low:enabled_sound",     PriorityLow.EnableSound);
        PriorityLow.EnableVibration  = sharedPref.getBoolean("priority_low:enabled_vibration", PriorityLow.EnableVibration);
        PriorityLow.RepeatSound      = sharedPref.getBoolean("priority_low:repeat_sound",      PriorityLow.RepeatSound);
        PriorityLow.SoundName        = sharedPref.getString( "priority_low:sound_name",        PriorityLow.SoundName);
        PriorityLow.SoundSource      = sharedPref.getString( "priority_low:sound_source",      PriorityLow.SoundSource);
        PriorityLow.LEDColor         = sharedPref.getInt(    "priority_low:led_color",         PriorityLow.LEDColor);

        PriorityNorm.EnableLED       = sharedPref.getBoolean("priority_norm:enabled_led",       PriorityNorm.EnableLED);
        PriorityNorm.EnableSound     = sharedPref.getBoolean("priority_norm:enabled_sound",     PriorityNorm.EnableSound);
        PriorityNorm.EnableVibration = sharedPref.getBoolean("priority_norm:enabled_vibration", PriorityNorm.EnableVibration);
        PriorityNorm.RepeatSound     = sharedPref.getBoolean("priority_norm:repeat_sound",      PriorityNorm.RepeatSound);
        PriorityNorm.SoundName       = sharedPref.getString( "priority_norm:sound_name",        PriorityNorm.SoundName);
        PriorityNorm.SoundSource     = sharedPref.getString( "priority_norm:sound_source",      PriorityNorm.SoundSource);
        PriorityNorm.LEDColor        = sharedPref.getInt(    "priority_norm:led_color",         PriorityNorm.LEDColor);

        PriorityHigh.EnableLED       = sharedPref.getBoolean("priority_high:enabled_led",       PriorityHigh.EnableLED);
        PriorityHigh.EnableSound     = sharedPref.getBoolean("priority_high:enabled_sound",     PriorityHigh.EnableSound);
        PriorityHigh.EnableVibration = sharedPref.getBoolean("priority_high:enabled_vibration", PriorityHigh.EnableVibration);
        PriorityHigh.RepeatSound     = sharedPref.getBoolean("priority_high:repeat_sound",      PriorityHigh.RepeatSound);
        PriorityHigh.SoundName       = sharedPref.getString( "priority_high:sound_name",        PriorityHigh.SoundName);
        PriorityHigh.SoundSource     = sharedPref.getString( "priority_high:sound_source",      PriorityHigh.SoundSource);
        PriorityHigh.LEDColor        = sharedPref.getInt(    "priority_high:led_color",         PriorityHigh.LEDColor);
    }

    public void save()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("Config", Context.MODE_PRIVATE);
        SharedPreferences.Editor e = sharedPref.edit();

        e.putInt(    "quota_curr",                      quota_curr);
        e.putInt(    "quota_max",                       quota_max);
        e.putInt(    "user_id",                         user_id);
        e.putString( "user_key",                        user_key);
        e.putString( "fcm_token_local",                 fcm_token_local);
        e.putString( "fcm_token_server",                fcm_token_server);

        e.putBoolean("app_enabled",                     Enabled);
        e.putInt(    "local_cache_size",                LocalCacheSize);

        e.putBoolean("priority_low:enabled_led",        PriorityLow.EnableLED);
        e.putBoolean("priority_low:enabled_sound",      PriorityLow.EnableSound);
        e.putBoolean("priority_low:enabled_vibration",  PriorityLow.EnableVibration);
        e.putBoolean("priority_low:repeat_sound",       PriorityLow.RepeatSound);
        e.putString( "priority_low:sound_name",         PriorityLow.SoundName);
        e.putString( "priority_low:sound_source",       PriorityLow.SoundSource);
        e.putInt(    "priority_low:led_color",          PriorityLow.LEDColor);

        e.putBoolean("priority_norm:enabled_led",       PriorityNorm.EnableLED);
        e.putBoolean("priority_norm:enabled_sound",     PriorityNorm.EnableSound);
        e.putBoolean("priority_norm:enabled_vibration", PriorityNorm.EnableVibration);
        e.putBoolean("priority_norm:repeat_sound",      PriorityNorm.RepeatSound);
        e.putString( "priority_norm:sound_name",        PriorityNorm.SoundName);
        e.putString( "priority_norm:sound_source",      PriorityNorm.SoundSource);
        e.putInt(    "priority_norm:led_color",         PriorityNorm.LEDColor);

        e.putBoolean("priority_high:enabled_led",       PriorityHigh.EnableLED);
        e.putBoolean("priority_high:enabled_sound",     PriorityHigh.EnableSound);
        e.putBoolean("priority_high:enabled_vibration", PriorityHigh.EnableVibration);
        e.putBoolean("priority_high:repeat_sound",      PriorityHigh.RepeatSound);
        e.putString( "priority_high:sound_name",        PriorityHigh.SoundName);
        e.putString( "priority_high:sound_source",      PriorityHigh.SoundSource);
        e.putInt(    "priority_high:led_color",         PriorityHigh.LEDColor);

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
