package com.blackforestbytes.simplecloudnotifier.service;

import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.os.Build;
import android.support.v4.app.NotificationCompat;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

public class NotificationService
{
    private final static String CHANNEL_ID = "CHAN_BFB_SCN_MESSAGES";

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
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O)
        {
            Context ctxt = SCNApp.getContext();

            NotificationChannel channel = new NotificationChannel(CHANNEL_ID, "Push notifications", NotificationManager.IMPORTANCE_HIGH);
            channel.setDescription("Messages from the API");
            channel.setLightColor(Color.rgb(255, 0, 0));
            channel.setVibrationPattern(new long[]{200});
            channel.enableLights(true);
            channel.enableVibration(true);

            NotificationManager notificationManager = ctxt.getSystemService(NotificationManager.class);
            if (notificationManager != null) notificationManager.createNotificationChannel(channel);
        }

    }

    public void show(CMessage msg)
    {
        Context ctxt = SCNApp.getContext();

        NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(ctxt, CHANNEL_ID)
                .setSmallIcon(R.drawable.ic_bfb)
                .setContentTitle(msg.Title)
                .setContentText(msg.Content)
                .setShowWhen(true)
                .setWhen(msg.Timestamp)
                .setPriority(NotificationCompat.PRIORITY_HIGH)
                .setAutoCancel(true);
        Intent intent = new Intent(ctxt, MainActivity.class);
        PendingIntent pi = PendingIntent.getActivity(ctxt, 0, intent, 0);
        mBuilder.setContentIntent(pi);
        NotificationManager mNotificationManager = (NotificationManager) ctxt.getSystemService(Context.NOTIFICATION_SERVICE);
        if (mNotificationManager != null) mNotificationManager.notify(0, mBuilder.build());
    }
}
