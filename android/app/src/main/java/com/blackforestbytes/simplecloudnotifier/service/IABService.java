package com.blackforestbytes.simplecloudnotifier.service;

import android.app.Activity;
import android.content.Context;
import android.util.Log;
import android.widget.Toast;

import com.android.billingclient.api.BillingClient;
import com.android.billingclient.api.BillingClientStateListener;
import com.android.billingclient.api.BillingFlowParams;
import com.android.billingclient.api.Purchase;
import com.android.billingclient.api.PurchasesUpdatedListener;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.lambda.Func0to0;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

import java.util.ArrayList;
import java.util.List;

import androidx.annotation.Nullable;

import static androidx.constraintlayout.widget.Constraints.TAG;

public class IABService implements PurchasesUpdatedListener
{
    public static final String IAB_PRO_MODE = "scn.pro.tier1";

    private final static Object _lock = new Object();
    private static IABService _inst = null;
    public static IABService inst()
    {
        synchronized (_lock)
        {
            if (_inst != null) return _inst;
            throw new Error("IABService == null");
        }
    }
    public static void startup(MainActivity a)
    {
        synchronized (_lock)
        {
            _inst = new IABService(a);
        }
    }

    private BillingClient client;
    private boolean isServiceConnected;
    private final List<Purchase> purchases = new ArrayList<>();

    public IABService(Context c)
    {
        client = BillingClient
                .newBuilder(c)
                .setListener(this)
                .build();

        startServiceConnection(this::queryPurchases, false);
    }

    public void queryPurchases()
    {
        Func0to0 queryToExecute = () ->
        {
            long time = System.currentTimeMillis();
            Purchase.PurchasesResult purchasesResult = client.queryPurchases(BillingClient.SkuType.INAPP);
            Log.i(TAG, "Querying purchases elapsed time: " + (System.currentTimeMillis() - time) + "ms");

            if (purchasesResult.getResponseCode() == BillingClient.BillingResponse.OK)
            {
                for (Purchase p : purchasesResult.getPurchasesList())
                {
                    handlePurchase(p);
                }

                boolean newProMode = getPurchaseCached(IAB_PRO_MODE) != null;
                if (newProMode != SCNSettings.inst().promode_local)
                {
                    refreshProModeListener();
                }
            }
            else
            {
                Log.w(TAG, "queryPurchases() got an error response code: " + purchasesResult.getResponseCode());
            }
        };

        executeServiceRequest(queryToExecute, false);
    }

    public void purchase(Activity a, String id)
    {
        executeServiceRequest(() ->
        {
            BillingFlowParams flowParams = BillingFlowParams
                    .newBuilder()
                    .setSku(id)
                    .setType(BillingClient.SkuType.INAPP) // SkuType.SUB for subscription
                    .build();
            client.launchBillingFlow(a, flowParams);
        }, true);
    }

    private void executeServiceRequest(Func0to0 runnable, final boolean userRequest) {
        if (isServiceConnected)
        {
            runnable.invoke();
        }
        else
        {
            // If billing service was disconnected, we try to reconnect 1 time.
            // (feel free to introduce your retry policy here).
            startServiceConnection(runnable, userRequest);
        }
    }
    public void destroy()
    {
        if (client != null && client.isReady()) {
            client.endConnection();
            client = null;
            isServiceConnected = false;
        }
    }

    @Override
    public void onPurchasesUpdated(int responseCode, @Nullable List<Purchase> purchases)
    {
        if (responseCode == BillingClient.BillingResponse.OK && purchases != null)
        {
            for (Purchase purchase : purchases)
            {
                handlePurchase(purchase);
            }
        }
        else if (responseCode == BillingClient.BillingResponse.ITEM_ALREADY_OWNED && purchases != null)
        {
            for (Purchase purchase : purchases)
            {
                handlePurchase(purchase);
            }
        }
    }

    private void handlePurchase(Purchase purchase)
    {
        Log.d(TAG, "Got a verified purchase: " + purchase);

        purchases.add(purchase);

        refreshProModeListener();
    }

    private void refreshProModeListener() {
        MainActivity ma = SCNApp.getMainActivity();
        if (ma != null) ma.adpTabs.tab3.updateProState();
        if (ma != null) ma.adpTabs.tab1.updateProState();
        SCNSettings.inst().updateProState(null);
    }

    public void startServiceConnection(final Func0to0 executeOnSuccess, final boolean userRequest)
    {
        client.startConnection(new BillingClientStateListener()
        {
            @Override
            public void onBillingSetupFinished(@BillingClient.BillingResponse int billingResponseCode)
            {
                if (billingResponseCode == BillingClient.BillingResponse.OK)
                {
                    isServiceConnected = true;
                    if (executeOnSuccess != null) executeOnSuccess.invoke();
                }
                else
                {
                    if (userRequest) SCNApp.showToast("Could not connect to google services", Toast.LENGTH_SHORT);
                }
            }

            @Override
            public void onBillingServiceDisconnected() {
                isServiceConnected = false;
            }
        });
    }

    public Purchase getPurchaseCached(String id)
    {
        for (Purchase p : purchases)
        {
            if (Str.equals(p.getSku(), id)) return p;
        }

        return null;
    }
}
