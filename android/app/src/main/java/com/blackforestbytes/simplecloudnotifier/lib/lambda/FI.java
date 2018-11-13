package com.blackforestbytes.simplecloudnotifier.lib.lambda;

import android.widget.SeekBar;

public final class FI
{
    private FI() throws InstantiationException { throw new InstantiationException(); }

    public static SeekBar.OnSeekBarChangeListener SeekBarChanged(Func3to0<SeekBar, Integer, Boolean> action)
    {
        return new SeekBar.OnSeekBarChangeListener() {
            @Override
            public void onProgressChanged(SeekBar seekBar, int progress, boolean fromUser) {
                action.invoke(seekBar, progress, fromUser);
            }

            @Override
            public void onStartTrackingTouch(SeekBar seekBar) {

            }

            @Override
            public void onStopTrackingTouch(SeekBar seekBar) {

            }
        };
    }

}
