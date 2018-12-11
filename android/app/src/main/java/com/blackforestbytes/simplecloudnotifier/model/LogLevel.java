package com.blackforestbytes.simplecloudnotifier.model;


import android.graphics.Color;

public enum LogLevel
{
    DEBUG,
    INFO,
    WARN,
    ERROR;

    public String toUIString()
    {
        switch (this)
        {
            case DEBUG: return "Debug";
            case INFO:  return "Info";
            case WARN:  return "Warning";
            case ERROR: return "Error";
            default:    return "???";
        }
    }

    public int getColor()
    {
        switch (this)
        {
            case DEBUG: return Color.GRAY;
            case WARN:  return Color.rgb(171, 145, 68);
            case INFO:  return Color.BLACK;
            case ERROR: return Color.RED;
            default:    return Color.MAGENTA;
        }
    }

    public int asInt()
    {
        switch (this)
        {
            case DEBUG: return 0;
            case WARN:  return 1;
            case INFO:  return 2;
            case ERROR: return 3;
            default:    return 999;
        }
    }

    public static LogLevel fromInt(int i)
    {
        if (i == 0) return LogLevel.DEBUG;
        if (i == 1) return LogLevel.WARN;
        if (i == 2) return LogLevel.INFO;
        if (i == 3) return LogLevel.ERROR;

        return LogLevel.ERROR; // ????
    }
}
