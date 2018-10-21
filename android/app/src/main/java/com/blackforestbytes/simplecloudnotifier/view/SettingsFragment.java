package com.blackforestbytes.simplecloudnotifier.view;

import android.media.AudioManager;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.ImageView;
import android.widget.Spinner;
import android.widget.Switch;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import xyz.aprildown.ultimatemusicpicker.UltimateMusicPicker;

public class SettingsFragment extends Fragment
{
    private Switch  prefAppEnabled;
    private Spinner prefLocalCacheSize;

    private Switch    prefMsgLowEnableSound;
    private TextView  tvMsgLowRingtone_value;
    private View      prevMsgLowRingtone;
    private Switch    prefMsgLowRepeatSound;
    private Switch    prefMsgLowEnableLED;
    private TextView  tvMsgLowLedColor;
    private ImageView prefMsgLowLedColor;
    private Switch    prefMsgLowEnableVibrations;

    private Switch    prefMsgNormEnableSound;
    private TextView  tvMsgNormRingtone_value;
    private View      prevMsgNormRingtone;
    private Switch    prefMsgNormRepeatSound;
    private Switch    prefMsgNormEnableLED;
    private TextView  tvMsgNormLedColor;
    private ImageView prefMsgNormLedColor;
    private Switch    prefMsgNormEnableVibrations;

    private Switch    prefMsgHighEnableSound;
    private TextView  tvMsgHighRingtone_value;
    private View      prevMsgHighRingtone;
    private Switch    prefMsgHighRepeatSound;
    private Switch    prefMsgHighEnableLED;
    private TextView  tvMsgHighLedColor;
    private ImageView prefMsgHighLedColor;
    private Switch    prefMsgHighEnableVibrations;

    public SettingsFragment()
    {
        // Required empty public constructor
    }

    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState)
    {
        View v = inflater.inflate(R.layout.fragment_settings, container, false);

        {
            prefAppEnabled             = v.findViewById(R.id.prefAppEnabled);
            prefLocalCacheSize         = v.findViewById(R.id.prefLocalCacheSize);

            prefMsgLowEnableSound      = v.findViewById(R.id.prefMsgLowEnableSound);
            tvMsgLowRingtone_value     = v.findViewById(R.id.tvMsgLowRingtone_value);
            prevMsgLowRingtone         = v.findViewById(R.id.prevMsgLowRingtone);
            prefMsgLowRepeatSound      = v.findViewById(R.id.prefMsgLowRepeatSound);
            prefMsgLowEnableLED        = v.findViewById(R.id.prefMsgLowEnableLED);
            tvMsgLowLedColor           = v.findViewById(R.id.tvMsgLowLedColor);
            prefMsgLowLedColor         = v.findViewById(R.id.prefMsgLowLedColor);
            prefMsgLowEnableVibrations = v.findViewById(R.id.prefMsgLowEnableVibrations);

            prefMsgNormEnableSound      = v.findViewById(R.id.prefMsgNormEnableSound);
            tvMsgNormRingtone_value     = v.findViewById(R.id.tvMsgNormRingtone_value);
            prevMsgNormRingtone         = v.findViewById(R.id.prevMsgNormRingtone);
            prefMsgNormRepeatSound      = v.findViewById(R.id.prefMsgNormRepeatSound);
            prefMsgNormEnableLED        = v.findViewById(R.id.prefMsgNormEnableLED);
            tvMsgNormLedColor           = v.findViewById(R.id.tvMsgNormLedColor);
            prefMsgNormLedColor         = v.findViewById(R.id.prefMsgNormLedColor);
            prefMsgNormEnableVibrations = v.findViewById(R.id.prefMsgNormEnableVibrations);

            prefMsgHighEnableSound      = v.findViewById(R.id.prefMsgHighEnableSound);
            tvMsgHighRingtone_value     = v.findViewById(R.id.tvMsgHighRingtone_value);
            prevMsgHighRingtone         = v.findViewById(R.id.prevMsgHighRingtone);
            prefMsgHighRepeatSound      = v.findViewById(R.id.prefMsgHighRepeatSound);
            prefMsgHighEnableLED        = v.findViewById(R.id.prefMsgHighEnableLED);
            tvMsgHighLedColor           = v.findViewById(R.id.tvMsgHighLedColor);
            prefMsgHighLedColor         = v.findViewById(R.id.prefMsgHighLedColor);
            prefMsgHighEnableVibrations = v.findViewById(R.id.prefMsgHighEnableVibrations);
        }

        {
            SCNSettings s = SCNSettings.inst();

            prefAppEnabled.setChecked(s.Enabled);
            prefAppEnabled.setOnCheckedChangeListener((a,b) -> onUpdate());

            ArrayAdapter<Integer> plcsa = new ArrayAdapter<>(v.getContext(), android.R.layout.simple_spinner_item, SCNSettings.CHOOSABLE_CACHE_SIZES);
            plcsa.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
            prefLocalCacheSize.setAdapter(plcsa);
            prefLocalCacheSize.setSelection(getCacheSizeIndex(s.LocalCacheSize));
            prefLocalCacheSize.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener()
            {
                @Override public void onItemSelected(AdapterView<?> parent, View view, int position, long id) { onUpdate(); }
                @Override public void onNothingSelected(AdapterView<?> parent) { onUpdate();  }
            });

            //TODO ...

            prevMsgLowRingtone.setOnClickListener((a) -> chooseRingtoneLow());
        }

        return v;
    }

    private void onUpdate()
    {
        SCNSettings s = SCNSettings.inst();

        s.Enabled        = prefAppEnabled.isChecked();
        s.LocalCacheSize = prefLocalCacheSize.getSelectedItemPosition()>=0 ? SCNSettings.CHOOSABLE_CACHE_SIZES[prefLocalCacheSize.getSelectedItemPosition()] : 100;

        s.save();
    }

    private int getCacheSizeIndex(int value)
    {
        for (int i = 0; i < SCNSettings.CHOOSABLE_CACHE_SIZES.length; i++)
        {
            if (SCNSettings.CHOOSABLE_CACHE_SIZES[i] == value) return i;
        }
        return 2;
    }

    private void chooseRingtoneLow()
    {
        new UltimateMusicPicker()
            .windowTitle("Choose notification sound")
            .removeSilent()
            .streamType(AudioManager.STREAM_ALARM)
            .ringtone()
            .notification()
            .alarm()
            .music()
            .goWithDialog(SCNApp.getMainActivity().getSupportFragmentManager());
    }
}