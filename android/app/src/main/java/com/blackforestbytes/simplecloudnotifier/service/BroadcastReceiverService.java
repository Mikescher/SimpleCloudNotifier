package com.blackforestbytes.simplecloudnotifier.service;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;

public class BroadcastReceiverService extends BroadcastReceiver
{
    public static final int STOP_NOTIFICATION_SOUND = 10022;
    public static final String ID_KEY = "com.blackforestbytes.simplecloudnotifier.BroadcastID";

    @Override
    public void onReceive(Context context, Intent intent)
    {
        if (intent == null) return;
        Bundle extras = intent.getExtras();
        if (extras == null) return;
        int notificationId = extras.getInt(ID_KEY, 0);

        if (notificationId == 10022) stopNotificationSound();
        else return;
    }

    private void stopNotificationSound()
    {
        SoundService.stopPlaying();
    }
}
