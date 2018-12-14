package com.blackforestbytes.simplecloudnotifier.service;

import android.util.Log;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.datatypes.Tuple4;
import com.blackforestbytes.simplecloudnotifier.lib.datatypes.Tuple5;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.LogLevel;
import com.blackforestbytes.simplecloudnotifier.model.PriorityEnum;
import com.blackforestbytes.simplecloudnotifier.model.QueryLog;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.model.ServerCommunication;
import com.blackforestbytes.simplecloudnotifier.model.SingleQuery;
import com.google.android.gms.common.util.JsonUtils;
import com.google.firebase.messaging.FirebaseMessagingService;
import com.google.firebase.messaging.RemoteMessage;

import org.joda.time.Instant;
import org.json.JSONObject;

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
            boolean trimmed   = Boolean.parseBoolean(remoteMessage.getData().get("trimmed"));


            SingleQuery q = new SingleQuery(LogLevel.INFO, Instant.now(), "FBM<recieve>", Str.Empty, new JSONObject(remoteMessage.getData()).toString(), 0, "SUCCESS");
            QueryLog.instance().add(q);

            if (trimmed)
            {
                ServerCommunication.expand(SCNSettings.inst().user_id, SCNSettings.inst().user_key, scn_id, null, (i1, i2, i3, i4, i5) -> recieveData(i4, i1, i2, i3, i5, false));
            }
            else
            {
                recieveData(time, title, content, prio, scn_id, false);
            }
        }
        catch (Exception e)
        {
            Log.e("FB:Err", e.toString());
            SCNApp.showToast("Recieved invalid message from server", Toast.LENGTH_LONG);
        }
    }

    public static void recieveData(long time, String title, String content, PriorityEnum prio, long scn_id, boolean alwaysAck)
    {
        if (CMessageList.inst().isAck(scn_id))
        {
            Log.w("FB::MessageReceived", "Recieved ack-ed message: " + scn_id);
            if (alwaysAck) ServerCommunication.ack(SCNSettings.inst().user_id, SCNSettings.inst().user_key, scn_id);
            return;
        }

        CMessage msg = CMessageList.inst().add(scn_id, time, title, content, prio);

        if (SCNApp.isBackground())
        {
            NotificationService.inst().showBackground(msg);
        }
        else
        {
            NotificationService.inst().showForeground(msg);
        }

        ServerCommunication.ack(SCNSettings.inst().user_id, SCNSettings.inst().user_key, scn_id);
    }
}