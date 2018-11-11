package com.blackforestbytes.simplecloudnotifier.lib.string;

import android.content.Context;
import android.util.Log;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.lambda.Func1to1;

import java.text.MessageFormat;
import java.util.List;

public class Str
{
	public final static String Empty = "";

	public static String format(String fmt, Object... data)
	{
		return MessageFormat.format(fmt, data);
	}

	public static String rformat(int fmtResId, Object... data)
	{
		Context inst = SCNApp.getContext();
		if (inst == null)
		{
			Log.e("StringFormat", "rformat::NoInstance --> inst==null for" + fmtResId);
			return "?ERR?";
		}

		return MessageFormat.format(inst.getResources().getString(fmtResId), data);
	}

	public static String firstLine(String content)
	{
		int idx = content.indexOf('\n');
		if (idx == -1) return content;

		if (idx == 0) return Str.Empty;

		if (content.charAt(idx-1) == '\r') return content.substring(0, idx-1);

		return content.substring(0, idx);
	}

	public static boolean isNullOrWhitespace(String str)
	{
		return str == null || str.length() == 0 || str.trim().length() == 0;
	}

	public static boolean isNullOrEmpty(String str)
	{
		return str == null || str.length() == 0;
	}

	public static boolean equals(String a, String b)
	{
		if (a == null) return (b == null);
		return a.equals(b);
	}

	public static String join(String sep, List<String> list)
	{
		StringBuilder b = new StringBuilder();
		boolean first = true;
		for (String v : list)
		{
			if (!first) b.append(sep);
			b.append(v);
			first = false;
		}
		return b.toString();
	}

	public static <T> String join(String sep, List<T> list, Func1to1<T, String> map)
	{
		StringBuilder b = new StringBuilder();
		boolean first = true;
		for (T v : list)
		{
			if (!first) b.append(sep);
			b.append(map.invoke(v));
			first = false;
		}
		return b.toString();
	}

	public static Integer tryParseToInt(String s)
	{
		try
		{
			return Integer.parseInt(s);
		}
		catch (Exception e)
		{
			return null;
		}
	}
}
