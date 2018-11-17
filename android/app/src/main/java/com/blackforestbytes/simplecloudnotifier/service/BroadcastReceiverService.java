package com.blackforestbytes.simplecloudnotifier.service;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;

import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

public class BroadcastReceiverService extends BroadcastReceiver
{
    public static final int NOTIF_SHOW_MAIN  = 10021;
    public static final int NOTIF_STOP_SOUND = 10022;
    public static final String ID_KEY = "com.blackforestbytes.simplecloudnotifier.BroadcastID";

    @Override
    public void onReceive(Context context, Intent intent)
    {
        if (intent == null) return;
        Bundle extras = intent.getExtras();
        if (extras == null) return;
        int notificationId = extras.getInt(ID_KEY, 0);

        if (notificationId == 0) return;
        else if (notificationId == NOTIF_SHOW_MAIN) showMain(context);
        else if (notificationId == NOTIF_STOP_SOUND) stopNotificationSound();
        else return;
    }

    private void stopNotificationSound()
    {
        SoundService.stop();
    }

    private void showMain(Context ctxt)
    {
        SoundService.stop();

        Intent intent = new Intent(ctxt, MainActivity.class);
        ctxt.startActivity(intent);
    }
}
