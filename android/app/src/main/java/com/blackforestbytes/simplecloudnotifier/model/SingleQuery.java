package com.blackforestbytes.simplecloudnotifier.model;

import android.content.SharedPreferences;
import android.os.BaseBundle;
import android.os.Bundle;

import com.blackforestbytes.simplecloudnotifier.lib.string.Str;

import org.joda.time.Instant;

public class SingleQuery
{
    public final Instant Timestamp;

    public final LogLevel Level;
    public final String Name;
    public final String URL;
    public final String Response;
    public final int ResponseCode;
    public final String ExceptionString;

    public SingleQuery(LogLevel l, Instant i, String n, String u, String r, int rc, String e)
    {
        Level=l;
        Timestamp=i;
        Name=n;
        URL=u;
        Response=r;
        ResponseCode=rc;
        ExceptionString=e;
    }

    public void save(SharedPreferences.Editor e, String base)
    {
        e.putInt(base+".Level", Level.asInt());
        e.putLong(base+".Timestamp", Timestamp.getMillis());
        e.putString(base+".Name", Name);
        e.putString(base+".URL", URL);
        e.putString(base+".Response", Response);
        e.putInt(base+".ResponseCode", ResponseCode);
        e.putString(base+".ExceptionString", ExceptionString);
    }

    public void save(BaseBundle e, String base)
    {
        e.putInt(base+".Level", Level.asInt());
        e.putLong(base+".Timestamp", Timestamp.getMillis());
        e.putString(base+".Name", Name);
        e.putString(base+".URL", URL);
        e.putString(base+".Response", Response);
        e.putInt(base+".ResponseCode", ResponseCode);
        e.putString(base+".ExceptionString", ExceptionString);
    }

    public static SingleQuery load(SharedPreferences e, String base)
    {
        return new SingleQuery
                (
                        LogLevel.fromInt(e.getInt(base+".Level", 0)),
                        new Instant(e.getLong(base+".Timestamp", 0)),
                        e.getString(base+".Name", Str.Empty),
                        e.getString(base+".URL", Str.Empty),
                        e.getString(base+".Response", Str.Empty),
                        e.getInt(base+".ResponseCode", -1),
                        e.getString(base+".ExceptionString", Str.Empty)
                );
    }

    public static SingleQuery load(BaseBundle e, String base)
    {
        return new SingleQuery
                (
                        LogLevel.fromInt(e.getInt(base+".Level", 0)),
                        new Instant(e.getLong(base+".Timestamp", 0)),
                        e.getString(base+".Name", Str.Empty),
                        e.getString(base+".URL", Str.Empty),
                        e.getString(base+".Response", Str.Empty),
                        e.getInt(base+".ResponseCode", -1),
                        e.getString(base+".ExceptionString", Str.Empty)
                );
    }
}
