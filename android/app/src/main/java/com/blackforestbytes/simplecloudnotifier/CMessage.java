package com.blackforestbytes.simplecloudnotifier;

import android.annotation.SuppressLint;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;
import java.util.TimeZone;

public class CMessage
{
    public final long Timestamp ;
    public final String Title;
    public final String Content;

    private static final SimpleDateFormat _format;
    static
    {
        _format = new SimpleDateFormat("yyyy'-'MM'-'dd HH':'mm':'ss", Locale.getDefault());
        _format.setTimeZone(TimeZone.getDefault());
    }

    public CMessage(long t, String mt, String mc)
    {
        Timestamp = t;
        Title = mt;
        Content = mc;
    }

    @SuppressLint("SimpleDateFormat")
    public String formatTimestamp()
    {
        return _format.format(new Date(Timestamp*1000));
    }
}
