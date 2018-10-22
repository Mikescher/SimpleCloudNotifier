package com.blackforestbytes.simplecloudnotifier.service;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.media.AudioAttributes;
import android.media.AudioManager;
import android.net.Uri;
import android.os.Build;
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
        if (channelLow == null) notifman.createNotificationChannel(channelLow = new NotificationChannel(CHANNEL_ID_LOW, "Push notifications (low priority)", NotificationManager.IMPORTANCE_LOW));
        NotificationChannel channelNorm = notifman.getNotificationChannel(CHANNEL_ID_NORM);
        if (channelNorm == null) notifman.createNotificationChannel(channelNorm = new NotificationChannel(CHANNEL_ID_NORM, "Push notifications (normal priority)", NotificationManager.IMPORTANCE_DEFAULT));
        NotificationChannel channelHigh = notifman.getNotificationChannel(CHANNEL_ID_HIGH);
        if (channelHigh == null) notifman.createNotificationChannel(channelHigh = new NotificationChannel(CHANNEL_ID_HIGH, "Push notifications (high priority)", NotificationManager.IMPORTANCE_HIGH));

        channelLow.setDescription("Messages from the API with priority set to low");
        channelLow.enableLights(SCNSettings.inst().PriorityLow.EnableLED);
        channelLow.setLightColor(SCNSettings.inst().PriorityLow.LEDColor);
        channelLow.enableVibration(SCNSettings.inst().PriorityLow.EnableVibration);
        channelLow.setVibrationPattern(new long[]{200});

        channelNorm.setDescription("Messages from the API with priority set to normal");
        channelNorm.enableLights(SCNSettings.inst().PriorityNorm.EnableLED);
        channelNorm.setLightColor(SCNSettings.inst().PriorityNorm.LEDColor);
        channelNorm.enableVibration(SCNSettings.inst().PriorityNorm.EnableVibration);
        channelNorm.setVibrationPattern(new long[]{200});

        channelHigh.setDescription("Messages from the API with priority set to high");
        channelHigh.enableLights(SCNSettings.inst().PriorityHigh.EnableLED);
        channelHigh.setLightColor(SCNSettings.inst().PriorityHigh.LEDColor);
        channelHigh.enableVibration(SCNSettings.inst().PriorityHigh.EnableVibration);
        channelHigh.setVibrationPattern(new long[]{200});
        channelLow.setBypassDnd(true);
        channelLow.setLockscreenVisibility(Notification.VISIBILITY_PUBLIC);
    }

    public void showForeground(CMessage msg)
    {
        SCNApp.showToast("Message recieved: " + msg.Title, Toast.LENGTH_LONG);
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
            if (ns.EnableVibration) mBuilder.setVibrate(new long[]{200});
            if (ns.EnableLED) mBuilder.setLights(ns.LEDColor, 500, 500);
        }

        if (ns.EnableSound && !ns.SoundSource.isEmpty())
        {
            mBuilder.setSound(Uri.parse(ns.SoundSource), AudioManager.STREAM_ALARM);
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
