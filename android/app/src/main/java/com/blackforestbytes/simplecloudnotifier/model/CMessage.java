package com.blackforestbytes.simplecloudnotifier.model;

import android.annotation.SuppressLint;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;
import java.util.TimeZone;

public class CMessage
{
    public boolean IsExpandedInAdapter = false;

    public final long SCN_ID;
    public final long Timestamp;
    public final String Title;
    public final String Content;
    public final PriorityEnum Priority;

    private static final SimpleDateFormat _format;
    static
    {
        _format = new SimpleDateFormat("yyyy'-'MM'-'dd HH':'mm':'ss", Locale.getDefault());
        _format.setTimeZone(TimeZone.getDefault());
    }

    public CMessage(long id, long t, String mt, String mc, PriorityEnum p)
    {
        SCN_ID = id;
        Timestamp = t;
        Title = mt;
        Content = mc;
        Priority = p;
    }

    @SuppressLint("SimpleDateFormat")
    public String formatTimestamp()
    {
        return _format.format(new Date(Timestamp*1000));
    }
}
