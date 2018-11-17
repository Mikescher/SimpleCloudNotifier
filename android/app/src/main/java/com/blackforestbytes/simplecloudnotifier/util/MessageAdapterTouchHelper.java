package com.blackforestbytes.simplecloudnotifier.util;

import android.graphics.Canvas;
import android.view.View;

import com.blackforestbytes.simplecloudnotifier.model.SCNSettings;
import com.blackforestbytes.simplecloudnotifier.view.MessageAdapter;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.ItemTouchHelper;
import androidx.recyclerview.widget.RecyclerView;

public class MessageAdapterTouchHelper extends ItemTouchHelper.SimpleCallback
{
    private MessageAdapterTouchHelperListener listener;

    private int dir = 0;

    public MessageAdapterTouchHelper(int dragDirs, int swipeDirs, MessageAdapterTouchHelperListener listener)
    {
        super(dragDirs, swipeDirs);
        this.dir = swipeDirs;
        this.listener = listener;
        updateEnabled();
    }

    public void updateEnabled()
    {
        int sdir = SCNSettings.inst().EnableDeleteSwipe ? ItemTouchHelper.LEFT : 0;
        if (dir == sdir) return;
        setDefaultSwipeDirs(dir = sdir);
    }

    @Override
    public boolean onMove(@NonNull RecyclerView recyclerView, @NonNull RecyclerView.ViewHolder viewHolder, @NonNull RecyclerView.ViewHolder target)
    {
        return true;
    }

    @Override
    public void onSelectedChanged(RecyclerView.ViewHolder viewHolder, int actionState)
    {
        if (viewHolder != null)
        {
            final View foregroundView = ((MessageAdapter.MessagePresenter) viewHolder).viewForeground;

            getDefaultUIUtil().onSelected(foregroundView);
        }
    }

    @Override
    public void onChildDrawOver(@NonNull Canvas c, @NonNull RecyclerView recyclerView, RecyclerView.ViewHolder viewHolder, float dX, float dY, int actionState, boolean isCurrentlyActive)
    {
        final View foregroundView = ((MessageAdapter.MessagePresenter) viewHolder).viewForeground;
        getDefaultUIUtil().onDrawOver(c, recyclerView, foregroundView, dX, dY, actionState, isCurrentlyActive);
    }

    @Override
    public void clearView(@NonNull RecyclerView recyclerView, @NonNull RecyclerView.ViewHolder viewHolder)
    {
        final View foregroundView = ((MessageAdapter.MessagePresenter) viewHolder).viewForeground;
        getDefaultUIUtil().clearView(foregroundView);
    }

    @Override
    public void onChildDraw(@NonNull Canvas c, @NonNull RecyclerView recyclerView, @NonNull RecyclerView.ViewHolder viewHolder, float dX, float dY, int actionState, boolean isCurrentlyActive)
    {
        final View foregroundView = ((MessageAdapter.MessagePresenter) viewHolder).viewForeground;

        getDefaultUIUtil().onDraw(c, recyclerView, foregroundView, dX, dY, actionState, isCurrentlyActive);
    }

    @Override
    public void onSwiped(@NonNull RecyclerView.ViewHolder viewHolder, int direction)
    {
        listener.onSwiped(viewHolder, direction, viewHolder.getAdapterPosition());
    }

    @Override
    public int convertToAbsoluteDirection(int flags, int layoutDirection)
    {
        return super.convertToAbsoluteDirection(flags, layoutDirection);
    }

    public interface MessageAdapterTouchHelperListener
    {
        void onSwiped(RecyclerView.ViewHolder viewHolder, int direction, int position);
    }
}
