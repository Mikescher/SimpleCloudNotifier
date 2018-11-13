package com.blackforestbytes.simplecloudnotifier.service;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.media.AudioManager;
import android.net.Uri;
import android.os.Build;
import android.os.VibrationEffect;
import android.os.Vibrator;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.NotificationSettings;
import com.blackforestbytes.simplecloudnotifier.model.PriorityEnum;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.model.SoundService;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

import androidx.annotation.RequiresApi;
import androidx.core.app.NotificationCompat;

public class NotificationService
{
    private final static String CHANNEL_ID  = "CHAN_BFB_SCN_MESSAGES";

    private final static Object _lock = new Object();
    private static NotificationService _inst = null;
    public static NotificationService inst()
    {
        synchronized (_lock)
        {
            if (_inst != null) return _inst;
            return _inst = new NotificationService();
        }
    }

    private NotificationService()
    {
        createChannels();
    }

    private void createChannels()
    {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) return;

        Context ctxt = SCNApp.getContext();
        NotificationManager notifman = ctxt.getSystemService(NotificationManager.class);
        if (notifman == null) return;

        NotificationChannel channel  = notifman.getNotificationChannel(CHANNEL_ID);
        if (channel == null)
        {
            channel = new NotificationChannel(CHANNEL_ID, "Push notifications", NotificationManager.IMPORTANCE_DEFAULT);
            channel.setDescription("Push notifications from the server");
            channel.setSound(null, null);
            channel.setVibrationPattern(null);
            notifman.createNotificationChannel(channel);
        }
    }

    public void showForeground(CMessage msg)
    {
        SCNApp.showToast("Message recieved: " + msg.Title, Toast.LENGTH_LONG);

        try
        {
            NotificationSettings ns = SCNSettings.inst().PriorityNorm;
            switch (msg.Priority)
            {
                case LOW:    ns = SCNSettings.inst().PriorityLow;  break;
                case NORMAL: ns = SCNSettings.inst().PriorityNorm; break;
                case HIGH:   ns = SCNSettings.inst().PriorityHigh; break;
            }

            SoundService.playForegroundNoLooping(ns.EnableSound, ns.SoundSource, ns.ForceVolume, ns.ForceVolumeValue);

            if (ns.EnableVibration)
            {
                Vibrator v = (Vibrator) SCNApp.getContext().getSystemService(Context.VIBRATOR_SERVICE);
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                    v.vibrate(VibrationEffect.createOneShot(1500, VibrationEffect.DEFAULT_AMPLITUDE));
                } else {
                    v.vibrate(1500);
                }
            }
        }
        catch (Exception e)
        {
            e.printStackTrace();
        }
    }

    public void showBackground(CMessage msg)
    {
        Context ctxt = SCNApp.getContext();

        NotificationSettings ns = SCNSettings.inst().PriorityNorm;
        switch (msg.Priority)
        {
            case LOW:    ns = SCNSettings.inst().PriorityLow;  break;
            case NORMAL: ns = SCNSettings.inst().PriorityNorm; break;
            case HIGH:   ns = SCNSettings.inst().PriorityHigh; break;
        }

        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O)
        {
            // old

            showBackground_old(msg, ctxt, ns);
        }
        else
        {
            // new

            showBackground_new(msg, ctxt, ns);
        }
    }

    private void showBackground_old(CMessage msg, Context ctxt, NotificationSettings ns) {
        NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(ctxt, CHANNEL_ID);
        mBuilder.setSmallIcon(R.drawable.ic_bfb);
        mBuilder.setContentTitle(msg.Title);
        mBuilder.setContentText(msg.Content);
        mBuilder.setShowWhen(true);
        mBuilder.setWhen(msg.Timestamp * 1000);
        mBuilder.setAutoCancel(true);

        if (msg.Priority == PriorityEnum.LOW)    mBuilder.setPriority(NotificationCompat.PRIORITY_LOW);
        if (msg.Priority == PriorityEnum.NORMAL) mBuilder.setPriority(NotificationCompat.PRIORITY_DEFAULT);
        if (msg.Priority == PriorityEnum.HIGH)   mBuilder.setPriority(NotificationCompat.PRIORITY_HIGH);
        if (ns.EnableVibration) mBuilder.setVibrate(new long[]{500});
        if (ns.EnableLED) mBuilder.setLights(ns.LEDColor, 500, 500);

        if (ns.EnableSound && !ns.SoundSource.isEmpty()) mBuilder.setSound(Uri.parse(ns.SoundSource), AudioManager.STREAM_NOTIFICATION);

        Intent intent = new Intent(ctxt, MainActivity.class);
        PendingIntent pi = PendingIntent.getActivity(ctxt, 0, intent, 0);
        mBuilder.setContentIntent(pi);
        NotificationManager mNotificationManager = (NotificationManager) ctxt.getSystemService(Context.NOTIFICATION_SERVICE);

        Notification n = mBuilder.build();
        if (ns.EnableSound && !ns.SoundSource.isEmpty() && ns.RepeatSound) n.flags |= Notification.FLAG_INSISTENT;

        if (mNotificationManager != null) mNotificationManager.notify(0, n);
    }

    @RequiresApi(api = Build.VERSION_CODES.O)
    private void showBackground_new(CMessage msg, Context ctxt, NotificationSettings ns) {
        NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(ctxt, CHANNEL_ID);
        mBuilder.setSmallIcon(R.drawable.ic_bfb);
        mBuilder.setContentTitle(msg.Title);
        mBuilder.setContentText(msg.Content);
        mBuilder.setShowWhen(true);
        mBuilder.setWhen(msg.Timestamp * 1000);
        mBuilder.setAutoCancel(true);

        if (ns.EnableLED) mBuilder.setLights(ns.LEDColor, 500, 500);

        if (msg.Priority == PriorityEnum.LOW)    mBuilder.setPriority(NotificationCompat.PRIORITY_LOW);
        if (msg.Priority == PriorityEnum.NORMAL) mBuilder.setPriority(NotificationCompat.PRIORITY_DEFAULT);
        if (msg.Priority == PriorityEnum.HIGH)   mBuilder.setPriority(NotificationCompat.PRIORITY_HIGH);

        if (ns.ForceVolume)
        {
            AudioManager aman = (AudioManager) SCNApp.getContext().getSystemService(Context.AUDIO_SERVICE);
            int maxVolume = aman.getStreamMaxVolume(AudioManager.STREAM_NOTIFICATION);
            aman.setStreamVolume(AudioManager.STREAM_MUSIC, (int)(maxVolume * (ns.ForceVolumeValue / 100.0)), 0);
        }

        Intent intent = new Intent(ctxt, MainActivity.class);
        PendingIntent pi = PendingIntent.getActivity(ctxt, 0, intent, 0);
        mBuilder.setContentIntent(pi);
        NotificationManager mNotificationManager = (NotificationManager) ctxt.getSystemService(Context.NOTIFICATION_SERVICE);
        if (mNotificationManager == null) return;

        Notification n = mBuilder.build();
        n.flags |= Notification.FLAG_AUTO_CANCEL;

        mNotificationManager.notify(0, n);

        if (ns.EnableSound && !Str.isNullOrWhitespace(ns.SoundSource))
        {
            if (ns.RepeatSound)
            {
                //TODO
            }
            else
            {
                SoundService.playForegroundNoLooping(ns.EnableSound, ns.SoundSource, ns.ForceVolume, ns.ForceVolumeValue);
            }
        }

        if (ns.EnableVibration)
        {
            Vibrator v = (Vibrator) SCNApp.getContext().getSystemService(Context.VIBRATOR_SERVICE);
            v.vibrate(VibrationEffect.createOneShot(1500, VibrationEffect.DEFAULT_AMPLITUDE));
        }

        if (ns.EnableLED)
        {
            //TODO
        }
    }

}
