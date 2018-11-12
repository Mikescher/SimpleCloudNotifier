package com.blackforestbytes.simplecloudnotifier.service;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.media.AudioAttributes;
import android.media.AudioManager;
import android.media.Ringtone;
import android.media.RingtoneManager;
import android.net.Uri;
import android.os.Build;
import android.os.VibrationEffect;
import android.os.Vibrator;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.NotificationSettings;
import com.blackforestbytes.simplecloudnotifier.model.PriorityEnum;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

import androidx.core.app.NotificationCompat;

public class NotificationService
{
    private final static String CHANNEL_ID_LOW  = "CHAN_BFB_SCN_MESSAGES_LOW";
    private final static String CHANNEL_ID_NORM = "CHAN_BFB_SCN_MESSAGES_NORM";
    private final static String CHANNEL_ID_HIGH = "CHAN_BFB_SCN_MESSAGES_HIGH";

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
        updateChannels();
    }

    public void updateChannels()
    {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) return;

        Context ctxt = SCNApp.getContext();
        NotificationManager notifman = ctxt.getSystemService(NotificationManager.class);
        if (notifman == null) return;

        NotificationChannel channelLow  = notifman.getNotificationChannel(CHANNEL_ID_LOW);
        if (channelLow == null) notifman.createNotificationChannel(channelLow = new NotificationChannel(CHANNEL_ID_LOW,    "Push notifications (0 - low)", NotificationManager.IMPORTANCE_LOW));
        NotificationChannel channelNorm = notifman.getNotificationChannel(CHANNEL_ID_NORM);
        if (channelNorm == null) notifman.createNotificationChannel(channelNorm = new NotificationChannel(CHANNEL_ID_NORM, "Push notifications (1 - normal)", NotificationManager.IMPORTANCE_DEFAULT));
        NotificationChannel channelHigh = notifman.getNotificationChannel(CHANNEL_ID_HIGH);
        if (channelHigh == null) notifman.createNotificationChannel(channelHigh = new NotificationChannel(CHANNEL_ID_HIGH, "Push notifications (2 - high)", NotificationManager.IMPORTANCE_HIGH));

        channelLow.setDescription("Messages from the API with priority set to low");
        updateSingleChannel(channelLow, SCNSettings.inst().PriorityLow, PriorityEnum.LOW);

        channelNorm.setDescription("Messages from the API with priority set to normal");
        updateSingleChannel(channelNorm, SCNSettings.inst().PriorityNorm, PriorityEnum.NORMAL);

        channelHigh.setDescription("Messages from the API with priority set to high");
        updateSingleChannel(channelHigh, SCNSettings.inst().PriorityHigh, PriorityEnum.HIGH);
    }

    private void updateSingleChannel(NotificationChannel c, NotificationSettings s, PriorityEnum p)
    {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) return;

        c.enableLights(s.EnableLED);
        c.setLightColor(s.LEDColor);
        c.enableVibration(s.EnableVibration);
        c.setVibrationPattern(new long[]{500});

        if (s.EnableSound)
        {
            AudioAttributes attr = new AudioAttributes.Builder()
                    .setContentType(AudioAttributes.CONTENT_TYPE_SONIFICATION)
                    .setLegacyStreamType(AudioManager.STREAM_NOTIFICATION)
                    .setUsage(AudioAttributes.USAGE_NOTIFICATION_EVENT)
                    .build();

            c.setSound(Uri.parse(s.SoundSource), attr);
        }
        else
        {
            c.setSound(null, null);
        }

        if (p == PriorityEnum.HIGH) c.setBypassDnd(true);
        if (p == PriorityEnum.HIGH) c.setLockscreenVisibility(Notification.VISIBILITY_PUBLIC);
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

            if (ns.EnableSound && !ns.SoundSource.isEmpty())
            {
                Ringtone rt = RingtoneManager.getRingtone(SCNApp.getContext(), Uri.parse(ns.SoundSource));
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.P) rt.setLooping(false);
                rt.play();
                new Thread(() -> { try { Thread.sleep(2*1000); } catch (InterruptedException e) { /* */ } rt.stop(); }).start();
            }

            if (ns.EnableVibration)
            {
                Vibrator v = (Vibrator) SCNApp.getContext().getSystemService(Context.VIBRATOR_SERVICE);
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                    v.vibrate(VibrationEffect.createOneShot(500, VibrationEffect.DEFAULT_AMPLITUDE));
                } else {
                    v.vibrate(500);
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

        String channel = CHANNEL_ID_NORM;
        NotificationSettings ns = SCNSettings.inst().PriorityNorm;
        switch (msg.Priority)
        {
            case LOW:    ns = SCNSettings.inst().PriorityLow;  channel = CHANNEL_ID_LOW;  break;
            case NORMAL: ns = SCNSettings.inst().PriorityNorm; channel = CHANNEL_ID_NORM; break;
            case HIGH:   ns = SCNSettings.inst().PriorityHigh; channel = CHANNEL_ID_HIGH; break;
        }

        NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(ctxt, channel);
        mBuilder.setSmallIcon(R.drawable.ic_bfb);
        mBuilder.setContentTitle(msg.Title);
        mBuilder.setContentText(msg.Content);
        mBuilder.setShowWhen(true);
        mBuilder.setWhen(msg.Timestamp);
        mBuilder.setAutoCancel(true);
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O)
        {
            if (msg.Priority == PriorityEnum.LOW) mBuilder.setPriority(NotificationCompat.PRIORITY_LOW);
            if (msg.Priority == PriorityEnum.NORMAL) mBuilder.setPriority(NotificationCompat.PRIORITY_DEFAULT);
            if (msg.Priority == PriorityEnum.HIGH) mBuilder.setPriority(NotificationCompat.PRIORITY_HIGH);
            if (ns.EnableVibration) mBuilder.setVibrate(new long[]{500});
            if (ns.EnableLED) mBuilder.setLights(ns.LEDColor, 500, 500);
        }

        if (ns.EnableSound && !ns.SoundSource.isEmpty())
        {
            mBuilder.setSound(Uri.parse(ns.SoundSource), AudioManager.STREAM_NOTIFICATION);
        }

        Intent intent = new Intent(ctxt, MainActivity.class);
        PendingIntent pi = PendingIntent.getActivity(ctxt, 0, intent, 0);
        mBuilder.setContentIntent(pi);
        NotificationManager mNotificationManager = (NotificationManager) ctxt.getSystemService(Context.NOTIFICATION_SERVICE);

        Notification n = mBuilder.build();
        if (ns.EnableSound && !ns.SoundSource.isEmpty() && ns.RepeatSound) n.flags |= Notification.FLAG_INSISTENT;

        if (mNotificationManager != null) mNotificationManager.notify(0, n);
    }
}
