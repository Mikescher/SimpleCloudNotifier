package com.blackforestbytes.simplecloudnotifier.view;

import android.annotation.SuppressLint;
import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.media.AudioAttributes;
import android.media.AudioManager;
import android.media.MediaPlayer;
import android.net.Uri;
import android.os.Bundle;
import android.text.Editable;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.SeekBar;
import android.widget.Spinner;
import android.widget.Switch;
import android.widget.TextView;
import android.widget.Toast;

import androidx.annotation.NonNull;
import androidx.core.content.FileProvider;
import androidx.fragment.app.Fragment;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.lambda.FI;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.service.IABService;
import com.blackforestbytes.simplecloudnotifier.util.TextChangedListener;

import org.jetbrains.annotations.NotNull;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.ObjectOutputStream;
import java.util.Map;

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
    private Switch    prefEnableDeleteSwipe;
    private EditText  prefPreviewLineCount;

    private Switch    prefMsgLowEnableSound;
    private TextView  prefMsgLowRingtone_value;
    private View      prefMsgLowRingtone_container;
    private Switch    prefMsgLowRepeatSound;
    private Switch    prefMsgLowEnableLED;
    private View      prefMsgLowLedColor_container;
    private ImageView prefMsgLowLedColor_value;
    private Switch    prefMsgLowEnableVibrations;
    private Switch    prefMsgLowForceVolume;
    private SeekBar   prefMsgLowVolume;
    private ImageView prefMsgLowVolumeTest;

    private Switch    prefMsgNormEnableSound;
    private TextView  prefMsgNormRingtone_value;
    private View      prefMsgNormRingtone_container;
    private Switch    prefMsgNormRepeatSound;
    private Switch    prefMsgNormEnableLED;
    private View      prefMsgNormLedColor_container;
    private ImageView prefMsgNormLedColor_value;
    private Switch    prefMsgNormEnableVibrations;
    private Switch    prefMsgNormForceVolume;
    private SeekBar   prefMsgNormVolume;
    private ImageView prefMsgNormVolumeTest;

    private Switch    prefMsgHighEnableSound;
    private TextView  prefMsgHighRingtone_value;
    private View      prefMsgHighRingtone_container;
    private Switch    prefMsgHighRepeatSound;
    private Switch    prefMsgHighEnableLED;
    private View      prefMsgHighLedColor_container;
    private ImageView prefMsgHighLedColor_value;
    private Switch    prefMsgHighEnableVibrations;
    private Switch    prefMsgHighForceVolume;
    private SeekBar   prefMsgHighVolume;
    private ImageView prefMsgHighVolumeTest;

    private Button    prefBtnImport;
    private Button    prefBtnExport;

    private int musicPickerSwitch = -1;

    private MediaPlayer[] mPlayers = new MediaPlayer[3];

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
        prefEnableDeleteSwipe         = v.findViewById(R.id.prefEnableDeleteSwipe);
        prefPreviewLineCount          = v.findViewById(R.id.prefPreviewLineCount);

        prefMsgLowEnableSound         = v.findViewById(R.id.prefMsgLowEnableSound);
        prefMsgLowRingtone_value      = v.findViewById(R.id.prefMsgLowRingtone_value);
        prefMsgLowRingtone_container  = v.findViewById(R.id.prefMsgLowRingtone_container);
        prefMsgLowRepeatSound         = v.findViewById(R.id.prefMsgLowRepeatSound);
        prefMsgLowEnableLED           = v.findViewById(R.id.prefMsgLowEnableLED);
        prefMsgLowLedColor_value      = v.findViewById(R.id.prefMsgLowLedColor_value);
        prefMsgLowLedColor_container  = v.findViewById(R.id.prefMsgLowLedColor_container);
        prefMsgLowEnableVibrations    = v.findViewById(R.id.prefMsgLowEnableVibrations);
        prefMsgLowForceVolume         = v.findViewById(R.id.prefMsgLowForceVolume);
        prefMsgLowVolume              = v.findViewById(R.id.prefMsgLowVolume);
        prefMsgLowVolumeTest          = v.findViewById(R.id.btnLowVolumeTest);

        prefMsgNormEnableSound        = v.findViewById(R.id.prefMsgNormEnableSound);
        prefMsgNormRingtone_value     = v.findViewById(R.id.prefMsgNormRingtone_value);
        prefMsgNormRingtone_container = v.findViewById(R.id.prefMsgNormRingtone_container);
        prefMsgNormRepeatSound        = v.findViewById(R.id.prefMsgNormRepeatSound);
        prefMsgNormEnableLED          = v.findViewById(R.id.prefMsgNormEnableLED);
        prefMsgNormLedColor_value     = v.findViewById(R.id.prefMsgNormLedColor_value);
        prefMsgNormLedColor_container = v.findViewById(R.id.prefMsgNormLedColor_container);
        prefMsgNormEnableVibrations   = v.findViewById(R.id.prefMsgNormEnableVibrations);
        prefMsgNormForceVolume         = v.findViewById(R.id.prefMsgNormForceVolume);
        prefMsgNormVolume              = v.findViewById(R.id.prefMsgNormVolume);
        prefMsgNormVolumeTest          = v.findViewById(R.id.btnNormVolumeTest);

        prefMsgHighEnableSound        = v.findViewById(R.id.prefMsgHighEnableSound);
        prefMsgHighRingtone_value     = v.findViewById(R.id.prefMsgHighRingtone_value);
        prefMsgHighRingtone_container = v.findViewById(R.id.prefMsgHighRingtone_container);
        prefMsgHighRepeatSound        = v.findViewById(R.id.prefMsgHighRepeatSound);
        prefMsgHighEnableLED          = v.findViewById(R.id.prefMsgHighEnableLED);
        prefMsgHighLedColor_value     = v.findViewById(R.id.prefMsgHighLedColor_value);
        prefMsgHighLedColor_container = v.findViewById(R.id.prefMsgHighLedColor_container);
        prefMsgHighEnableVibrations   = v.findViewById(R.id.prefMsgHighEnableVibrations);
        prefMsgHighForceVolume        = v.findViewById(R.id.prefMsgHighForceVolume);
        prefMsgHighVolume             = v.findViewById(R.id.prefMsgHighVolume);
        prefMsgHighVolumeTest         = v.findViewById(R.id.btnHighVolumeTest);

        prefBtnExport                 = v.findViewById(R.id.prefExport);
        prefBtnImport                 = v.findViewById(R.id.prefImport);

        ArrayAdapter<Integer> plcsa = new ArrayAdapter<>(v.getContext(), android.R.layout.simple_spinner_item, SCNSettings.CHOOSABLE_CACHE_SIZES);
        plcsa.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
        prefLocalCacheSize.setAdapter(plcsa);
    }

    @SuppressLint("SetTextI18n")
    private void updateUI()
    {
        SCNSettings s = SCNSettings.inst();
        Context c = getContext();
        if (c == null) return;

        if (prefAppEnabled.isChecked() != s.Enabled) prefAppEnabled.setChecked(s.Enabled);
        if (prefEnableDeleteSwipe.isChecked() != s.EnableDeleteSwipe) prefEnableDeleteSwipe.setChecked(s.EnableDeleteSwipe);
        if (!prefPreviewLineCount.getText().toString().equals(Integer.toString(s.PreviewLineCount))) prefPreviewLineCount.setText(Integer.toString(s.PreviewLineCount));

        prefUpgradeAccount.setVisibility(     SCNSettings.inst().promode_local ? View.GONE    : View.VISIBLE);
        prefUpgradeAccount_info.setVisibility(SCNSettings.inst().promode_local ? View.GONE    : View.VISIBLE);
        prefUpgradeAccount_msg.setVisibility( SCNSettings.inst().promode_local ? View.VISIBLE : View.GONE   );

        if (prefLocalCacheSize.getSelectedItemPosition() != getCacheSizeIndex(s.LocalCacheSize)) prefLocalCacheSize.setSelection(getCacheSizeIndex(s.LocalCacheSize));

        if (prefMsgLowEnableSound.isChecked() != s.PriorityLow.EnableSound) prefMsgLowEnableSound.setChecked(s.PriorityLow.EnableSound);
        if (!prefMsgLowRingtone_value.getText().equals(s.PriorityLow.SoundName)) prefMsgLowRingtone_value.setText(s.PriorityLow.SoundName);
        if (prefMsgLowRepeatSound.isChecked() != s.PriorityLow.RepeatSound) prefMsgLowRepeatSound.setChecked(s.PriorityLow.RepeatSound);
        if (prefMsgLowEnableLED.isChecked() != s.PriorityLow.EnableLED) prefMsgLowEnableLED.setChecked(s.PriorityLow.EnableLED);
        prefMsgLowLedColor_value.setColorFilter(s.PriorityLow.LEDColor);
        if (prefMsgLowEnableVibrations.isChecked() != s.PriorityLow.EnableVibration) prefMsgLowEnableVibrations.setChecked(s.PriorityLow.EnableVibration);
        if (prefMsgLowForceVolume.isChecked() != s.PriorityLow.ForceVolume) prefMsgLowForceVolume.setChecked(s.PriorityLow.ForceVolume);
        if (prefMsgLowVolume.getMax() != 100) prefMsgLowVolume.setMax(100);
        if (prefMsgLowVolume.getProgress() != s.PriorityLow.ForceVolumeValue) prefMsgLowVolume.setProgress(s.PriorityLow.ForceVolumeValue);
        if (prefMsgLowVolume.isEnabled() != s.PriorityLow.ForceVolume) prefMsgLowVolume.setEnabled(s.PriorityLow.ForceVolume);
        if (prefMsgLowVolumeTest.isEnabled() != s.PriorityLow.ForceVolume) prefMsgLowVolumeTest.setEnabled(s.PriorityLow.ForceVolume);
        if (s.PriorityLow.ForceVolume) prefMsgLowVolumeTest.setColorFilter(null); else prefMsgLowVolumeTest.setColorFilter(Color.argb(150,200,200,200));

        if (prefMsgNormEnableSound.isChecked() != s.PriorityNorm.EnableSound) prefMsgNormEnableSound.setChecked(s.PriorityNorm.EnableSound);
        if (!prefMsgNormRingtone_value.getText().equals(s.PriorityNorm.SoundName)) prefMsgNormRingtone_value.setText(s.PriorityNorm.SoundName);
        if (prefMsgNormRepeatSound.isChecked() != s.PriorityNorm.RepeatSound) prefMsgNormRepeatSound.setChecked(s.PriorityNorm.RepeatSound);
        if (prefMsgNormEnableLED.isChecked() != s.PriorityNorm.EnableLED) prefMsgNormEnableLED.setChecked(s.PriorityNorm.EnableLED);
        prefMsgNormLedColor_value.setColorFilter(s.PriorityNorm.LEDColor);
        if (prefMsgNormEnableVibrations.isChecked() != s.PriorityNorm.EnableVibration) prefMsgNormEnableVibrations.setChecked(s.PriorityNorm.EnableVibration);
        if (prefMsgNormForceVolume.isChecked() != s.PriorityNorm.ForceVolume) prefMsgNormForceVolume.setChecked(s.PriorityNorm.ForceVolume);
        if (prefMsgNormVolume.getMax() != 100) prefMsgNormVolume.setMax(100);
        if (prefMsgNormVolume.getProgress() != s.PriorityNorm.ForceVolumeValue) prefMsgNormVolume.setProgress(s.PriorityNorm.ForceVolumeValue);
        if (prefMsgNormVolume.isEnabled() != s.PriorityNorm.ForceVolume) prefMsgNormVolume.setEnabled(s.PriorityNorm.ForceVolume);
        if (prefMsgNormVolumeTest.isEnabled() != s.PriorityNorm.ForceVolume) prefMsgNormVolumeTest.setEnabled(s.PriorityNorm.ForceVolume);
        if (s.PriorityNorm.ForceVolume) prefMsgNormVolumeTest.setColorFilter(null); else prefMsgNormVolumeTest.setColorFilter(Color.argb(150,200,200,200));

        if (prefMsgHighEnableSound.isChecked() != s.PriorityHigh.EnableSound) prefMsgHighEnableSound.setChecked(s.PriorityHigh.EnableSound);
        if (!prefMsgHighRingtone_value.getText().equals(s.PriorityHigh.SoundName)) prefMsgHighRingtone_value.setText(s.PriorityHigh.SoundName);
        if (prefMsgHighRepeatSound.isChecked() != s.PriorityHigh.RepeatSound) prefMsgHighRepeatSound.setChecked(s.PriorityHigh.RepeatSound);
        if (prefMsgHighEnableLED.isChecked() != s.PriorityHigh.EnableLED) prefMsgHighEnableLED.setChecked(s.PriorityHigh.EnableLED);
        prefMsgHighLedColor_value.setColorFilter(s.PriorityHigh.LEDColor);
        if (prefMsgHighEnableVibrations.isChecked() != s.PriorityHigh.EnableVibration) prefMsgHighEnableVibrations.setChecked(s.PriorityHigh.EnableVibration);
        if (prefMsgHighForceVolume.isChecked() != s.PriorityHigh.ForceVolume) prefMsgHighForceVolume.setChecked(s.PriorityHigh.ForceVolume);
        if (prefMsgHighVolume.getMax() != 100) prefMsgHighVolume.setMax(100);
        if (prefMsgHighVolume.getProgress() != s.PriorityHigh.ForceVolumeValue) prefMsgHighVolume.setProgress(s.PriorityHigh.ForceVolumeValue);
        if (prefMsgHighVolume.isEnabled() != s.PriorityHigh.ForceVolume) prefMsgHighVolume.setEnabled(s.PriorityHigh.ForceVolume);
        if (prefMsgHighVolumeTest.isEnabled() != s.PriorityHigh.ForceVolume) prefMsgHighVolumeTest.setEnabled(s.PriorityHigh.ForceVolume);
        if (s.PriorityHigh.ForceVolume) prefMsgHighVolumeTest.setColorFilter(null); else prefMsgHighVolumeTest.setColorFilter(Color.argb(150,200,200,200));
    }

    private void initListener()
    {
        SCNSettings s = SCNSettings.inst();

        prefAppEnabled.setOnCheckedChangeListener((a,b) -> { boolean prev=s.Enabled; s.Enabled=b; saveAndUpdate(); updateEnabled(prev, b); });
        prefEnableDeleteSwipe.setOnCheckedChangeListener((a,b) -> { s.EnableDeleteSwipe=b; saveAndUpdate(); });
        prefPreviewLineCount.addTextChangedListener(new TextChangedListener<EditText>(prefPreviewLineCount) {
            @Override
            public void onTextChanged(EditText target, Editable ed) {
                if (!ed.toString().isEmpty()) try { s.PreviewLineCount=Integer.parseInt(ed.toString()); saveAndUpdate(); } catch (Exception e) { /* */ }
            }
        });

        prefLocalCacheSize.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener()
        {
            @Override public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                s.LocalCacheSize = prefLocalCacheSize.getSelectedItemPosition()>=0 ? SCNSettings.CHOOSABLE_CACHE_SIZES[prefLocalCacheSize.getSelectedItemPosition()] : 100;
                saveAndUpdate();
            }
            @Override public void onNothingSelected(AdapterView<?> parent) { /* */  }
        });

        prefUpgradeAccount.setOnClickListener(a -> onUpgradeAccount());

        prefBtnExport.setOnClickListener(a -> onExport());
        prefBtnImport.setOnClickListener(a -> onImport());

        prefMsgLowEnableSound.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableSound=b; saveAndUpdate(); });
        prefMsgLowRingtone_container.setOnClickListener(a -> chooseRingtoneLow());
        prefMsgLowRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.RepeatSound=b; saveAndUpdate(); });
        prefMsgLowEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableLED=b; saveAndUpdate(); });
        prefMsgLowLedColor_container.setOnClickListener(a -> chooseLEDColorLow());
        prefMsgLowEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.EnableVibration=b; saveAndUpdate(); });
        prefMsgLowForceVolume.setOnCheckedChangeListener((a,b) -> { s.PriorityLow.ForceVolume=b; saveAndUpdate(); });
        prefMsgLowVolume.setOnSeekBarChangeListener(FI.SeekBarChanged((a,b,c) -> { if (c) { s.PriorityLow.ForceVolumeValue=b; saveAndUpdate(); updateVolume(0, b); } }));
        prefMsgLowVolumeTest.setOnClickListener((v) -> { if (s.PriorityLow.ForceVolume) playTestSound(0, prefMsgLowVolumeTest, s.PriorityLow.SoundSource, s.PriorityLow.ForceVolumeValue); });

        prefMsgNormEnableSound.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.EnableSound=b; saveAndUpdate(); });
        prefMsgNormRingtone_container.setOnClickListener(a -> chooseRingtoneNorm());
        prefMsgNormRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.RepeatSound=b; saveAndUpdate(); });
        prefMsgNormEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.EnableLED=b; saveAndUpdate(); });
        prefMsgNormLedColor_container.setOnClickListener(a -> chooseLEDColorNorm());
        prefMsgNormEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.EnableVibration=b; saveAndUpdate(); });
        prefMsgNormForceVolume.setOnCheckedChangeListener((a,b) -> { s.PriorityNorm.ForceVolume=b; saveAndUpdate(); });
        prefMsgNormVolume.setOnSeekBarChangeListener(FI.SeekBarChanged((a,b,c) -> { if (c) { s.PriorityNorm.ForceVolumeValue=b; saveAndUpdate(); updateVolume(1, b); } }));
        prefMsgNormVolumeTest.setOnClickListener((v) -> { if (s.PriorityNorm.ForceVolume) playTestSound(1, prefMsgNormVolumeTest, s.PriorityNorm.SoundSource, s.PriorityNorm.ForceVolumeValue); });

        prefMsgHighEnableSound.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.EnableSound=b; saveAndUpdate(); });
        prefMsgHighRingtone_container.setOnClickListener(a -> chooseRingtoneHigh());
        prefMsgHighRepeatSound.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.RepeatSound=b; saveAndUpdate(); });
        prefMsgHighEnableLED.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.EnableLED=b; saveAndUpdate(); });
        prefMsgHighLedColor_container.setOnClickListener(a -> chooseLEDColorHigh());
        prefMsgHighEnableVibrations.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.EnableVibration=b; saveAndUpdate(); });
        prefMsgHighForceVolume.setOnCheckedChangeListener((a,b) -> { s.PriorityHigh.ForceVolume=b; saveAndUpdate(); });
        prefMsgHighVolume.setOnSeekBarChangeListener(FI.SeekBarChanged((a,b,c) -> { if (c) { s.PriorityHigh.ForceVolumeValue=b; saveAndUpdate(); updateVolume(2, b); } }));
        prefMsgHighVolumeTest.setOnClickListener((v) -> { if (s.PriorityHigh.ForceVolume) playTestSound(2, prefMsgHighVolumeTest, s.PriorityHigh.SoundSource, s.PriorityHigh.ForceVolumeValue); });
    }

    private void onExport()
    {
        Context ctxt = getContext();
        if (ctxt == null) return;

        try
        {
            File outputDir = ctxt.getCacheDir(); // context being the Activity pointer
            File outputFile = File.createTempFile("scn_export_", ".dat", outputDir);

            ObjectOutputStream output = new ObjectOutputStream(new FileOutputStream(outputFile));

            Map<String, ?> d1 = ctxt.getSharedPreferences("Config", Context.MODE_PRIVATE).getAll();
            Map<String, ?> d2 = ctxt.getSharedPreferences("IAB", Context.MODE_PRIVATE).getAll();
            Map<String, ?> d3 = ctxt.getSharedPreferences("CMessageList", Context.MODE_PRIVATE).getAll();
            Map<String, ?> d4 = ctxt.getSharedPreferences("QueryLog", Context.MODE_PRIVATE).getAll();

            output.writeObject(d1);
            output.writeObject(d2);
            output.writeObject(d3);
            output.writeObject(d4);

            Intent intent = new Intent(Intent.ACTION_SEND);

            Uri uri = FileProvider.getUriForFile(ctxt, "com.blackforestbytes.simplecloudnotifier.fileprovider", outputFile);
            intent.putExtra(Intent.EXTRA_STREAM, uri);
            intent.setType("*/*");
            intent.setFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION);

            startActivity(Intent.createChooser(intent, "Export"));
        }
        catch (IOException e)
        {
            Log.e("Export:Err", e.toString());
            SCNApp.showToast("Export failed", Toast.LENGTH_LONG);
        }
    }

    private void onImport()
    {
        SCNApp.getMainActivity().setContentView(R.layout.activity_main);

        Intent intent = new Intent()
                .setType("*/*")
                .setAction(Intent.ACTION_GET_CONTENT);

        ((MainActivity)getActivity()).startActivityForResult(Intent.createChooser(intent, "Select a file"), 1991);
    }

    private void updateEnabled(boolean prev, boolean now)
    {
        if (!prev && now)
        {
            SCNApp.showToast("SimpleCloudNotifier is now enabled", Toast.LENGTH_SHORT);
        }
        else if (prev && !now)
        {
            SCNApp.showToast("SimpleCloudNotifier is now disabled\nYou won't recieve new messages.", Toast.LENGTH_LONG);
        }
    }

    private void updateVolume(int idx, int volume)
    {
        if (mPlayers[idx] != null && mPlayers[idx].isPlaying())
        {
            AudioManager aman = (AudioManager) SCNApp.getContext().getSystemService(Context.AUDIO_SERVICE);
            int maxVolume = aman.getStreamMaxVolume(AudioManager.STREAM_NOTIFICATION);
            aman.setStreamVolume(AudioManager.STREAM_NOTIFICATION, (int)(maxVolume * (volume / 100.0)), 0);
        }
    }

    private void stopSound(final int idx, final ImageView iv)
    {
        if (mPlayers[idx] != null && mPlayers[idx].isPlaying())
        {
            mPlayers[idx].stop();
            mPlayers[idx].release();
            iv.setImageResource(R.drawable.ic_play);
            mPlayers[idx] = null;
        }
    }

    private void playTestSound(final int idx, final ImageView iv, String src, int volume)
    {
        if (mPlayers[idx] != null && mPlayers[idx].isPlaying())
        {
            mPlayers[idx].stop();
            mPlayers[idx].release();
            iv.setImageResource(R.drawable.ic_play);
            mPlayers[idx] = null;
            return;
        }

        if (Str.isNullOrWhitespace(src)) return;
        if (volume == 0) return;

        Context ctxt = getContext();
        if (ctxt == null) return;

        iv.setImageResource(R.drawable.ic_pause);

        AudioManager aman = (AudioManager) SCNApp.getContext().getSystemService(Context.AUDIO_SERVICE);
        int maxVolume = aman.getStreamMaxVolume(AudioManager.STREAM_NOTIFICATION);
        aman.setStreamVolume(AudioManager.STREAM_NOTIFICATION, (int)(maxVolume * (volume / 100.0)), 0);

        MediaPlayer player = mPlayers[idx] = new MediaPlayer();
        player.setAudioAttributes(new AudioAttributes.Builder().setLegacyStreamType(AudioManager.STREAM_NOTIFICATION).build());
        player.setAudioStreamType(AudioManager.STREAM_NOTIFICATION);

        try
        {
            player.setDataSource(ctxt, Uri.parse(src));
            player.setLooping(false);
            player.setOnCompletionListener(  mp -> SCNApp.runOnUiThread(() -> { mp.stop(); iv.setImageResource(R.drawable.ic_play); mPlayers[idx]=null; mp.release(); }));
            player.setOnSeekCompleteListener(mp -> SCNApp.runOnUiThread(() -> { mp.stop(); iv.setImageResource(R.drawable.ic_play); mPlayers[idx]=null; mp.release(); }));
            player.prepare();
            player.start();
        }
        catch (IOException e)
        {
            Log.e("SFRAG:play", e.toString());
        }
    }

    private void saveAndUpdate()
    {
        SCNSettings.inst().save();
        updateUI();
        SCNApp.getMainActivity().adpTabs.tab1.updateDeleteSwipeEnabled();
    }

    private void onUpgradeAccount()
    {
        IABService.inst().purchase(getActivity(), IABService.IAB_PRO_MODE);
    }

    public void updateProState()
    {
        boolean pmode = IABService.inst().getPurchaseCachedSimple(IABService.IAB_PRO_MODE);

        if (prefUpgradeAccount != null)      prefUpgradeAccount.setVisibility(     pmode ? View.GONE    : View.VISIBLE);
        if (prefUpgradeAccount_info != null) prefUpgradeAccount_info.setVisibility(pmode ? View.GONE    : View.VISIBLE);
        if (prefUpgradeAccount_msg != null)  prefUpgradeAccount_msg.setVisibility( pmode ? View.VISIBLE : View.GONE   );
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
        ump.streamType(AudioManager.STREAM_NOTIFICATION);
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
        ump.streamType(AudioManager.STREAM_NOTIFICATION);
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
        ump.streamType(AudioManager.STREAM_NOTIFICATION);
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

    public void onViewpagerHide()
    {
        stopSound(0, prefMsgLowVolumeTest);
        stopSound(1, prefMsgNormVolumeTest);
        stopSound(2, prefMsgHighVolumeTest);
    }
}