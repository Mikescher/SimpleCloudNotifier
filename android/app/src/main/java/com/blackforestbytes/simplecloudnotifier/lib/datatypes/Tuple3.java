package com.blackforestbytes.simplecloudnotifier.lib.datatypes;

public class Tuple3<T1, T2, T3>
{
	public final T1 Item1;
	public final T2 Item2;
	public final T3 Item3;

	public Tuple3(T1 i1, T2 i2, T3 i3)
	{
		Item1 = i1;
		Item2 = i2;
		Item3 = i3;
	}
}
