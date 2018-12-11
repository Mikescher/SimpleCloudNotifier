package com.blackforestbytes.simplecloudnotifier.view;

import android.content.Intent;
import android.graphics.Color;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.blackforestbytes.simplecloudnotifier.R;
import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.model.CMessage;
import com.blackforestbytes.simplecloudnotifier.model.CMessageList;
import com.google.android.material.button.MaterialButton;

import java.lang.ref.WeakReference;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import java.util.WeakHashMap;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

public class MessageAdapter extends RecyclerView.Adapter
{
    private final View vNoElements;
    private final LinearLayoutManager manLayout;
    private final RecyclerView viewRecycler;

    private WeakHashMap<MessagePresenter, Boolean> viewHolders = new WeakHashMap<>();

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
        view.setMessage(msg, position);

        viewHolders.put(view, true);
    }

    @Override
    public void onViewRecycled(@NonNull RecyclerView.ViewHolder holder)
    {
        if (holder instanceof MessagePresenter) viewHolders.remove(holder);
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

    public CMessage removeItem(int position)
    {
        CMessage i = CMessageList.inst().removeFromBack(position);
        notifyDataSetChanged();
        return i;
    }

    public void restoreItem(CMessage item, int position)
    {
        CMessageList.inst().insert(position, item);
        notifyDataSetChanged();
    }

    public class MessagePresenter extends RecyclerView.ViewHolder implements View.OnClickListener
    {
        private TextView tvTimestamp;
        private TextView tvTitle;
        private TextView tvMessage;
        private ImageView ivPriority;

        public RelativeLayout viewForeground;
        public RelativeLayout viewBackground;

        public MaterialButton btnShare;
        public MaterialButton btnDelete;

        private CMessage data;
        private int datapos;

        MessagePresenter(View itemView)
        {
            super(itemView);
            tvTimestamp    = itemView.findViewById(R.id.tvTimestamp);
            tvTitle        = itemView.findViewById(R.id.tvTitle);
            tvMessage      = itemView.findViewById(R.id.tvMessage);
            ivPriority     = itemView.findViewById(R.id.ivPriority);
            viewForeground = itemView.findViewById(R.id.layoutFront);
            viewBackground = itemView.findViewById(R.id.layoutBack);
            btnShare       = itemView.findViewById(R.id.btnShare);
            btnDelete      = itemView.findViewById(R.id.btnDelete);

            itemView.setOnClickListener(this);
            tvTimestamp.setOnClickListener(this);
            tvTitle.setOnClickListener(this);
            tvMessage.setOnClickListener(this);
            ivPriority.setOnClickListener(this);
            viewForeground.setOnClickListener(this);

            btnShare.setOnClickListener(v ->
            {
                if (data == null) return;
                Intent sharingIntent = new Intent(android.content.Intent.ACTION_SEND);
                sharingIntent.setType("text/plain");
                sharingIntent.putExtra(android.content.Intent.EXTRA_SUBJECT, data.Title);
                sharingIntent.putExtra(android.content.Intent.EXTRA_TEXT, data.Content);
                SCNApp.getMainActivity().startActivity(Intent.createChooser(sharingIntent, "Share message"));

            });
            btnDelete.setOnClickListener(v -> { if (data != null) SCNApp.getMainActivity().adpTabs.tab1.deleteMessage(datapos); });

        }

        void setMessage(CMessage msg, int pos)
        {
            tvTimestamp.setText(msg.formatTimestamp());
            tvTitle.setText(msg.Title);
            tvMessage.setText(msg.Content);

            switch (msg.Priority)
            {
                case LOW:
                    ivPriority.setVisibility(View.VISIBLE);
                    ivPriority.setImageResource(R.drawable.priority_low);
                    ivPriority.setColorFilter(Color.rgb(176, 176, 176));
                    break;
                case NORMAL:
                    ivPriority.setVisibility(View.GONE);
                    ivPriority.setColorFilter(Color.rgb(176, 176, 176));
                    break;
                case HIGH:
                    ivPriority.setVisibility(View.VISIBLE);
                    ivPriority.setImageResource(R.drawable.priority_high);
                    ivPriority.setColorFilter(Color.rgb(200, 0, 0));
                    break;
            }

            data = msg;
            datapos = pos;

            if (msg.IsExpandedInAdapter) expand(true); else collapse(true);
        }

        private void expand(boolean force)
        {
            if (data != null && data.IsExpandedInAdapter && !force) return;
            if (data != null) data.IsExpandedInAdapter = true;
            if (tvMessage != null) tvMessage.setMaxLines(999);
            if (btnDelete != null) btnDelete.setVisibility(View.VISIBLE);
            if (btnShare != null) btnShare.setVisibility(View.VISIBLE);

        }

        private void collapse(boolean force)
        {
            if (data != null && !data.IsExpandedInAdapter && !force) return;
            if (data != null) data.IsExpandedInAdapter = false;
            if (tvMessage != null) tvMessage.setMaxLines(6);
            if (btnDelete != null) btnDelete.setVisibility(View.GONE);
            if (btnShare != null) btnShare.setVisibility(View.GONE);
        }

        @Override
        public void onClick(View v)
        {
            if (data.IsExpandedInAdapter)
            {
                collapse(false);
                return;
            }

            for (MessagePresenter holder : MessageAdapter.this.viewHolders.keySet())
            {
                if (holder == null) continue;
                if (holder == this) continue;
                holder.collapse(false);
            }

            expand(false);
        }
    }
}
