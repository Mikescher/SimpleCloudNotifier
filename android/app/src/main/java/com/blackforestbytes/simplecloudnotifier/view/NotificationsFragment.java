package com.blackforestbytes.simplecloudnotifier.view;

import android.graphics.Color;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.blackforestbytes.simplecloudnotifier.util.MessageAdapterTouchHelper;
import com.google.android.gms.ads.doubleclick.PublisherAdRequest;
import com.google.android.gms.ads.doubleclick.PublisherAdView;
import com.google.android.material.snackbar.Snackbar;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.recyclerview.widget.ItemTouchHelper;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

public class NotificationsFragment extends Fragment implements MessageAdapterTouchHelper.MessageAdapterTouchHelperListener
{
    private PublisherAdView adView;
    private MessageAdapter adpMessages;

    public MessageAdapterTouchHelper touchHelper;

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
        rvMessages.setAdapter(adpMessages = new MessageAdapter(v.findViewById(R.id.tvNoElements), lman, rvMessages));

        ItemTouchHelper.SimpleCallback itemTouchHelperCallback = touchHelper = new MessageAdapterTouchHelper(0, ItemTouchHelper.LEFT, this);
        new ItemTouchHelper(itemTouchHelperCallback).attachToRecyclerView(rvMessages);

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

    @Override
    public void onSwiped(RecyclerView.ViewHolder viewHolder, int direction, int position)
    {
        if (viewHolder instanceof MessageAdapter.MessagePresenter)
        {
            final int deletedIndex = viewHolder.getAdapterPosition();

            final CMessage deletedItem = adpMessages.removeItem(viewHolder.getAdapterPosition());
            String name = deletedItem.Title;

            Snackbar snackbar = Snackbar.make(SCNApp.getMainActivity().layoutRoot, name + " removed", Snackbar.LENGTH_LONG);
            snackbar.setAction("UNDO", view -> adpMessages.restoreItem(deletedItem, deletedIndex));
            snackbar.setActionTextColor(Color.YELLOW);
            snackbar.show();
        }
    }

    public void updateDeleteSwipeEnabled()
    {
        if (touchHelper != null) touchHelper.updateEnabled();
    }
}
