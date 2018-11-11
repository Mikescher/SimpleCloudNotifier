package com.blackforestbytes.simplecloudnotifier.lib.android;

import android.util.Log;

public final class ThreadUtils
{
	public static void safeSleep(int millisMin, int millisMax)
	{
		safeSleep(millisMin + (int)(Math.random()*(millisMax-millisMin)));
	}

	public static void safeSleep(int millis)
	{
		try
		{
			Thread.sleep(millis);
		}
		catch (InterruptedException e)
		{
			Log.d("ThreadUtils", e.toString());
		}
	}
}
