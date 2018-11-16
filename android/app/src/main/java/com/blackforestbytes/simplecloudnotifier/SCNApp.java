package com.blackforestbytes.simplecloudnotifier;

import android.app.Application;
import android.content.Context;
import android.widget.Toast;

import com.android.billingclient.api.BillingClient;
import com.blackforestbytes.simplecloudnotifier.view.AccountFragment;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;
import com.blackforestbytes.simplecloudnotifier.view.TabAdapter;

import java.lang.ref.WeakReference;

import androidx.lifecycle.Lifecycle;
import androidx.lifecycle.LifecycleObserver;
import androidx.lifecycle.OnLifecycleEvent;
import androidx.lifecycle.ProcessLifecycleOwner;

public class SCNApp extends Application implements LifecycleObserver
{
    private static SCNApp instance;
    private static WeakReference<MainActivity> mainActivity = new WeakReference<>(null);

    public static final boolean LOCAL_DEBUG = BuildConfig.DEBUG;
    public static final boolean DEBUG       = BuildConfig.DEBUG || !BuildConfig.VERSION_NAME.endsWith(".0");
    public static final boolean RELEASE     = !DEBUG;

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

    public static MainActivity getMainActivity()
    {
        return mainActivity.get();
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

/*
 ==TODO==

[ ]  - test notification channels

[ ]  - Delete single message (swipe right)
[ ]  - Query non-ack-ed messages in app
[ ]  - periodically get non-ack (option - even when not in-app)

[ ]  - publish (+ HN post ?)

[ ]  - Use for mscom server errrors
[ ]  - Use for bfb server errors
[ ]  - Use for transmission state
[ ]     - Message on connnection lost (seperate process - resend until succ)
[ ]     - Message on connnection regained
[ ]     - Message on seed-count changed

*/