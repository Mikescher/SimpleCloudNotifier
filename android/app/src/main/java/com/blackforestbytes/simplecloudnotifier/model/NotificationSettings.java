package com.blackforestbytes.simplecloudnotifier.model;

import android.graphics.Color;
import android.net.Uri;

public class NotificationSettings
{
    public boolean EnableSound = false;
    public String SoundName = "";
    public String SoundSource = Uri.EMPTY.toString();
    public boolean RepeatSound = false;

    public boolean EnableLED = false;
    public int LEDColor = Color.BLUE;

    public boolean EnableVibration = false;
}
