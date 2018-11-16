package com.blackforestbytes.simplecloudnotifier.view;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.google.android.gms.ads.doubleclick.PublisherAdRequest;
import com.google.android.gms.ads.doubleclick.PublisherAdView;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

public class NotificationsFragment extends Fragment
{
    private PublisherAdView adView;

    public NotificationsFragment()
    {
        // Required empty public constructor
    }

    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState)
    {
        View v = inflater.inflate(R.layout.fragment_notifications, container, false);

        RecyclerView rvMessages = v.findViewById(R.id.rvMessages);
        LinearLayoutManager lman = new LinearLayoutManager(this.getContext(), RecyclerView.VERTICAL, false);
        rvMessages.setLayoutManager(lman);
        rvMessages.setAdapter(new MessageAdapter(v.findViewById(R.id.tvNoElements), lman, rvMessages));
        //lman.scrollToPosition(0);

        adView = v.findViewById(R.id.adBanner);
        PublisherAdRequest adRequest = new PublisherAdRequest.Builder().build();
        adView.loadAd(adRequest);

        adView.setVisibility(SCNSettings.inst().promode_local ? View.GONE : View.VISIBLE);

        return v;
    }

    public void updateProState()
    {
        if (adView != null) adView.setVisibility(IABService.inst().getPurchaseCached(IABService.IAB_PRO_MODE) != null ? View.GONE : View.VISIBLE);
    }
}
