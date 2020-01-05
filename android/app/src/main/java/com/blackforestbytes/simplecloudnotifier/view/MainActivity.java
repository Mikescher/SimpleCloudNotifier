package com.blackforestbytes.simplecloudnotifier.view;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.net.Uri;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.QueryLog;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.blackforestbytes.simplecloudnotifier.service.NotificationService;
import com.blackforestbytes.simplecloudnotifier.view.debug.QueryLogActivity;
import com.google.android.material.tabs.TabLayout;

import androidx.appcompat.app.AppCompatActivity;
import androidx.appcompat.widget.Toolbar;
import androidx.viewpager.widget.PagerAdapter;
import androidx.viewpager.widget.ViewPager;

import java.io.File;
import java.io.FileInputStream;
import java.io.ObjectInputStream;
import java.util.Map;
import java.util.Set;

public class MainActivity extends AppCompatActivity
{
    public TabAdapter adpTabs;
    public RelativeLayout layoutRoot;

    @Override
    protected void onCreate(Bundle savedInstanceState)
    {
        QueryLog.inst();

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        NotificationService.inst();
        CMessageList.inst();

        layoutRoot = findViewById(R.id.layoutRoot);

        Toolbar toolbar = findViewById(R.id.toolbar);
        toolbar.setOnClickListener(this::onToolbackClicked);
        setSupportActionBar(toolbar);

        ViewPager viewPager = findViewById(R.id.pager);
        PagerAdapter adapter = adpTabs = new TabAdapter(getSupportFragmentManager());
        viewPager.setAdapter(adapter);

        TabLayout tabLayout = findViewById(R.id.tab_layout);
        tabLayout.setTabGravity(TabLayout.GRAVITY_FILL);
        tabLayout.setupWithViewPager(viewPager);

        viewPager.addOnPageChangeListener(new ViewPager.OnPageChangeListener() {
            @Override
            public void onPageScrolled(int position, float positionOffset, int positionOffsetPixels) { /* */ }

            @Override
            public void onPageSelected(int position)
            {
                if (position != 2) adpTabs.tab3.onViewpagerHide();
            }

            @Override
            public void onPageScrollStateChanged(int state) {

            }
        });

        SCNApp.register(this);
        IABService.startup(this);
        SCNSettings.inst().work(this);
    }

    @Override
    protected void onStop()
    {
        super.onStop();

        SCNSettings.inst().save();
        CMessageList.inst().fullSave();
    }

    @Override
    protected void onDestroy()
    {
        super.onDestroy();

        CMessageList.inst().fullSave();
        IABService.inst().destroy();
    }

    private int clickCount = 0;
    private long lastClick = 0;
    private void onToolbackClicked(View v)
    {
        long now = System.currentTimeMillis();
        if (now - lastClick > 200) clickCount=0;
        clickCount++;
        lastClick = now;

        if (clickCount == 4) startActivity(new Intent(this, QueryLogActivity.class));
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        if(requestCode == 1991 && resultCode == RESULT_OK)
        {
            Uri uri = data.getData(); //The uri with the location of the file

            Context ctxt = this;

            try
            {
                ObjectInputStream stream = new ObjectInputStream(getContentResolver().openInputStream(uri));

                Map<String, ?> d1 = (Map<String, ?>)stream.readObject();
                Map<String, ?> d2 = (Map<String, ?>)stream.readObject();
                Map<String, ?> d3 = (Map<String, ?>)stream.readObject();
                Map<String, ?> d4 = (Map<String, ?>)stream.readObject();

                stream.close();

                runOnUiThread(() ->
                {

                    SharedPreferences.Editor e1 = ctxt.getSharedPreferences("Config", Context.MODE_PRIVATE).edit();
                    SharedPreferences.Editor e2 = ctxt.getSharedPreferences("IAB", Context.MODE_PRIVATE).edit();
                    SharedPreferences.Editor e3 = ctxt.getSharedPreferences("CMessageList", Context.MODE_PRIVATE).edit();
                    SharedPreferences.Editor e4 = ctxt.getSharedPreferences("QueryLog", Context.MODE_PRIVATE).edit();

                    e1.clear();
                    for (Map.Entry<String, ?> entry : d1.entrySet())
                    {
                        if (entry.getValue() instanceof String)  e1.putString(entry.getKey(),    (String)entry.getValue());
                        if (entry.getValue() instanceof Boolean) e1.putBoolean(entry.getKey(),   (Boolean)entry.getValue());
                        if (entry.getValue() instanceof Float)   e1.putFloat(entry.getKey(),     (Float)entry.getValue());
                        if (entry.getValue() instanceof Integer) e1.putInt(entry.getKey(),       (Integer)entry.getValue());
                        if (entry.getValue() instanceof Long)    e1.putLong(entry.getKey(),      (Long)entry.getValue());
                        if (entry.getValue() instanceof Set<?>)  e1.putStringSet(entry.getKey(), (Set<String>)entry.getValue());
                    }

                    e2.clear();
                    for (Map.Entry<String, ?> entry : d2.entrySet())
                    {
                        if (entry.getValue() instanceof String)  e2.putString(entry.getKey(),    (String)entry.getValue());
                        if (entry.getValue() instanceof Boolean) e2.putBoolean(entry.getKey(),   (Boolean)entry.getValue());
                        if (entry.getValue() instanceof Float)   e2.putFloat(entry.getKey(),     (Float)entry.getValue());
                        if (entry.getValue() instanceof Integer) e2.putInt(entry.getKey(),       (Integer)entry.getValue());
                        if (entry.getValue() instanceof Long)    e2.putLong(entry.getKey(),      (Long)entry.getValue());
                        if (entry.getValue() instanceof Set<?>)  e2.putStringSet(entry.getKey(), (Set<String>)entry.getValue());
                    }

                    e2.clear();
                    for (Map.Entry<String, ?> entry : d3.entrySet())
                    {
                        if (entry.getValue() instanceof String)  e3.putString(entry.getKey(),    (String)entry.getValue());
                        if (entry.getValue() instanceof Boolean) e3.putBoolean(entry.getKey(),   (Boolean)entry.getValue());
                        if (entry.getValue() instanceof Float)   e3.putFloat(entry.getKey(),     (Float)entry.getValue());
                        if (entry.getValue() instanceof Integer) e3.putInt(entry.getKey(),       (Integer)entry.getValue());
                        if (entry.getValue() instanceof Long)    e3.putLong(entry.getKey(),      (Long)entry.getValue());
                        if (entry.getValue() instanceof Set<?>)  e3.putStringSet(entry.getKey(), (Set<String>)entry.getValue());
                    }

                    e4.clear();
                    for (Map.Entry<String, ?> entry : d4.entrySet())
                    {
                        if (entry.getValue() instanceof String)  e4.putString(entry.getKey(),    (String)entry.getValue());
                        if (entry.getValue() instanceof Boolean) e4.putBoolean(entry.getKey(),   (Boolean)entry.getValue());
                        if (entry.getValue() instanceof Float)   e4.putFloat(entry.getKey(),     (Float)entry.getValue());
                        if (entry.getValue() instanceof Integer) e4.putInt(entry.getKey(),       (Integer)entry.getValue());
                        if (entry.getValue() instanceof Long)    e4.putLong(entry.getKey(),      (Long)entry.getValue());
                        if (entry.getValue() instanceof Set<?>)  e4.putStringSet(entry.getKey(), (Set<String>)entry.getValue());
                    }

                    e1.apply();
                    e2.apply();
                    e3.apply();
                    e4.apply();


                    SCNSettings.inst().reloadPrefs();
                    IABService.inst().reloadPrefs();
                    CMessageList.inst().reloadPrefs();
                    QueryLog.inst().reloadPrefs();


                    Toolbar toolbar = findViewById(R.id.toolbar);
                    setSupportActionBar(toolbar);

                    ViewPager viewPager = findViewById(R.id.pager);
                    PagerAdapter adapter = adpTabs = new TabAdapter(getSupportFragmentManager());
                    viewPager.setAdapter(adapter);

                    TabLayout tabLayout = findViewById(R.id.tab_layout);
                    tabLayout.setupWithViewPager(viewPager);


                    SCNSettings.inst().work(this);

                    SCNApp.showToast("Backup imported", Toast.LENGTH_LONG);

                    finish();
                });
            }
            catch (Exception e)
            {
                Log.e("Import:Err", e.toString());
                SCNApp.showToast("Import failed", Toast.LENGTH_LONG);
            }
        }
    }
}
