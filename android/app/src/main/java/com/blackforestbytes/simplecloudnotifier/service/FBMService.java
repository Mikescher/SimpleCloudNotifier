package com.blackforestbytes.simplecloudnotifier.service;

import android.util.Log;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.PriorityEnum;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
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
            Log.i("FB::MessageReceived", "From: " + remoteMessage.getFrom());
            Log.i("FB::MessageReceived", "Payload: " + remoteMessage.getData());
            if (remoteMessage.getNotification() != null) Log.i("FB::MessageReceived", "Notify_Title: " + remoteMessage.getNotification().getTitle());
            if (remoteMessage.getNotification() != null) Log.i("FB::MessageReceived", "Notify_Body: " + remoteMessage.getNotification().getBody());

            long time         = Long.parseLong(remoteMessage.getData().get("timestamp"));
            String title      = remoteMessage.getData().get("title");
            String content    = remoteMessage.getData().get("body");
            PriorityEnum prio = PriorityEnum.parseAPI(remoteMessage.getData().get("priority"));

            CMessage msg = CMessageList.inst().add(time, title, content, prio);


            if (SCNApp.isBackground())
            {
                NotificationService.inst().show(msg);
            }
            else
            {
                SCNApp.showToast("Message recieved: " + title, Toast.LENGTH_LONG);
            }
        }
        catch (Exception e)
        {
            Log.e("FB:Err", e.toString());
            SCNApp.showToast("Recieved invalid message from server", Toast.LENGTH_LONG);
        }
    }
}