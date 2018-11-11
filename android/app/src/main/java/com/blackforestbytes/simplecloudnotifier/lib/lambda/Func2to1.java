package com.blackforestbytes.simplecloudnotifier.lib.lambda;

@FunctionalInterface
public interface Func2to1<TInput1, TInput2, TResult> {
	TResult invoke(TInput1 value1, TInput2 value2);
}
