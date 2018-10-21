package com.blackforestbytes.simplecloudnotifier.model;

import android.graphics.Color;

public class NotificationSettings
{
    public boolean EnableSound = false;
    public String SoundSource = "";
    public boolean RepeatSound = false;

    public boolean EnableLED = false;
    public int LEDColor = Color.BLUE;

    public boolean EnableVibration = false;
}
