package com.blackforestbytes.simplecloudnotifier.view;

import android.net.Uri;
import android.os.Bundle;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.NotificationService;
import com.google.android.material.tabs.TabLayout;

import org.jetbrains.annotations.NotNull;

import androidx.appcompat.app.AppCompatActivity;
import androidx.appcompat.widget.Toolbar;
import androidx.viewpager.widget.PagerAdapter;
import androidx.viewpager.widget.ViewPager;
import xyz.aprildown.ultimatemusicpicker.MusicPickerListener;

public class MainActivity extends AppCompatActivity implements MusicPickerListener
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

        SCNApp.register(this);

        SCNSettings.inst().work(this);
    }

    @Override
    protected void onStop()
    {
        super.onStop();

        CMessageList.inst().fullSave();
    }

    @Override
    public void onMusicPick(@NotNull Uri uri, @NotNull String s) {

    }

    @Override
    public void onPickCanceled() {

    }
}
