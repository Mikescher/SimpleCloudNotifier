package com.blackforestbytes.simplecloudnotifier.service;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;
import android.widget.Toast;

import com.android.billingclient.api.BillingClient;
import com.android.billingclient.api.BillingClientStateListener;
import com.android.billingclient.api.BillingFlowParams;
import com.android.billingclient.api.BillingResult;
import com.android.billingclient.api.Purchase;
import com.android.billingclient.api.PurchasesUpdatedListener;
import com.android.billingclient.api.SkuDetails;
import com.android.billingclient.api.SkuDetailsParams;
import com.android.billingclient.api.SkuDetailsResponseListener;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.datatypes.Tuple2;
import com.blackforestbytes.simplecloudnotifier.lib.datatypes.Tuple3;
import com.blackforestbytes.simplecloudnotifier.lib.lambda.Func0to0;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.view.MainActivity;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Dictionary;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

import androidx.annotation.NonNull;
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

    public enum SimplePurchaseState { YES, NO, UNINITIALIZED }

    private BillingClient client;
    private boolean isServiceConnected;
    private final List<Purchase> purchases = new ArrayList<>();
    private boolean _isInitialized = false;

    private final Map<String, Boolean> _localCache= new HashMap<>();

    public IABService(Context c)
    {
        _isInitialized = false;

        loadCache();

        client = BillingClient
                .newBuilder(c)
                .setListener(this)
                .build();

        startServiceConnection(this::queryPurchases, false);
        startServiceConnection(this::querySkuDetails, false);
    }

    public void reloadPrefs()
    {
        loadCache();
    }

    private void loadCache()
    {
        _localCache.clear();
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("IAB", Context.MODE_PRIVATE);
        int count = sharedPref.getInt("c", 0);
        for (int i=0; i < count; i++)
        {
            String  k = sharedPref.getString("["+i+"]->key", null);
            boolean v = sharedPref.getBoolean("["+i+"]->value", false);
            if (k==null)continue;
            _localCache.put(k, v);
        }
    }

    private void saveCache()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("IAB", Context.MODE_PRIVATE);
        SharedPreferences.Editor editor= sharedPref.edit();

        editor.putInt("c", _localCache.size());
        int i = 0;
        for (Map.Entry<String, Boolean> e : _localCache.entrySet())
        {
            editor.putString("["+i+"]->key", e.getKey());
            editor.putBoolean("["+i+"]->value", e.getValue());
            i++;
        }
        editor.apply();
    }

    @SuppressWarnings("ConstantConditions")
    private synchronized void updateCache(String k, boolean v)
    {
        if (_localCache.containsKey(k) && _localCache.get(k)==v) return;

        _localCache.put(k, v);
        saveCache();
    }

    public void queryPurchases()
    {
        Func0to0 queryToExecute = () ->
        {
            long time = System.currentTimeMillis();
            Purchase.PurchasesResult purchasesResult = client.queryPurchases(BillingClient.SkuType.INAPP);
            Log.i(TAG, "Querying purchases elapsed time: " + (System.currentTimeMillis() - time) + "ms");

            if (purchasesResult.getResponseCode() == BillingClient.BillingResponseCode.OK)
            {
                for (Purchase p : Objects.requireNonNull(purchasesResult.getPurchasesList()))
                {
                    handlePurchase(p, false);
                }

                _isInitialized = true;

                boolean newProMode = getPurchaseCachedSimple(IAB_PRO_MODE);
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

    public void querySkuDetails() {
    }

    public void purchase(Activity a, String id)
    {
        Func0to0 queryRequest = () -> {
            // Query the purchase async
            SkuDetailsParams.Builder params = SkuDetailsParams.newBuilder();
            params.setSkusList(Collections.singletonList(id)).setType(BillingClient.SkuType.INAPP);
            client.querySkuDetailsAsync(params.build(), (billingResult, skuDetailsList) ->
            {
                if (billingResult.getResponseCode() != BillingClient.BillingResponseCode.OK || skuDetailsList == null || skuDetailsList.size() != 1)
                {
                    SCNApp.showToast("Could not find product", Toast.LENGTH_SHORT);
                    return;
                }

                executeServiceRequest(() ->
                {
                    BillingFlowParams flowParams = BillingFlowParams
                            .newBuilder()
                            .setSkuDetails(skuDetailsList.get(0))
                            .build();
                    client.launchBillingFlow(a, flowParams);
                }, true);
            });
        };
        executeServiceRequest(queryRequest, false);

    }

    private void executeServiceRequest(Func0to0 runnable, final boolean userRequest)
    {
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
    public void onPurchasesUpdated(@NonNull BillingResult billingResult, @Nullable List<Purchase> purchases)
    {
        if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.OK && purchases != null)
        {
            for (Purchase purchase : purchases)
            {
                handlePurchase(purchase, true);
            }
        }
        else if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.ITEM_ALREADY_OWNED && purchases != null)
        {
            for (Purchase purchase : purchases)
            {
                handlePurchase(purchase, true);
            }
        }
    }

    private void handlePurchase(Purchase purchase, boolean triggerUpdate)
    {
        Log.d(TAG, "Got a verified purchase: " + purchase);

        purchases.add(purchase);

        if (triggerUpdate) refreshProModeListener();

        updateCache(purchase.getSku(), true);
    }

    private void refreshProModeListener()
    {
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
            public void onBillingSetupFinished(@NonNull BillingResult billingResult)
            {
                if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.OK)
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

    public boolean getPurchaseCachedSimple(String id)
    {
        return getPurchaseCachedExtended(id).Item1;
    }

    @SuppressWarnings("ConstantConditions")
    public Tuple3<Boolean, Boolean, String> getPurchaseCachedExtended(String id)
    {
        // <state, initialized, token>

        if (!_isInitialized)
        {
            if (_localCache.containsKey(id) && _localCache.get(id)) return new Tuple3<>(true, false, Str.Empty);
        }

        for (Purchase p : purchases)
        {
            if (Str.equals(p.getSku(), id))
            {
                updateCache(id, true);
                return new Tuple3<>(true, true, p.getPurchaseToken());
            }
        }

        updateCache(id, false);
        return new Tuple3<>(false, true, Str.Empty);
    }
}
