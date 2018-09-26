package com.blackforestbytes.simplecloudnotifier;

import android.app.Application;
import android.arch.lifecycle.Lifecycle;
import android.arch.lifecycle.LifecycleObserver;
import android.arch.lifecycle.OnLifecycleEvent;
import android.arch.lifecycle.ProcessLifecycleOwner;
import android.content.Context;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.service.NotificationService;
import com.blackforestbytes.simplecloudnotifier.view.AccountFragment;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;
import com.blackforestbytes.simplecloudnotifier.view.TabAdapter;

import java.lang.ref.WeakReference;

public class SCNApp extends Application implements LifecycleObserver
{
    private static SCNApp instance;
    private static WeakReference<MainActivity> mainActivity;

    private static boolean isBackground = true;

    public SCNApp()
    {
        instance = this;
        ProcessLifecycleOwner.get().getLifecycle().addObserver(this);
    }

    public static Context getContext()
    {
        return instance;
    }

    public static boolean isBackground()
    {
        return isBackground;
    }

    public static void showToast(final String msg, final int duration)
    {
        final MainActivity a = mainActivity.get();
        if (a != null)
        {
            a.runOnUiThread(() -> Toast.makeText(a, msg, duration).show());
        }
    }

    public static boolean runOnUiThread(Runnable r)
    {
        final MainActivity a = mainActivity.get();
        if (a != null) {a.runOnUiThread(r); return true;}
        return false;
    }

    public static void refreshAccountTab()
    {
        runOnUiThread(() ->
        {
            MainActivity a = mainActivity.get();
            if (a == null) return;

            TabAdapter ta = a.adpTabs;
            if (ta == null) return;

            AccountFragment tf = ta.tab2;
            if (tf == null) return;

            tf.updateUI();
        });
    }

    public static void register(MainActivity a)
    {
        mainActivity = new WeakReference<>(a);
    }

    @OnLifecycleEvent(Lifecycle.Event.ON_STOP)
    public void onAppBackgrounded()
    {
        isBackground = true;
    }

    @OnLifecycleEvent(Lifecycle.Event.ON_START)
    public void onAppForegrounded()
    {
        isBackground = false;
    }
}