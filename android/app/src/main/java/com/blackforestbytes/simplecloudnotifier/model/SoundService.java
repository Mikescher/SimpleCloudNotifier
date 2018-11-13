package com.blackforestbytes.simplecloudnotifier.model;

import android.content.Context;
import android.media.AudioManager;
import android.media.MediaPlayer;
import android.media.Ringtone;
import android.media.RingtoneManager;
import android.net.Uri;
import android.os.Build;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.android.ThreadUtils;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;

public class SoundService
{
    public static void playForegroundNoLooping(boolean enableSound, String soundSource, boolean forceVolume, int forceVolumeValue)
    {
        if (!enableSound) return;
        if (Str.isNullOrWhitespace(soundSource)) return;

        if (forceVolume)
        {
            AudioManager aman = (AudioManager) SCNApp.getContext().getSystemService(Context.AUDIO_SERVICE);
            int maxVolume = aman.getStreamMaxVolume(AudioManager.STREAM_MUSIC);
            aman.setStreamVolume(AudioManager.STREAM_MUSIC, (int)(maxVolume * (forceVolumeValue / 100.0)), 0);

            MediaPlayer player = MediaPlayer.create(SCNApp.getMainActivity(), Uri.parse(soundSource));
            player.setLooping(false);
            player.setOnCompletionListener(  mp -> { mp.stop(); mp.release(); });
            player.setOnSeekCompleteListener(mp -> { mp.stop(); mp.release(); });
            player.start();
        }
        else
        {
            Ringtone rt = RingtoneManager.getRingtone(SCNApp.getContext(), Uri.parse(soundSource));
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.P) rt.setLooping(false);
            rt.play();

            new Thread(() -> { ThreadUtils.safeSleep(5*1000); rt.stop(); }).start();
        }
    }
}
