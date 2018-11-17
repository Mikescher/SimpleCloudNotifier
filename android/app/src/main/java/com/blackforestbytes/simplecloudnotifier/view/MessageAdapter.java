package com.blackforestbytes.simplecloudnotifier.view;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

public class MessageAdapter extends RecyclerView.Adapter
{
    private final View vNoElements;
    private final LinearLayoutManager manLayout;
    private final RecyclerView viewRecycler;

    public MessageAdapter(View noElementsView, LinearLayoutManager layout, RecyclerView recycler)
    {
        vNoElements  = noElementsView;
        manLayout    = layout;
        viewRecycler = recycler;
        CMessageList.inst().register(this);

        vNoElements.setVisibility(getItemCount()>0 ? View.GONE : View.VISIBLE);
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
        CMessage msg = CMessageList.inst().tryGetFromBack(position);
        MessagePresenter view = (MessagePresenter) holder;
        view.setMessage(msg);
    }

    @Override
    public int getItemCount()
    {
        return CMessageList.inst().size();
    }

    public void customNotifyItemInserted(int idx)
    {
        notifyItemInserted(idx);
        vNoElements.setVisibility(getItemCount()>0 ? View.GONE : View.VISIBLE);
    }

    public void customNotifyDataSetChanged()
    {
        notifyDataSetChanged();
        vNoElements.setVisibility(getItemCount()>0 ? View.GONE : View.VISIBLE);
    }

    public void scrollToTop()
    {
        manLayout.smoothScrollToPosition(viewRecycler, null, 0);
    }


    public void removeItem(int position)
    {
        CMessageList.inst().remove(position);
        notifyItemRemoved(position);
    }

    public void restoreItem(CMessage item, int position)
    {
        CMessageList.inst().insert(position, item);
        notifyItemInserted(position);
    }

    public class MessagePresenter extends RecyclerView.ViewHolder implements View.OnClickListener
    {
        private TextView tvTimestamp;
        private TextView tvTitle;
        private TextView tvMessage;
        private ImageView ivPriority;

        public RelativeLayout viewForeground;
        public RelativeLayout viewBackground;

        private CMessage data;

        MessagePresenter(View itemView)
        {
            super(itemView);
            tvTimestamp    = itemView.findViewById(R.id.tvTimestamp);
            tvTitle        = itemView.findViewById(R.id.tvTitle);
            tvMessage      = itemView.findViewById(R.id.tvMessage);
            ivPriority     = itemView.findViewById(R.id.ivPriority);
            viewForeground = itemView.findViewById(R.id.layoutFront);
            viewBackground = itemView.findViewById(R.id.layoutBack);
            itemView.setOnClickListener(this);
        }

        void setMessage(CMessage msg)
        {
            tvTimestamp.setText(msg.formatTimestamp());
            tvTitle.setText(msg.Title);
            tvMessage.setText(msg.Content);

            switch (msg.Priority)
            {
                case LOW:
                    ivPriority.setVisibility(View.VISIBLE);
                    ivPriority.setImageResource(R.drawable.priority_low);
                    break;
                case NORMAL:
                    ivPriority.setVisibility(View.GONE);
                    break;
                case HIGH:
                    ivPriority.setVisibility(View.VISIBLE);
                    ivPriority.setImageResource(R.drawable.priority_high);
                    break;
            }

            data = msg;
        }

        @Override
        public void onClick(View v)
        {
            //SCNApp.showToast(data.Title, Toast.LENGTH_LONG);
        }
    }
}
