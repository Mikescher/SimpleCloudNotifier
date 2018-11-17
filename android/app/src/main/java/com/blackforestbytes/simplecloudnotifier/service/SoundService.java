package com.blackforestbytes.simplecloudnotifier.service;

import android.content.Context;
import android.media.AudioAttributes;
import android.media.AudioManager;
import android.media.MediaPlayer;
import android.net.Uri;
import android.util.Log;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;

import java.io.IOException;

public class SoundService
{
    private static MediaPlayer mpLast = null;

    public static void play(boolean enableSound, String soundSource, boolean forceVolume, int forceVolumeValue, boolean loop)
    {
        if (!enableSound) return;
        if (Str.isNullOrWhitespace(soundSource)) return;

        stop();

        if (forceVolume)
        {
            AudioManager aman = (AudioManager) SCNApp.getContext().getSystemService(Context.AUDIO_SERVICE);
            int maxVolume = aman.getStreamMaxVolume(AudioManager.STREAM_NOTIFICATION);
            aman.setStreamVolume(AudioManager.STREAM_NOTIFICATION, (int)(maxVolume * (forceVolumeValue / 100.0)), 0);
        }

        try
        {
            MediaPlayer player = new MediaPlayer();
            player.setAudioAttributes(new AudioAttributes.Builder().setLegacyStreamType(AudioManager.STREAM_NOTIFICATION).build());
            player.setAudioStreamType(AudioManager.STREAM_NOTIFICATION);
            player.setDataSource(SCNApp.getContext(), Uri.parse(soundSource));
            player.setLooping(loop);
            player.setOnCompletionListener(  mp -> { mp.stop(); mp.release(); });
            player.setOnSeekCompleteListener(mp -> { mp.stop(); mp.release(); });
            player.prepare();
            player.start();
            mpLast = player;
        }
        catch (IOException e)
        {
            Log.e("Sound::play", e.toString());
        }
    }

    public static void stop()
    {
        if (mpLast != null && mpLast.isPlaying()) { mpLast.stop(); mpLast.release(); mpLast = null; }
    }
}
