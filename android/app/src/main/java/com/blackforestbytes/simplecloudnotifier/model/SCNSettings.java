package com.blackforestbytes.simplecloudnotifier.model;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;
import android.view.View;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.datatypes.Tuple3;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.google.firebase.installations.FirebaseInstallations;

public class SCNSettings
{
    private final static Object _lock = new Object();
    private static volatile SCNSettings _inst = null;
    public static SCNSettings inst()
    {
        SCNSettings local = _inst;
        if (local == null)
        {
            synchronized (_lock)
            {
                local = _inst;
                if (local == null) _inst = local = new SCNSettings();
            }
        }
        return local;
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

    public String  promode_token;
    public boolean promode_local;
    public boolean promode_server;

    // ------------------------------------------------------------

    public boolean Enabled = true;
    public int LocalCacheSize = 500;
    public boolean EnableDeleteSwipe = false;
    public int PreviewLineCount = 6;

    public final NotificationSettings PriorityLow  = new NotificationSettings(PriorityEnum.LOW);
    public final NotificationSettings PriorityNorm = new NotificationSettings(PriorityEnum.NORMAL);
    public final NotificationSettings PriorityHigh = new NotificationSettings(PriorityEnum.HIGH);

    // ------------------------------------------------------------

    public SCNSettings()
    {
        reloadPrefs();
    }

    public void reloadPrefs()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("Config", Context.MODE_PRIVATE);

        quota_curr       = sharedPref.getInt(   "quota_curr",        0);
        quota_max        = sharedPref.getInt(   "quota_max",         0);
        user_id          = sharedPref.getInt(   "user_id",          -1);
        user_key         = sharedPref.getString("user_key",         "");
        fcm_token_local  = sharedPref.getString("fcm_token_local",  "");
        fcm_token_server = sharedPref.getString("fcm_token_server", "");
        promode_local    = sharedPref.getBoolean("promode_local",   false);
        promode_server   = sharedPref.getBoolean("promode_server",  false);
        promode_token    = sharedPref.getString("promode_token",    "");

        Enabled                     = sharedPref.getBoolean("app_enabled",  Enabled);
        LocalCacheSize              = sharedPref.getInt("local_cache_size", LocalCacheSize);
        EnableDeleteSwipe           = sharedPref.getBoolean("do_del_swipe",  EnableDeleteSwipe);
        PreviewLineCount            = sharedPref.getInt("preview_line_count", PreviewLineCount);

        PriorityLow.EnableLED         = sharedPref.getBoolean("priority_low:enabled_led",         PriorityLow.EnableLED);
        PriorityLow.EnableSound       = sharedPref.getBoolean("priority_low:enabled_sound",       PriorityLow.EnableSound);
        PriorityLow.EnableVibration   = sharedPref.getBoolean("priority_low:enabled_vibration",   PriorityLow.EnableVibration);
        PriorityLow.RepeatSound       = sharedPref.getBoolean("priority_low:repeat_sound",        PriorityLow.RepeatSound);
        PriorityLow.SoundName         = sharedPref.getString( "priority_low:sound_name",          PriorityLow.SoundName);
        PriorityLow.SoundSource       = sharedPref.getString( "priority_low:sound_source",        PriorityLow.SoundSource);
        PriorityLow.LEDColor          = sharedPref.getInt(    "priority_low:led_color",           PriorityLow.LEDColor);
        PriorityLow.ForceVolume       = sharedPref.getBoolean("priority_low:force_volume",        PriorityLow.ForceVolume);
        PriorityLow.ForceVolumeValue  = sharedPref.getInt(    "priority_low:force_volume_value",  PriorityLow.ForceVolumeValue);

        PriorityNorm.EnableLED        = sharedPref.getBoolean("priority_norm:enabled_led",        PriorityNorm.EnableLED);
        PriorityNorm.EnableSound      = sharedPref.getBoolean("priority_norm:enabled_sound",      PriorityNorm.EnableSound);
        PriorityNorm.EnableVibration  = sharedPref.getBoolean("priority_norm:enabled_vibration",  PriorityNorm.EnableVibration);
        PriorityNorm.RepeatSound      = sharedPref.getBoolean("priority_norm:repeat_sound",       PriorityNorm.RepeatSound);
        PriorityNorm.SoundName        = sharedPref.getString( "priority_norm:sound_name",         PriorityNorm.SoundName);
        PriorityNorm.SoundSource      = sharedPref.getString( "priority_norm:sound_source",       PriorityNorm.SoundSource);
        PriorityNorm.LEDColor         = sharedPref.getInt(    "priority_norm:led_color",          PriorityNorm.LEDColor);
        PriorityNorm.ForceVolume      = sharedPref.getBoolean("priority_norm:force_volume",       PriorityNorm.ForceVolume);
        PriorityNorm.ForceVolumeValue = sharedPref.getInt(    "priority_norm:force_volume_value", PriorityNorm.ForceVolumeValue);

        PriorityHigh.EnableLED        = sharedPref.getBoolean("priority_high:enabled_led",        PriorityHigh.EnableLED);
        PriorityHigh.EnableSound      = sharedPref.getBoolean("priority_high:enabled_sound",      PriorityHigh.EnableSound);
        PriorityHigh.EnableVibration  = sharedPref.getBoolean("priority_high:enabled_vibration",  PriorityHigh.EnableVibration);
        PriorityHigh.RepeatSound      = sharedPref.getBoolean("priority_high:repeat_sound",       PriorityHigh.RepeatSound);
        PriorityHigh.SoundName        = sharedPref.getString( "priority_high:sound_name",         PriorityHigh.SoundName);
        PriorityHigh.SoundSource      = sharedPref.getString( "priority_high:sound_source",       PriorityHigh.SoundSource);
        PriorityHigh.LEDColor         = sharedPref.getInt(    "priority_high:led_color",          PriorityHigh.LEDColor);
        PriorityHigh.ForceVolume      = sharedPref.getBoolean("priority_high:force_volume",       PriorityHigh.ForceVolume);
        PriorityHigh.ForceVolumeValue = sharedPref.getInt(    "priority_high:force_volume_value", PriorityHigh.ForceVolumeValue);
    }

    public void save()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("Config", Context.MODE_PRIVATE);
        SharedPreferences.Editor e = sharedPref.edit();

        e.putInt(    "quota_curr",                       quota_curr);
        e.putInt(    "quota_max",                        quota_max);
        e.putInt(    "user_id",                          user_id);
        e.putString( "user_key",                         user_key);
        e.putString( "fcm_token_local",                  fcm_token_local);
        e.putString( "fcm_token_server",                 fcm_token_server);
        e.putBoolean("promode_local",                    promode_local);
        e.putBoolean("promode_server",                   promode_server);
        e.putString( "promode_token",                    promode_token);

        e.putBoolean("app_enabled",                      Enabled);
        e.putInt(    "local_cache_size",                 LocalCacheSize);
        e.putBoolean("do_del_swipe",                     EnableDeleteSwipe);
        e.putInt(    "preview_line_count",               PreviewLineCount);

        e.putBoolean("priority_low:enabled_led",         PriorityLow.EnableLED);
        e.putBoolean("priority_low:enabled_sound",       PriorityLow.EnableSound);
        e.putBoolean("priority_low:enabled_vibration",   PriorityLow.EnableVibration);
        e.putBoolean("priority_low:repeat_sound",        PriorityLow.RepeatSound);
        e.putString( "priority_low:sound_name",          PriorityLow.SoundName);
        e.putString( "priority_low:sound_source",        PriorityLow.SoundSource);
        e.putInt(    "priority_low:led_color",           PriorityLow.LEDColor);
        e.putBoolean("priority_low:force_volume",        PriorityLow.ForceVolume);
        e.putInt(    "priority_low:force_volume_value",  PriorityLow.ForceVolumeValue);

        e.putBoolean("priority_norm:enabled_led",        PriorityNorm.EnableLED);
        e.putBoolean("priority_norm:enabled_sound",      PriorityNorm.EnableSound);
        e.putBoolean("priority_norm:enabled_vibration",  PriorityNorm.EnableVibration);
        e.putBoolean("priority_norm:repeat_sound",       PriorityNorm.RepeatSound);
        e.putString( "priority_norm:sound_name",         PriorityNorm.SoundName);
        e.putString( "priority_norm:sound_source",       PriorityNorm.SoundSource);
        e.putInt(    "priority_norm:led_color",          PriorityNorm.LEDColor);
        e.putBoolean("priority_norm:force_volume",       PriorityNorm.ForceVolume);
        e.putInt(    "priority_norm:force_volume_value", PriorityNorm.ForceVolumeValue);

        e.putBoolean("priority_high:enabled_led",        PriorityHigh.EnableLED);
        e.putBoolean("priority_high:enabled_sound",      PriorityHigh.EnableSound);
        e.putBoolean("priority_high:enabled_vibration",  PriorityHigh.EnableVibration);
        e.putBoolean("priority_high:repeat_sound",       PriorityHigh.RepeatSound);
        e.putString( "priority_high:sound_name",         PriorityHigh.SoundName);
        e.putString( "priority_high:sound_source",       PriorityHigh.SoundSource);
        e.putInt(    "priority_high:led_color",          PriorityHigh.LEDColor);
        e.putBoolean("priority_high:force_volume",       PriorityHigh.ForceVolume);
        e.putInt(    "priority_high:force_volume_value", PriorityHigh.ForceVolumeValue);

        e.apply();
    }

    public boolean isConnected()
    {
        return user_id>=0 && user_key != null && !user_key.isEmpty();
    }

    public String createOnlineURL(boolean longurl)
    {
        String base = longurl ? ServerCommunication.PAGE_URL_LONG : ServerCommunication.PAGE_URL_SHORT;

        if (!isConnected()) return base;
        return base + "index.php?preset_user_id="+user_id+"&preset_user_key="+user_key;
    }

    public void setServerToken(String token, View loader, boolean force)
    {
        if (isConnected())
        {
            fcm_token_local = token;
            save();
            if (!fcm_token_local.equals(fcm_token_server) || force) ServerCommunication.updateFCMToken(user_id, user_key, fcm_token_local, loader);
        }
        else
        {
            fcm_token_local = token;
            save();
            ServerCommunication.register(fcm_token_local, loader, promode_local, promode_token);
            updateProState(loader);
        }
    }

    // called at app start
    public void work(Activity a, boolean force)
    {
        FirebaseInstallations.getInstance().getId().addOnSuccessListener(a, newToken ->
        {
            Log.d("FB::GetInstanceId", newToken);
            SCNSettings.inst().setServerToken(newToken, null, force);
        }).addOnCompleteListener(r ->
        {
            if (isConnected()) ServerCommunication.info(user_id, user_key, null);
        });

        updateProState(null);
    }

    // reset account key
    public void reset(View loader)
    {
        if (!isConnected()) return;

        ServerCommunication.resetSecret(user_id, user_key, loader);
    }

    // refresh account data
    public void refresh(View loader, Activity a)
    {
        if (isConnected())
        {
            ServerCommunication.info(user_id, user_key, loader);

            if (promode_server != promode_local) updateProState(loader);

            if (!Str.equals(fcm_token_local, fcm_token_server)) work(a, false);
        }
        else
        {
            // get token then register
            FirebaseInstallations.getInstance().getId().addOnSuccessListener(a, newToken ->
            {
                Log.d("FB::GetInstanceId", newToken);
                SCNSettings.inst().setServerToken(newToken, loader, false); // does register in here
            }).addOnCompleteListener(r ->
            {
                if (isConnected()) ServerCommunication.info(user_id, user_key, null); // info again for safety
            });
        }
    }

    public void updateProState(View loader)
    {
        Tuple3<Boolean, Boolean, String> state = IABService.inst().getPurchaseCachedExtended(IABService.IAB_PRO_MODE);
        if (!state.Item2) return; // not initialized

        boolean promode_real = state.Item1;

        if (promode_real != promode_local || promode_real != promode_server)
        {
            promode_local = promode_real;
            promode_token = promode_real ? state.Item3 : "";
            save();

            updateProStateOnServer(loader);
        }
    }

    public void updateProStateOnServer(View loader)
    {
        if (!isConnected()) return;

        ServerCommunication.upgrade(user_id, user_key, loader, promode_local, promode_token);
    }
}
