package com.blackforestbytes.simplecloudnotifier.view;

import android.os.Bundle;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.blackforestbytes.simplecloudnotifier.service.NotificationService;
import com.google.android.material.tabs.TabLayout;

import androidx.appcompat.app.AppCompatActivity;
import androidx.appcompat.widget.Toolbar;
import androidx.viewpager.widget.PagerAdapter;
import androidx.viewpager.widget.ViewPager;

public class MainActivity extends AppCompatActivity
{
    public TabAdapter adpTabs;

    @Override
    protected void onCreate(Bundle savedInstanceState)
    {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        NotificationService.inst();
        CMessageList.inst();

        Toolbar toolbar = findViewById(R.id.toolbar);
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
}
