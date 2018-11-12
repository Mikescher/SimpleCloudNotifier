package com.blackforestbytes.simplecloudnotifier.model;

import android.graphics.Color;
import android.media.RingtoneManager;
import android.net.Uri;

public class NotificationSettings
{
    public boolean EnableSound;
    public String SoundName;
    public String SoundSource;
    public boolean RepeatSound;

    public boolean EnableLED;
    public int LEDColor;

    public boolean EnableVibration;

    public NotificationSettings(PriorityEnum p)
    {
        EnableSound     = (p == PriorityEnum.HIGH);
        SoundName       = (p == PriorityEnum.HIGH) ? "Default" : "";
        SoundSource     = (p == PriorityEnum.HIGH) ? RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION).toString() : Uri.EMPTY.toString();
        RepeatSound     = false;
        EnableLED       = (p == PriorityEnum.HIGH) || (p == PriorityEnum.NORMAL);
        LEDColor        = Color.BLUE;
        EnableVibration = (p == PriorityEnum.HIGH) || (p == PriorityEnum.NORMAL);
    }
}
