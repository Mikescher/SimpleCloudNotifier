package com.blackforestbytes.simplecloudnotifier.lib.datatypes;

public class Tuple4<T1, T2, T3, T4>
{
	public final T1 Item1;
	public final T2 Item2;
	public final T3 Item3;
	public final T4 Item4;

	public Tuple4(T1 i1, T2 i2, T3 i3, T4 i4)
	{
		Item1 = i1;
		Item2 = i2;
		Item3 = i3;
		Item4 = i4;
	}
}