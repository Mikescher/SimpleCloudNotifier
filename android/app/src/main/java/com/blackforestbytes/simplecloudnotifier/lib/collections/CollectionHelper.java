package com.blackforestbytes.simplecloudnotifier.lib.collections;

import com.blackforestbytes.simplecloudnotifier.lib.lambda.Func1to1;

import java.util.*;

public final class CollectionHelper
{
	public static <T, C> List<T> unique(List<T> input, Func1to1<T, C> mapping)
	{
		List<T> output = new ArrayList<>(input.size());

		HashSet<C> seen = new HashSet<>();

		for (T v : input) if (seen.add(mapping.invoke(v))) output.add(v);

		return output;
	}

	public static <T> List<T> sort(List<T> input, Comparator<T> comparator)
	{
		List<T> output = new ArrayList<>(input);
		Collections.sort(output, comparator);
		return output;
	}

	public static <T, U extends Comparable<U>> List<T> sort(List<T> input, Func1to1<T, U> mapper)
	{
		return sort(input, mapper, 1);
	}

	public static <T, U extends Comparable<U>> List<T> sort(List<T> input, Func1to1<T, U> mapper, int sortMod)
	{
		List<T> output = new ArrayList<>(input);
		Collections.sort(output, (o1, o2) -> sortMod * mapper.invoke(o1).compareTo(mapper.invoke(o2)));
		return output;
	}
}
