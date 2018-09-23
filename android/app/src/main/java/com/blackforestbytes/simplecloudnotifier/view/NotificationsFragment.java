package com.blackforestbytes.simplecloudnotifier.view;

import android.content.Context;
import android.net.Uri;
import android.os.Bundle;
import android.support.annotation.NonNull;
import android.support.v4.app.Fragment;
import android.support.v7.widget.LinearLayoutManager;
import android.support.v7.widget.RecyclerView;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

import com.blackforestbytes.simplecloudnotifier.R;

public class NotificationsFragment extends Fragment
{
    public NotificationsFragment()
    {
        // Required empty public constructor
    }

    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState)
    {
        View v = inflater.inflate(R.layout.fragment_notifications, container, false);

        RecyclerView rvMessages = v.findViewById(R.id.rvMessages);
        rvMessages.setLayoutManager(new LinearLayoutManager(this.getContext(), RecyclerView.VERTICAL, true));
        rvMessages.setAdapter(new MessageAdapter());

        return v;
    }
}
