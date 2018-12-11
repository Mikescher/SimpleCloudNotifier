package com.blackforestbytes.simplecloudnotifier.util;

import android.content.Context;
import android.content.res.Resources;
import android.content.res.TypedArray;
import android.util.AttributeSet;
import android.util.TypedValue;
import android.widget.ScrollView;

import com.blackforestbytes.simplecloudnotifier.R;

public class MaxHeightScrollView extends ScrollView
{
	public int maxHeight = Integer.MAX_VALUE;//dp

	public MaxHeightScrollView(Context context)
	{
		super(context);
	}

	public MaxHeightScrollView(Context context, AttributeSet attrs)
	{
		super(context, attrs);

		TypedArray a = context.getTheme().obtainStyledAttributes(attrs, R.styleable.MaxHeightScrollView, 0, 0);
		try {
			maxHeight = a.getInteger(R.styleable.MaxHeightScrollView_maxHeightOverride, Integer.MAX_VALUE);
		} finally {
			a.recycle();
		}
	}

	public MaxHeightScrollView(Context context, AttributeSet attrs, int defStyleAttr)
	{
		super(context, attrs, defStyleAttr);

		TypedArray a = context.getTheme().obtainStyledAttributes(attrs, R.styleable.MaxHeightScrollView, 0, 0);
		try {
			maxHeight = a.getInteger(R.styleable.MaxHeightScrollView_maxHeightOverride, Integer.MAX_VALUE);
		} finally {
			a.recycle();
		}
	}

	@Override
	protected void onMeasure(int widthMeasureSpec, int heightMeasureSpec)
	{
		heightMeasureSpec = MeasureSpec.makeMeasureSpec(dpToPx(getResources(), maxHeight), MeasureSpec.AT_MOST);
		super.onMeasure(widthMeasureSpec, heightMeasureSpec);
	}

	private int dpToPx(Resources res, int dp)
	{
		return (int) TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, dp, res.getDisplayMetrics());
	}
}
