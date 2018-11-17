package com.blackforestbytes.simplecloudnotifier.service;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.graphics.Color;
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
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

import androidx.annotation.RequiresApi;
import androidx.core.app.NotificationCompat;

public class NotificationService
{
    private final static String CHANNEL_P0_ID  = "CHAN_BFB_SCN_MESSAGES_P0";
    private final static String CHANNEL_P1_ID  = "CHAN_BFB_SCN_MESSAGES_P1";
    private final static String CHANNEL_P2_ID  = "CHAN_BFB_SCN_MESSAGES_P2";

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

        {
            NotificationChannel channel0  = notifman.getNotificationChannel(CHANNEL_P0_ID);
            if (channel0 == null)
            {
                channel0 = new NotificationChannel(CHANNEL_P0_ID, "Push notifications (low priority)", NotificationManager.IMPORTANCE_DEFAULT);
                channel0.setDescription("Push notifications from the server with low priority.\nGo to the in-app settings to configure ringtone, volume and vibrations");
                channel0.setSound(null, null);
                channel0.setVibrationPattern(null);
                channel0.setLightColor(Color.BLUE);
                notifman.createNotificationChannel(channel0);
            }
        }
        {
            NotificationChannel channel1  = notifman.getNotificationChannel(CHANNEL_P1_ID);
            if (channel1 == null)
            {
                channel1 = new NotificationChannel(CHANNEL_P1_ID, "Push notifications (normal priority)", NotificationManager.IMPORTANCE_DEFAULT);
                channel1.setDescription("Push notifications from the server with low priority.\nGo to the in-app settings to configure ringtone, volume and vibrations");
                channel1.setSound(null, null);
                channel1.setVibrationPattern(null);
                channel1.setLightColor(Color.BLUE);
                notifman.createNotificationChannel(channel1);
            }
        }
        {
            NotificationChannel channel2  = notifman.getNotificationChannel(CHANNEL_P2_ID);
            if (channel2 == null)
            {
                channel2 = new NotificationChannel(CHANNEL_P1_ID, "Push notifications (high priority)", NotificationManager.IMPORTANCE_DEFAULT);
                channel2.setDescription("Push notifications from the server with low priority.\nGo to the in-app settings to configure ringtone, volume and vibrations");
                channel2.setSound(null, null);
                channel2.setVibrationPattern(null);
                channel2.setLightColor(Color.BLUE);
                notifman.createNotificationChannel(channel2);
            }
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

            SoundService.play(ns.EnableSound, ns.SoundSource, ns.ForceVolume, ns.ForceVolumeValue, false);

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
            showBackground_old(msg, ctxt, ns, msg.Priority);
        }
        else
        {
            // new
            showBackground_new(msg, ctxt, ns, msg.Priority);
        }
    }

    private String getChannel(PriorityEnum p)
    {
        switch (p)
        {
            case LOW:    return CHANNEL_P0_ID;
            case NORMAL: return CHANNEL_P1_ID;
            case HIGH:   return CHANNEL_P2_ID;

            default:     return CHANNEL_P0_ID;
        }
    }

    private void showBackground_old(CMessage msg, Context ctxt, NotificationSettings ns, PriorityEnum prio)
    {
        NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(ctxt, getChannel(prio));
        mBuilder.setSmallIcon(R.drawable.ic_notification_white);
        mBuilder.setLargeIcon(BitmapFactory.decodeResource(ctxt.getResources(), R.mipmap.ic_notification_full));
        mBuilder.setContentTitle(msg.Title);
        mBuilder.setContentText(msg.Content);
        mBuilder.setShowWhen(true);
        mBuilder.setWhen(msg.Timestamp * 1000);
        mBuilder.setAutoCancel(true);
        mBuilder.setCategory(Notification.CATEGORY_MESSAGE);
        mBuilder.setGroup("com.blackforestbytes.simplecloudnotifier.notifications.group."+prio.toString());

        if (msg.Priority == PriorityEnum.LOW)    mBuilder.setPriority(NotificationCompat.PRIORITY_LOW);
        if (msg.Priority == PriorityEnum.NORMAL) mBuilder.setPriority(NotificationCompat.PRIORITY_DEFAULT);
        if (msg.Priority == PriorityEnum.HIGH)   mBuilder.setPriority(NotificationCompat.PRIORITY_HIGH);
        if (ns.EnableVibration) mBuilder.setVibrate(new long[]{500});
        if (ns.EnableLED) mBuilder.setLights(ns.LEDColor, 500, 500);

        if (ns.EnableSound && !ns.SoundSource.isEmpty() && !ns.RepeatSound) mBuilder.setSound(Uri.parse(ns.SoundSource), AudioManager.STREAM_NOTIFICATION);

        Intent intent = new Intent(ctxt, MainActivity.class);
        PendingIntent pi = PendingIntent.getActivity(ctxt, 0, intent, 0);
        mBuilder.setContentIntent(pi);
        NotificationManager mNotificationManager = (NotificationManager) ctxt.getSystemService(Context.NOTIFICATION_SERVICE);

        if (ns.EnableSound && !ns.SoundSource.isEmpty() && ns.RepeatSound)
        {
            Intent intnt_stop = new Intent(SCNApp.getContext(), BroadcastReceiverService.class);
            intnt_stop.putExtra(BroadcastReceiverService.ID_KEY, BroadcastReceiverService.NOTIF_STOP_SOUND);
            PendingIntent pi_stop = PendingIntent.getBroadcast(SCNApp.getContext().getApplicationContext(), BroadcastReceiverService.NOTIF_STOP_SOUND, intnt_stop, 0);
            mBuilder.addAction(new NotificationCompat.Action(-1, "Stop", pi_stop));
            mBuilder.setDeleteIntent(pi_stop);

            SoundService.play(ns.EnableSound, ns.SoundSource, ns.ForceVolume, ns.ForceVolumeValue, ns.RepeatSound);
        }

        Notification n = mBuilder.build();

        if (mNotificationManager != null) mNotificationManager.notify((int)msg.SCN_ID, n);
    }

    @RequiresApi(api = Build.VERSION_CODES.O)
    private void showBackground_new(CMessage msg, Context ctxt, NotificationSettings ns, PriorityEnum prio)
    {
        NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(ctxt, getChannel(prio));
        mBuilder.setSmallIcon(R.drawable.ic_notification_white);
        mBuilder.setLargeIcon(BitmapFactory.decodeResource(ctxt.getResources(), R.mipmap.ic_notification_full));
        mBuilder.setContentTitle(msg.Title);
        mBuilder.setContentText(msg.Content);
        mBuilder.setShowWhen(true);
        mBuilder.setWhen(msg.Timestamp * 1000);
        mBuilder.setAutoCancel(true);
        mBuilder.setCategory(Notification.CATEGORY_MESSAGE);
        mBuilder.setGroup("com.blackforestbytes.simplecloudnotifier.notifications.group."+prio.toString());

        if (ns.EnableLED) mBuilder.setLights(ns.LEDColor, 500, 500);

        if (msg.Priority == PriorityEnum.LOW)    mBuilder.setPriority(NotificationCompat.PRIORITY_LOW);
        if (msg.Priority == PriorityEnum.NORMAL) mBuilder.setPriority(NotificationCompat.PRIORITY_DEFAULT);
        if (msg.Priority == PriorityEnum.HIGH)   mBuilder.setPriority(NotificationCompat.PRIORITY_HIGH);

        Intent intnt_click = new Intent(SCNApp.getContext(), BroadcastReceiverService.class);
        intnt_click.putExtra(BroadcastReceiverService.ID_KEY, BroadcastReceiverService.NOTIF_SHOW_MAIN);
        PendingIntent pi = PendingIntent.getBroadcast(ctxt, 0, intnt_click, 0);
        mBuilder.setContentIntent(pi);
        NotificationManager mNotificationManager = (NotificationManager) ctxt.getSystemService(Context.NOTIFICATION_SERVICE);
        if (mNotificationManager == null) return;

        if (ns.EnableSound && !Str.isNullOrWhitespace(ns.SoundSource))
        {
            if (ns.RepeatSound)
            {
                Intent intnt_stop = new Intent(SCNApp.getContext(), BroadcastReceiverService.class);
                intnt_stop.putExtra(BroadcastReceiverService.ID_KEY, BroadcastReceiverService.NOTIF_STOP_SOUND);
                PendingIntent pi_stop = PendingIntent.getBroadcast(ctxt, BroadcastReceiverService.NOTIF_STOP_SOUND, intnt_stop, 0);
                mBuilder.addAction(new NotificationCompat.Action(-1, "Stop", pi_stop));
                mBuilder.setDeleteIntent(pi_stop);
            }

            SoundService.play(ns.EnableSound, ns.SoundSource, ns.ForceVolume, ns.ForceVolumeValue, ns.RepeatSound);
        }

        Notification n = mBuilder.build();
        n.flags |= Notification.FLAG_AUTO_CANCEL;

        mNotificationManager.notify((int)msg.SCN_ID, n);

        if (ns.EnableVibration)
        {
            Vibrator v = (Vibrator) SCNApp.getContext().getSystemService(Context.VIBRATOR_SERVICE);
            v.vibrate(VibrationEffect.createOneShot(1500, VibrationEffect.DEFAULT_AMPLITUDE));
        }

        //if (ns.EnableLED) {  } // no LED in Android-O -- configure via Channel
    }

}
