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

    public boolean ForceVolume;
    public int ForceVolumeValue;

    public boolean EnableLED;
    public int LEDColor;

    public boolean EnableVibration;

    public NotificationSettings(PriorityEnum p)
    {
        EnableSound      = (p == PriorityEnum.HIGH);
        SoundName        = "Default";
        SoundSource      = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION).toString();
        RepeatSound      = false;
        EnableLED        = (p == PriorityEnum.HIGH) || (p == PriorityEnum.NORMAL);
        LEDColor         = Color.BLUE;
        EnableVibration  = (p == PriorityEnum.HIGH) || (p == PriorityEnum.NORMAL);
        ForceVolume      = false;
        ForceVolumeValue = 50;
    }
}
