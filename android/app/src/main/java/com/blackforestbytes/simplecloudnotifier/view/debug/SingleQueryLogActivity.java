package com.blackforestbytes.simplecloudnotifier.view.debug;

import androidx.appcompat.app.AppCompatActivity;

import android.annotation.SuppressLint;
import android.os.Bundle;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.lib.string.CompactJsonFormatter;
import com.blackforestbytes.simplecloudnotifier.model.SingleQuery;

import org.joda.time.format.DateTimeFormat;
import org.joda.time.format.DateTimeFormatter;

import java.util.Objects;

public class SingleQueryLogActivity extends AppCompatActivity
{
    @Override
    @SuppressLint("SetTextI18n")
    protected void onCreate(Bundle savedInstanceState)
    {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_singlequerylog);

        SingleQuery q = SingleQuery.load(getIntent().getBundleExtra("query"), "data");

        this.<TextView>findViewById(R.id.tvQL_Timestamp).setText(q.Timestamp.toString(DateTimeFormat.forPattern("yyyy-MM-dd HH:mm:ss")));
        this.<TextView>findViewById(R.id.tvQL_Level).setText(q.Level.toUIString());
        this.<TextView>findViewById(R.id.tvQL_Level).setTextColor(q.Level.getColor());
        this.<TextView>findViewById(R.id.tvQL_Name).setText(q.Name);
        this.<TextView>findViewById(R.id.tvQL_URL).setText(q.URL.replace("?", "\r\n?").replace("&", "\r\n&"));
        this.<TextView>findViewById(R.id.tvQL_Response).setText(CompactJsonFormatter.formatJSON(q.Response, 999));
        this.<TextView>findViewById(R.id.tvQL_ResponseCode).setText(Integer.toString(q.ResponseCode));
        this.<TextView>findViewById(R.id.tvQL_ExceptionString).setText(q.ExceptionString);
    }
}
