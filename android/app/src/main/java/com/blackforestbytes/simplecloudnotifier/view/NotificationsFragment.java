package com.blackforestbytes.simplecloudnotifier.view;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

import com.blackforestbytes.simplecloudnotifier.R;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

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
        rvMessages.setAdapter(new MessageAdapter(v.findViewById(R.id.tvNoElements)));

        return v;
    }
}
