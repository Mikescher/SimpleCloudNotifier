package com.blackforestbytes.simplecloudnotifier.view.debug;

import android.content.Context;
import android.graphics.Color;
import androidx.annotation.NonNull;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.model.SingleQuery;

import org.joda.time.format.DateTimeFormat;
import org.joda.time.format.DateTimeFormatter;

public class QueryLogAdapter extends ArrayAdapter<SingleQuery>
{
    public static DateTimeFormatter UI_FULLTIME_FORMATTER = DateTimeFormat.forPattern("HH:mm:ss");

    public QueryLogAdapter(@NonNull Context context, @NonNull SingleQuery[] objects)
    {
        super(context, R.layout.adapter_querylog, objects);
    }

    @NonNull
    @Override
    public View getView(int position, View convertView, @NonNull ViewGroup parent)
    {
        View v = convertView;

        if (v == null) {
            LayoutInflater vi;
            vi = LayoutInflater.from(getContext());
            v = vi.inflate(R.layout.adapter_querylog, parent, false);
        }

        SingleQuery p = getItem(position);

        if (p != null)
        {
            TextView tt1 = v.findViewById(R.id.list_item_debuglogrow_time);
            if (tt1 != null) tt1.setText(p.Timestamp.toString(UI_FULLTIME_FORMATTER));
            if (tt1 != null) tt1.setTextColor(Color.BLACK);

            TextView tt2 = v.findViewById(R.id.list_item_debuglogrow_level);
            if (tt2 != null) tt2.setText(p.Level.toUIString());
            if (tt2 != null) tt2.setTextColor(Color.BLACK);

            TextView tt3 = v.findViewById(R.id.list_item_debuglogrow_info);
            if (tt3 != null) tt3.setText("");
            if (tt3 != null) tt3.setTextColor(Color.BLUE);

            TextView tt4 = v.findViewById(R.id.list_item_debuglogrow_id);
            if (tt4 != null) tt4.setText(p.Name);
            if (tt4 != null) tt4.setTextColor(p.Level.getColor());

            TextView tt5 = v.findViewById(R.id.list_item_debuglogrow_message);
            if (tt5 != null) tt5.setText(p.ExceptionString.length()> 40 ? p.ExceptionString.substring(0, 40-3)+"..." : p.ExceptionString);
            if (tt5 != null) tt5.setTextColor(p.Level.getColor());
        }

        return v;
    }
}
