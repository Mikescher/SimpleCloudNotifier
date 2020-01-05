package com.blackforestbytes.simplecloudnotifier.model;

import android.content.Context;
import android.content.SharedPreferences;

import com.blackforestbytes.simplecloudnotifier.lib.collections.CollectionHelper;
import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.view.MessageAdapter;
import com.blackforestbytes.simplecloudnotifier.SCNApp;

import java.lang.ref.WeakReference;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;

public class CMessageList
{
    private final Object msg_lock = new Object();

    public ArrayList<CMessage> Messages;
    public Set<String> AllAcks;

    private ArrayList<WeakReference<MessageAdapter>> _listener = new ArrayList<>();

    private final static Object _lock = new Object();
    private static CMessageList _inst = null;
    public static CMessageList inst()
    {
        synchronized (_lock)
        {
            if (_inst != null) return _inst;
            return _inst = new CMessageList();
        }
    }

    private CMessageList()
    {
        reloadPrefs();
    }

    public void reloadPrefs()
    {
        synchronized (msg_lock)
        {
            Messages = new ArrayList<>();
            AllAcks  = new HashSet<>();

            SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("CMessageList", Context.MODE_PRIVATE);
            int count = sharedPref.getInt("message_count", 0);
            for (int i=0; i < count; i++)
            {
                long time         = sharedPref.getLong("message["+i+"].timestamp", 0);
                String title      = sharedPref.getString("message["+i+"].title", "");
                String content    = sharedPref.getString("message["+i+"].content", "");
                PriorityEnum prio = PriorityEnum.parseAPI(sharedPref.getInt("message["+i+"].priority", 1));
                long scnid        = sharedPref.getLong("message["+i+"].scnid", 0);

                Messages.add(new CMessage(scnid, time, title, content, prio));
            }

            AllAcks = sharedPref.getStringSet("acks", new HashSet<>());
        }
    }

    public CMessage add(final long scnid, final long time, final String title, final String content, final PriorityEnum pe)
    {
        CMessage msg = new CMessage(scnid, time, title, content, pe);

        boolean run = SCNApp.runOnUiThread(() ->
        {
            SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("CMessageList", Context.MODE_PRIVATE);
            int count = sharedPref.getInt("message_count", 0);

            synchronized (msg_lock)
            {
                Messages.add(msg);
                AllAcks.add(Long.toHexString(msg.SCN_ID));

                while (Messages.size()>SCNSettings.inst().LocalCacheSize) Messages.remove(0);
            }

            if (Messages.size()>1 && Messages.get(Messages.size()-2).Timestamp < msg.Timestamp)
            {
                // quick save

                SharedPreferences.Editor e = sharedPref.edit();

                e.putInt(   "message_count",                count+1);
                e.putLong(  "message["+count+"].timestamp", time);
                e.putString("message["+count+"].title",     title);
                e.putString("message["+count+"].content",   content);
                e.putInt(   "message["+count+"].priority",  pe.ID);
                e.putLong(  "message["+count+"].scnid",     scnid);

                e.putStringSet("acks", AllAcks);

                e.apply();
            }
            else
            {
                // full save

                fullSave(); // does sort in here
            }


            for (WeakReference<MessageAdapter> ref : _listener)
            {
                MessageAdapter a = ref.get();
                if (a == null) continue;
                a.customNotifyDataSetChanged();
                a.scrollToTop();
            }
            CleanUpListener();
        });

        if (!run)
        {
            synchronized (msg_lock)
            {
                Messages.add(new CMessage(scnid, time, title, content, pe));
                AllAcks.add(Long.toHexString(msg.SCN_ID));
            }
            fullSave(); // does sort in here
        }

        return msg;
    }

    public void clear()
    {
        Messages.clear();
        fullSave();

        for (WeakReference<MessageAdapter> ref : _listener)
        {
            MessageAdapter a = ref.get();
            if (a == null) continue;
            a.customNotifyDataSetChanged();
        }
        CleanUpListener();
    }

    public void fullSave()
    {
        synchronized (msg_lock)
        {
            CollectionHelper.sort_inplace(Messages, (a,b) -> Long.compare(a.Timestamp, b.Timestamp));

            SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("CMessageList", Context.MODE_PRIVATE);
            SharedPreferences.Editor e = sharedPref.edit();

            e.clear();

            e.putInt("message_count", Messages.size());

            for (int i = 0; i < Messages.size(); i++)
            {
                e.putLong(  "message["+i+"].timestamp", Messages.get(i).Timestamp);
                e.putString("message["+i+"].title",     Messages.get(i).Title);
                e.putString("message["+i+"].content",   Messages.get(i).Content);
                e.putInt(   "message["+i+"].priority",  Messages.get(i).Priority.ID);
                e.putLong(  "message["+i+"].scnid",     Messages.get(i).SCN_ID);
            }

            e.putStringSet("acks", AllAcks);

            e.apply();
        }
    }

    public CMessage tryGet(int pos)
    {
        synchronized (msg_lock)
        {
            if (pos < 0 || pos >= Messages.size()) return null;
            return Messages.get(pos);
        }
    }

    public CMessage tryGetFromBack(int pos)
    {
        return tryGet(Messages.size() - pos - 1);
    }

    public int size()
    {
        synchronized (msg_lock)
        {
            return Messages.size();
        }
    }

    public void register(MessageAdapter adp)
    {
        _listener.add(new WeakReference<>(adp));
        CleanUpListener();
    }

    private void CleanUpListener()
    {
        for (int i=_listener.size()-1; i >= 0; i--)
        {
            if (_listener.get(i).get() == null) _listener.remove(i);
        }
    }

    public boolean isAck(long id)
    {
        synchronized (msg_lock)
        {
            return AllAcks.contains(Long.toHexString(id));
        }
    }

    public CMessage removeFromBack(int pos)
    {
        CMessage r;
        synchronized (msg_lock)
        {
            int index = Messages.size() - pos - 1;
            r = Messages.remove(index);
        }
        fullSave(); // does sort in here
        return r;
    }

    public void insert(int index, CMessage item)
    {
        synchronized (msg_lock)
        {
            Messages.add(index, item);
        }
        fullSave(); // does sort in here
    }
}
