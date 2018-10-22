package com.blackforestbytes.simplecloudnotifier.view;

import android.content.Context;
import android.media.AudioManager;
import android.net.Uri;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.Spinner;
import android.widget.Switch;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;

import org.jetbrains.annotations.NotNull;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import top.defaults.colorpicker.ColorPickerPopup;
import xyz.aprildown.ultimatemusicpicker.MusicPickerListener;
import xyz.aprildown.ultimatemusicpicker.UltimateMusicPicker;

public class SettingsFragment extends Fragment implements MusicPickerListener
{
    private Switch    prefAppEnabled;
    private Spinner   prefLocalCacheSize;
    private Button    prefUpgradeAccount;

    private Switch    prefMsgLowEnableSound;
    private TextView  tvMsgLowRingtone_value;
    private View      prefMsgLowRingtone;
    private Switch    prefMsgLowRepeatSound;
    private Switch    prefMsgLowEnableLED;
    private View      prefMsgLowLedColor_container;
    private ImageView prefMsgLowLedColor_value;
    private Switch    prefMsgLowEnableVibrations;

    private Switch    prefMsgNormEnableSound;
    private TextView  tvMsgNormRingtone_value;
    private View      prefMsgNormRingtone;
    private Switch    prefMsgNormRepeatSound;
    private Switch    prefMsgNormEnableLED;
    private View      prefMsgNormLedColor_container;
    private ImageView prefMsgNormLedColor_value;
    private Switch    prefMsgNormEnableVibrations;

    private Switch    prefMsgHighEnableSound;
    private TextView  tvMsgHighRingtone_value;
    private View      prefMsgHighRingtone;
    private Switch    prefMsgHighRepeatSound;
    private Switch    prefMsgHighEnableLED;
    private View      prefMsgHighLedColor_container;
    private ImageView prefMsgHighLedColor_value;
    private Switch    prefMsgHighEnableVibrations;

    private int musicPickerSwitch = -1;

    public SettingsFragment()
    {
        // Required empty public constructor
    }

    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState)
    {
        View v = inflater.inflate(R.layout.fragment_settings, container, false);

        initFields(v);
        updateUI();
        initListener();

        return v;
    }

    private void initFields(View v)
    {
        prefAppEnabled                = v.findViewById(R.id.prefAppEnabled);
        prefLocalCacheSize            = v.findViewById(R.id.prefLocalCacheSize);
        prefUpgradeAccount            = v.findViewById(R.id.prefUpgradeAccount);

        prefMsgLowEnableSound         = v.findViewById(R.id.prefMsgLowEnableSound);
        tvMsgLowRingtone_value        = v.findViewById(R.id.tvMsgLowRingtone_value);
        prefMsgLowRingtone            = v.findViewById(R.id.prefMsgLowRingtone);
        prefMsgLowRepeatSound         = v.findViewById(R.id.prefMsgLowRepeatSound);
        prefMsgLowEnableLED           = v.findViewById(R.id.prefMsgLowEnableLED);
        prefMsgLowLedColor_value      = v.findViewById(R.id.prefMsgLowLedColor_value);
        prefMsgLowLedColor_container  = v.findViewById(R.id.prefMsgLowLedColor_container);
        prefMsgLowEnableVibrations    = v.findViewById(R.id.prefMsgLowEnableVibrations);

        prefMsgNormEnableSound        = v.findViewById(R.id.prefMsgNormEnableSound);
        tvMsgNormRingtone_value       = v.findViewById(R.id.tvMsgNormRingtone_value);
        prefMsgNormRingtone           = v.findViewById(R.id.prefMsgNormRingtone);
        prefMsgNormRepeatSound        = v.findViewById(R.id.prefMsgNormRepeatSound);
        prefMsgNormEnableLED          = v.findViewById(R.id.prefMsgNormEnableLED);
        prefMsgNormLedColor_value     = v.findViewById(R.id.prefMsgNormLedColor_value);
        prefMsgNormLedColor_container = v.findViewById(R.id.prefMsgNormLedColor_container);
        prefMsgNormEnableVibrations   = v.findViewById(R.id.prefMsgNormEnableVibrations);

        prefMsgHighEnableSound        = v.findViewById(R.id.prefMsgHighEnableSound);
        tvMsgHighRingtone_value       = v.findViewById(R.id.tvMsgHighRingtone_value);
        prefMsgHighRingtone           = v.findViewById(R.id.prefMsgHighRingtone);
        prefMsgHighRepeatSound        = v.findViewById(R.id.prefMsgHighRepeatSound);
        prefMsgHighEnableLED          = v.findViewById(R.id.prefMsgHighEnableLED);
        prefMsgHighLedColor_value     = v.findViewById(R.id.prefMsgHighLedColor_value);
        prefMsgHighLedColor_container = v.findViewById(R.id.prefMsgHighLedColor_container);
        prefMsgHighEnableVibrations   = v.findViewById(R.id.prefMsgHighEnableVibrations);
    }

    private void updateUI()
    {
        SCNSettings s = SCNSettings.inst();
        Context c = getContext();
        if (c == null) return;

        prefAppEnabled.setChecked(s.Enabled);

        ArrayAdapter<Integer> plcsa = new ArrayAdapter<>(c, android.R.layout.simple_spinner_item, SCNSettings.CHOOSABLE_CACHE_SIZES);
        plcsa.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
        prefLocalCacheSize.setAdapter(plcsa);
        prefLocalCacheSize.setSelection(getCacheSizeIndex(s.LocalCacheSize));

        prefMsgLowEnableSound.setChecked(s.PriorityLow.EnableSound);
        tvMsgLowRingtone_value.setText(s.PriorityLow.SoundName);
        prefMsgLowRepeatSound.setChecked(s.PriorityLow.RepeatSound);
        prefMsgLowEnableLED.setChecked(s.PriorityLow.EnableLED);
        prefMsgLowLedColor_value.setColorFilter(s.PriorityLow.LEDColor);
        prefMsgLowEnableVibrations.setChecked(s.PriorityLow.EnableVibration);
    }

    private void initListener()
    {
        SCNSettings s = SCNSettings.inst();

        prefAppEnabled.setOnCheckedChangeListener((a,b) -> { s.Enabled=b; saveAndUpdate(); });

        prefLocalCacheSize.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener()
        {
            @Override public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                s.LocalCacheSize = prefLocalCacheSize.getSelectedItemPosition()>=0 ? SCNSettings.CHOOSABLE_CACHE_SIZES[prefLocalCacheSize.getSelectedItemPosition()] : 100;
                saveAndUpdate();
            }
            @Override public void onNothingSelected(AdapterView<?> parent) { /* */  }
        });

        prefUpgradeAccount.setOnClickListener(a -> onUpgradeAccount());

        prefMsgLowEnableSound.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableSound=b; saveAndUpdate(); });
        prefMsgLowRingtone.setOnClickListener(a -> chooseRingtoneLow());
        prefMsgLowRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.RepeatSound=b; saveAndUpdate(); });
        prefMsgLowEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableLED=b; saveAndUpdate(); });
        prefMsgLowLedColor_container.setOnClickListener(a -> chooseLEDColorLow());
        prefMsgLowEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableVibration=b; saveAndUpdate(); });
    }

    private void saveAndUpdate()
    {
        SCNSettings.inst().save();
        updateUI();
    }

    private void onUpgradeAccount()
    {
        //TODO
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
        musicPickerSwitch = 1;
        new UltimateMusicPicker()
                .windowTitle("Choose notification sound")
                .removeSilent()
                .streamType(AudioManager.STREAM_ALARM)
                .ringtone()
                .notification()
                .alarm()
                .music()
                .goWithDialog(getChildFragmentManager());
    }

    private void chooseRingtoneNorm()
    {
        musicPickerSwitch = 2;
        new UltimateMusicPicker()
                .windowTitle("Choose notification sound")
                .removeSilent()
                .streamType(AudioManager.STREAM_ALARM)
                .ringtone()
                .notification()
                .alarm()
                .music()
                .goWithDialog(getChildFragmentManager());
    }

    private void chooseRingtoneHigh()
    {
        musicPickerSwitch = 3;
        new UltimateMusicPicker()
                .windowTitle("Choose notification sound")
                .removeSilent()
                .streamType(AudioManager.STREAM_ALARM)
                .ringtone()
                .notification()
                .alarm()
                .music()
                .goWithDialog(getChildFragmentManager());
    }

    private void chooseLEDColorLow()
    {
        new ColorPickerPopup.Builder(getContext())
                .initialColor(SCNSettings.inst().PriorityLow.LEDColor) // Set initial color
                .enableBrightness(true) // Enable brightness slider or not
                .okTitle("Choose")
                .cancelTitle("Cancel")
                .showIndicator(true)
                .showValue(false)
                .build()
                .show(getView(), new ColorPickerPopup.ColorPickerObserver()
                {
                    @Override
                    public void onColorPicked(int color) {
                        SCNSettings.inst().PriorityLow.LEDColor = color;
                        saveAndUpdate();
                    }

                    @Override
                    public void onColor(int color, boolean fromUser) { }
                });
    }

    @Override
    public void onMusicPick(@NotNull Uri uri, @NotNull String s)
    {
        if (musicPickerSwitch == 1) { SCNSettings.inst().PriorityLow.SoundSource =uri.toString(); SCNSettings.inst().PriorityLow.SoundName =s; saveAndUpdate(); }
        if (musicPickerSwitch == 2) { SCNSettings.inst().PriorityNorm.SoundSource=uri.toString(); SCNSettings.inst().PriorityNorm.SoundName=s; saveAndUpdate(); }
        if (musicPickerSwitch == 3) { SCNSettings.inst().PriorityHigh.SoundSource=uri.toString(); SCNSettings.inst().PriorityHigh.SoundName=s; saveAndUpdate(); }

        musicPickerSwitch = -1;
    }

    @Override
    public void onPickCanceled()
    {
        musicPickerSwitch = -1;
    }
}