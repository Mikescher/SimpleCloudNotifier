package com.blackforestbytes.simplecloudnotifier.view.debug;

import android.content.Intent;
import android.os.Bundle;

import com.blackforestbytes.simplecloudnotifier.model.QueryLog;
import com.blackforestbytes.simplecloudnotifier.model.SingleQuery;

import androidx.appcompat.app.AppCompatActivity;

import android.widget.ListView;

import com.blackforestbytes.simplecloudnotifier.R;

public class QueryLogActivity extends AppCompatActivity
{

    @Override
    protected void onCreate(Bundle savedInstanceState)
    {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_querylog);

        ListView lvMain = findViewById(R.id.lvQueryList);
        SingleQuery[] arr = QueryLog.inst().get().toArray(new SingleQuery[0]);
        QueryLogAdapter a = new QueryLogAdapter(this, arr);
        lvMain.setAdapter(a);

        lvMain.setOnItemClickListener((parent, view, position, id) ->
        {
            if (position >= 0 && position < arr.length)
            {
                Intent i = new Intent(QueryLogActivity.this, SingleQueryLogActivity.class);
                Bundle b = new Bundle();
                arr[position].save(b, "data");
                i.putExtra("query", b);
                startActivity(i);
            }
        });


    }

}
