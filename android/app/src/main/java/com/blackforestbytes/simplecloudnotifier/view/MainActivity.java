package com.blackforestbytes.simplecloudnotifier.view;

import android.content.Intent;
import android.icu.text.SymbolTable;
import android.os.Bundle;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.TextView;
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

public class MainActivity extends AppCompatActivity
{
    public TabAdapter adpTabs;
    public RelativeLayout layoutRoot;

    @Override
    protected void onCreate(Bundle savedInstanceState)
    {
        QueryLog.instance();

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
}
