package com.blackforestbytes.simplecloudnotifier.view;

import android.annotation.SuppressLint;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.support.annotation.NonNull;
import android.support.v4.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageButton;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;

import net.glxn.qrgen.android.QRCode;
import net.glxn.qrgen.core.image.ImageType;

import static android.content.Context.CLIPBOARD_SERVICE;

public class AccountFragment extends Fragment
{
    public AccountFragment()
    {
        // Required empty public constructor
    }

    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState)
    {
        View v = inflater.inflate(R.layout.fragment_account, container, false);

        updateUI(v);

        v.findViewById(R.id.btnCopyUserID).setOnClickListener(cv ->
        {
            ClipboardManager clipboard = (ClipboardManager) cv.getContext().getSystemService(CLIPBOARD_SERVICE);
            clipboard.setPrimaryClip(ClipData.newPlainText("UserID", String.valueOf(SCNSettings.inst().user_id)));
            SCNApp.showToast("Copied userID to clipboard", 1000);
        });

        v.findViewById(R.id.btnCopyUserKey).setOnClickListener(cv ->
        {
            ClipboardManager clipboard = (ClipboardManager) cv.getContext().getSystemService(CLIPBOARD_SERVICE);
            clipboard.setPrimaryClip(ClipData.newPlainText("UserKey", String.valueOf(SCNSettings.inst().user_key)));
            SCNApp.showToast("Copied key to clipboard", 1000);
        });

        v.findViewById(R.id.btnAccountReset).setOnClickListener(cv ->
        {
            View lpnl = v.findViewById(R.id.loadingPanel);
            lpnl.setVisibility(View.VISIBLE);
            SCNSettings.inst().reset(lpnl);
        });

        v.findViewById(R.id.btnClearLocalStorage).setOnClickListener(cv ->
        {
            CMessageList.inst().clear();
            SCNApp.showToast("Notifications cleared", 1000);
        });

        v.findViewById(R.id.btnQR).setOnClickListener(cv ->
        {
            Intent browserIntent = new Intent(Intent.ACTION_VIEW, Uri.parse(SCNSettings.inst().createOnlineURL()));
            startActivity(browserIntent);
        });

        v.findViewById(R.id.btnRefresh).setOnClickListener(cv ->
        {
            View lpnl = v.findViewById(R.id.loadingPanel);
            lpnl.setVisibility(View.VISIBLE);
            SCNSettings.inst().refresh(lpnl, getActivity());
        });

        return v;
    }

    public void updateUI()
    {
        updateUI(getView());
    }

    @SuppressLint("DefaultLocale")
    public void updateUI(View v)
    {
        if (v == null) return;
        TextView tvUserID  = v.findViewById(R.id.tvUserID);
        TextView tvUserKey = v.findViewById(R.id.tvUserKey);
        TextView tvQuota   = v.findViewById(R.id.tvQuota);
        ImageButton btnQR  = v.findViewById(R.id.btnQR);

        SCNSettings s = SCNSettings.inst();

        if (s.isConnected())
        {
            tvUserID.setText(String.valueOf(s.user_id));
            tvUserKey.setText(s.user_key);
            tvQuota.setText(String.format("%d / %d", s.quota_curr, s.quota_max));
            btnQR.setImageBitmap(QRCode.from(s.createOnlineURL()).to(ImageType.PNG).withSize(512, 512).bitmap());
        }
        else
        {
            tvUserID.setText(R.string.str_not_connected);
            tvUserKey.setText(R.string.str_not_connected);
            tvQuota.setText(R.string.str_not_connected);
            btnQR.setImageResource(R.drawable.qr_default);
        }
    }
}
