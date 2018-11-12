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

import com.android.billingclient.api.Purchase;
import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.blackforestbytes.simplecloudnotifier.service.NotificationService;

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
    private TextView  prefUpgradeAccount_msg;
    private TextView  prefUpgradeAccount_info;

    private Switch    prefMsgLowEnableSound;
    private TextView  prefMsgLowRingtone_value;
    private View      prefMsgLowRingtone_container;
    private Switch    prefMsgLowRepeatSound;
    private Switch    prefMsgLowEnableLED;
    private View      prefMsgLowLedColor_container;
    private ImageView prefMsgLowLedColor_value;
    private Switch    prefMsgLowEnableVibrations;

    private Switch    prefMsgNormEnableSound;
    private TextView  prefMsgNormRingtone_value;
    private View      prefMsgNormRingtone_container;
    private Switch    prefMsgNormRepeatSound;
    private Switch    prefMsgNormEnableLED;
    private View      prefMsgNormLedColor_container;
    private ImageView prefMsgNormLedColor_value;
    private Switch    prefMsgNormEnableVibrations;

    private Switch    prefMsgHighEnableSound;
    private TextView  prefMsgHighRingtone_value;
    private View      prefMsgHighRingtone_container;
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
        prefUpgradeAccount_msg        = v.findViewById(R.id.prefUpgradeAccount2);
        prefUpgradeAccount_info       = v.findViewById(R.id.prefUpgradeAccount_info);

        prefMsgLowEnableSound         = v.findViewById(R.id.prefMsgLowEnableSound);
        prefMsgLowRingtone_value      = v.findViewById(R.id.prefMsgLowRingtone_value);
        prefMsgLowRingtone_container  = v.findViewById(R.id.prefMsgLowRingtone_container);
        prefMsgLowRepeatSound         = v.findViewById(R.id.prefMsgLowRepeatSound);
        prefMsgLowEnableLED           = v.findViewById(R.id.prefMsgLowEnableLED);
        prefMsgLowLedColor_value      = v.findViewById(R.id.prefMsgLowLedColor_value);
        prefMsgLowLedColor_container  = v.findViewById(R.id.prefMsgLowLedColor_container);
        prefMsgLowEnableVibrations    = v.findViewById(R.id.prefMsgLowEnableVibrations);

        prefMsgNormEnableSound        = v.findViewById(R.id.prefMsgNormEnableSound);
        prefMsgNormRingtone_value     = v.findViewById(R.id.prefMsgNormRingtone_value);
        prefMsgNormRingtone_container = v.findViewById(R.id.prefMsgNormRingtone_container);
        prefMsgNormRepeatSound        = v.findViewById(R.id.prefMsgNormRepeatSound);
        prefMsgNormEnableLED          = v.findViewById(R.id.prefMsgNormEnableLED);
        prefMsgNormLedColor_value     = v.findViewById(R.id.prefMsgNormLedColor_value);
        prefMsgNormLedColor_container = v.findViewById(R.id.prefMsgNormLedColor_container);
        prefMsgNormEnableVibrations   = v.findViewById(R.id.prefMsgNormEnableVibrations);

        prefMsgHighEnableSound        = v.findViewById(R.id.prefMsgHighEnableSound);
        prefMsgHighRingtone_value     = v.findViewById(R.id.prefMsgHighRingtone_value);
        prefMsgHighRingtone_container = v.findViewById(R.id.prefMsgHighRingtone_container);
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

        if (prefAppEnabled.isChecked() != s.Enabled) prefAppEnabled.setChecked(s.Enabled);

        prefUpgradeAccount.setVisibility(     SCNSettings.inst().promode_local ? View.GONE    : View.VISIBLE);
        prefUpgradeAccount_info.setVisibility(SCNSettings.inst().promode_local ? View.GONE    : View.VISIBLE);
        prefUpgradeAccount_msg.setVisibility( SCNSettings.inst().promode_local ? View.VISIBLE : View.GONE   );

        ArrayAdapter<Integer> plcsa = new ArrayAdapter<>(c, android.R.layout.simple_spinner_item, SCNSettings.CHOOSABLE_CACHE_SIZES);
        plcsa.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
        prefLocalCacheSize.setAdapter(plcsa);
        prefLocalCacheSize.setSelection(getCacheSizeIndex(s.LocalCacheSize));

        if (prefMsgLowEnableSound.isChecked() != s.PriorityLow.EnableSound) prefMsgLowEnableSound.setChecked(s.PriorityLow.EnableSound);
        if (!prefMsgLowRingtone_value.getText().equals(s.PriorityLow.SoundName)) prefMsgLowRingtone_value.setText(s.PriorityLow.SoundName);
        if (prefMsgLowRepeatSound.isChecked() != s.PriorityLow.RepeatSound) prefMsgLowRepeatSound.setChecked(s.PriorityLow.RepeatSound);
        if (prefMsgLowEnableLED.isChecked() != s.PriorityLow.EnableLED) prefMsgLowEnableLED.setChecked(s.PriorityLow.EnableLED);
        prefMsgLowLedColor_value.setColorFilter(s.PriorityLow.LEDColor);
        if (prefMsgLowEnableVibrations.isChecked() != s.PriorityLow.EnableVibration) prefMsgLowEnableVibrations.setChecked(s.PriorityLow.EnableVibration);

        if (prefMsgNormEnableSound.isChecked() != s.PriorityNorm.EnableSound) prefMsgNormEnableSound.setChecked(s.PriorityNorm.EnableSound);
        if (!prefMsgNormRingtone_value.getText().equals(s.PriorityNorm.SoundName)) prefMsgNormRingtone_value.setText(s.PriorityNorm.SoundName);
        if (prefMsgNormRepeatSound.isChecked() != s.PriorityNorm.RepeatSound) prefMsgNormRepeatSound.setChecked(s.PriorityNorm.RepeatSound);
        if (prefMsgNormEnableLED.isChecked() != s.PriorityNorm.EnableLED) prefMsgNormEnableLED.setChecked(s.PriorityNorm.EnableLED);
        prefMsgNormLedColor_value.setColorFilter(s.PriorityNorm.LEDColor);
        if (prefMsgNormEnableVibrations.isChecked() != s.PriorityNorm.EnableVibration) prefMsgNormEnableVibrations.setChecked(s.PriorityNorm.EnableVibration);

        if (prefMsgHighEnableSound.isChecked() != s.PriorityHigh.EnableSound) prefMsgHighEnableSound.setChecked(s.PriorityHigh.EnableSound);
        if (!prefMsgHighRingtone_value.getText().equals(s.PriorityHigh.SoundName)) prefMsgHighRingtone_value.setText(s.PriorityHigh.SoundName);
        if (prefMsgHighRepeatSound.isChecked() != s.PriorityHigh.RepeatSound) prefMsgHighRepeatSound.setChecked(s.PriorityHigh.RepeatSound);
        if (prefMsgHighEnableLED.isChecked() != s.PriorityHigh.EnableLED) prefMsgHighEnableLED.setChecked(s.PriorityHigh.EnableLED);
        prefMsgHighLedColor_value.setColorFilter(s.PriorityHigh.LEDColor);
        if (prefMsgHighEnableVibrations.isChecked() != s.PriorityHigh.EnableVibration) prefMsgHighEnableVibrations.setChecked(s.PriorityHigh.EnableVibration);
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
        prefMsgLowRingtone_container.setOnClickListener(a -> chooseRingtoneLow());
        prefMsgLowRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.RepeatSound=b; saveAndUpdate(); });
        prefMsgLowEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableLED=b; saveAndUpdate(); });
        prefMsgLowLedColor_container.setOnClickListener(a -> chooseLEDColorLow());
        prefMsgLowEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableVibration=b; saveAndUpdate(); });

        prefMsgNormEnableSound.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.EnableSound=b; saveAndUpdate(); });
        prefMsgNormRingtone_container.setOnClickListener(a -> chooseRingtoneNorm());
        prefMsgNormRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.RepeatSound=b; saveAndUpdate(); });
        prefMsgNormEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.EnableLED=b; saveAndUpdate(); });
        prefMsgNormLedColor_container.setOnClickListener(a -> chooseLEDColorNorm());
        prefMsgNormEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.EnableVibration=b; saveAndUpdate(); });

        prefMsgHighEnableSound.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.EnableSound=b; saveAndUpdate(); });
        prefMsgHighRingtone_container.setOnClickListener(a -> chooseRingtoneHigh());
        prefMsgHighRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.RepeatSound=b; saveAndUpdate(); });
        prefMsgHighEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.EnableLED=b; saveAndUpdate(); });
        prefMsgHighLedColor_container.setOnClickListener(a -> chooseLEDColorHigh());
        prefMsgHighEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.EnableVibration=b; saveAndUpdate(); });
    }

    private void saveAndUpdate()
    {
        SCNSettings.inst().save();
        updateUI();
        NotificationService.inst().updateChannels();
    }

    private void onUpgradeAccount()
    {
        IABService.inst().purchase(getActivity(), IABService.IAB_PRO_MODE);
    }

    public void updateProState()
    {
        Purchase p = IABService.inst().getPurchaseCached(IABService.IAB_PRO_MODE);

        if (prefUpgradeAccount != null)      prefUpgradeAccount.setVisibility(     p != null ? View.GONE    : View.VISIBLE);
        if (prefUpgradeAccount_info != null) prefUpgradeAccount_info.setVisibility(p != null ? View.GONE    : View.VISIBLE);
        if (prefUpgradeAccount_msg != null)  prefUpgradeAccount_msg.setVisibility( p != null ? View.VISIBLE : View.GONE   );
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
        UltimateMusicPicker ump = new UltimateMusicPicker();
        ump.windowTitle("Choose notification sound");
        ump.removeSilent();
        ump.streamType(AudioManager.STREAM_ALARM);
        ump.ringtone();
        ump.notification();
        ump.alarm();
        ump.music();
        if (!SCNSettings.inst().PriorityLow.SoundSource.isEmpty())ump.selectUri(Uri.parse(SCNSettings.inst().PriorityLow.SoundSource));
        ump.goWithDialog(getChildFragmentManager());
    }

    private void chooseRingtoneNorm()
    {
        musicPickerSwitch = 2;
        UltimateMusicPicker ump = new UltimateMusicPicker();
        ump.windowTitle("Choose notification sound");
        ump.removeSilent();
        ump.streamType(AudioManager.STREAM_ALARM);
        ump.ringtone();
        ump.notification();
        ump.alarm();
        ump.music();
        if (!SCNSettings.inst().PriorityNorm.SoundSource.isEmpty())ump.defaultUri(Uri.parse(SCNSettings.inst().PriorityNorm.SoundSource));
        ump.goWithDialog(getChildFragmentManager());
    }

    private void chooseRingtoneHigh()
    {
        musicPickerSwitch = 3;
        UltimateMusicPicker ump = new UltimateMusicPicker();
        ump.windowTitle("Choose notification sound");
        ump.removeSilent();
        ump.streamType(AudioManager.STREAM_ALARM);
        ump.ringtone();
        ump.notification();
        ump.alarm();
        ump.music();
        if (!SCNSettings.inst().PriorityHigh.SoundSource.isEmpty())ump.defaultUri(Uri.parse(SCNSettings.inst().PriorityHigh.SoundSource));
        ump.goWithDialog(getChildFragmentManager());
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

    private void chooseLEDColorNorm()
    {
        new ColorPickerPopup.Builder(getContext())
                .initialColor(SCNSettings.inst().PriorityNorm.LEDColor) // Set initial color
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
                        SCNSettings.inst().PriorityNorm.LEDColor = color;
                        saveAndUpdate();
                    }

                    @Override
                    public void onColor(int color, boolean fromUser) { }
                });
    }

    private void chooseLEDColorHigh()
    {
        new ColorPickerPopup.Builder(getContext())
                .initialColor(SCNSettings.inst().PriorityHigh.LEDColor) // Set initial color
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
                        SCNSettings.inst().PriorityHigh.LEDColor = color;
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