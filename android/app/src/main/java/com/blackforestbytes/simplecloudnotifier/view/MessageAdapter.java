package com.blackforestbytes.simplecloudnotifier.view;

import android.support.annotation.NonNull;
import android.support.v7.widget.RecyclerView;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;
import android.widget.Toast;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;

public class MessageAdapter extends RecyclerView.Adapter
{
    public MessageAdapter()
    {
        CMessageList.inst().register(this);
    }

    @NonNull
    @Override
    public RecyclerView.ViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType)
    {
        View myView = LayoutInflater.from(parent.getContext()).inflate(R.layout.message_card, parent, false);
        return new MessagePresenter(myView);
    }

    @Override
    public void onBindViewHolder(@NonNull RecyclerView.ViewHolder holder, int position)
    {
        CMessage msg = CMessageList.inst().tryGet(position);
        MessagePresenter view = (MessagePresenter) holder;
        view.setMessage(msg);
    }

    @Override
    public int getItemCount()
    {
        return CMessageList.inst().size();
    }

    private class MessagePresenter extends RecyclerView.ViewHolder implements View.OnClickListener
    {
        private TextView tvTimestamp;
        private TextView tvTitle;
        private TextView tvMessage;

        private CMessage data;

        MessagePresenter(View itemView)
        {
            super(itemView);
            tvTimestamp = itemView.findViewById(R.id.tvTimestamp);
            tvTitle = itemView.findViewById(R.id.tvTitle);
            tvMessage = itemView.findViewById(R.id.tvMessage);
            itemView.setOnClickListener(this);
        }

        void setMessage(CMessage msg)
        {
            tvTimestamp.setText(msg.formatTimestamp());
            tvTitle.setText(msg.Title);
            tvMessage.setText(msg.Content);
            data = msg;
        }


        @Override
        public void onClick(View v)
        {
            SCNApp.showToast(data.Title, Toast.LENGTH_LONG);
        }
    }
}
