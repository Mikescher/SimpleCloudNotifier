package com.blackforestbytes.simplecloudnotifier.service;

import android.util.Log;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.PriorityEnum;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.model.ServerCommunication;
import com.google.firebase.messaging.FirebaseMessagingService;
import com.google.firebase.messaging.RemoteMessage;

public class FBMService extends FirebaseMessagingService
{
    @Override
    public void onNewToken(String token)
    {
        Log.i("Firebase::NewToken", token);
        SCNSettings.inst().setServerToken(token, null);
    }

    @Override
    public void onMessageReceived(RemoteMessage remoteMessage)
    {
        try
        {
            if (!SCNSettings.inst().Enabled) return;

            Log.i("FB::MessageReceived", "From: " + remoteMessage.getFrom());
            Log.i("FB::MessageReceived", "Payload: " + remoteMessage.getData());
            if (remoteMessage.getNotification() != null) Log.i("FB::MessageReceived", "Notify_Title: " + remoteMessage.getNotification().getTitle());
            if (remoteMessage.getNotification() != null) Log.i("FB::MessageReceived", "Notify_Body: " + remoteMessage.getNotification().getBody());

            long time         = Long.parseLong(remoteMessage.getData().get("timestamp"));
            String title      = remoteMessage.getData().get("title");
            String content    = remoteMessage.getData().get("body");
            PriorityEnum prio = PriorityEnum.parseAPI(remoteMessage.getData().get("priority"));
            long scn_id       = Long.parseLong(remoteMessage.getData().get("scn_msg_id"));

            recieveData(time, title, content, prio, scn_id, false);
        }
        catch (Exception e)
        {
            Log.e("FB:Err", e.toString());
            SCNApp.showToast("Recieved invalid message from server", Toast.LENGTH_LONG);
        }
    }

    public static void recieveData(long time, String title, String content, PriorityEnum prio, long scn_id, boolean alwaysAck)
    {
        CMessage msg = CMessageList.inst().add(scn_id, time, title, content, prio);

        if (CMessageList.inst().isAck(scn_id))
        {
            Log.w("FB::MessageReceived", "Recieved ack-ed message: " + scn_id);
            if (alwaysAck) ServerCommunication.ack(SCNSettings.inst().user_id, SCNSettings.inst().user_key, msg);
            return;
        }

        if (SCNApp.isBackground())
        {
            NotificationService.inst().showBackground(msg);
        }
        else
        {
            NotificationService.inst().showForeground(msg);
        }

        ServerCommunication.ack(SCNSettings.inst().user_id, SCNSettings.inst().user_key, msg);
    }
}