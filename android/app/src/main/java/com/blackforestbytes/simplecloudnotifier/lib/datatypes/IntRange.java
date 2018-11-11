package com.blackforestbytes.simplecloudnotifier.lib.datatypes;

public class IntRange
{
	private int Start;
	public int Start() { return Start; }

	private int End;
	public int End() { return End; }

	public IntRange(int s, int e) { Start = s; End = e; }

	private IntRange() {  }
}
