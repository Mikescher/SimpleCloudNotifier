package com.blackforestbytes.simplecloudnotifier.model;

import android.content.Context;
import android.content.SharedPreferences;

import com.blackforestbytes.simplecloudnotifier.lib.string.Str;
import com.blackforestbytes.simplecloudnotifier.view.MessageAdapter;
import com.blackforestbytes.simplecloudnotifier.SCNApp;

import java.lang.ref.WeakReference;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.Set;

public class CMessageList
{
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

    public CMessage add(final long scnid, final long time, final String title, final String content, final PriorityEnum pe)
    {
        CMessage msg = new CMessage(scnid, time, title, content, pe);

        boolean run = SCNApp.runOnUiThread(() ->
        {
            SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("CMessageList", Context.MODE_PRIVATE);
            int count = sharedPref.getInt("message_count", 0);

            SharedPreferences.Editor e = sharedPref.edit();

            Messages.add(msg);
            AllAcks.add(Long.toHexString(msg.SCN_ID));

            while (Messages.size()>SCNSettings.inst().LocalCacheSize) Messages.remove(0);

            e.putInt(   "message_count",                count+1);
            e.putLong(  "message["+count+"].timestamp", time);
            e.putString("message["+count+"].title",     title);
            e.putString("message["+count+"].content",   content);
            e.putInt(   "message["+count+"].priority",  pe.ID);
            e.putLong(  "message["+count+"].scnid",     scnid);

            e.putStringSet("acks", AllAcks);

            e.apply();

            for (WeakReference<MessageAdapter> ref : _listener)
            {
                MessageAdapter a = ref.get();
                if (a == null) continue;
                a.customNotifyItemInserted(count);
                a.scrollToTop();
            }
            CleanUpListener();
        });

        if (!run)
        {
            Messages.add(new CMessage(scnid, time, title, content, pe));
            AllAcks.add(Long.toHexString(msg.SCN_ID));
            fullSave();
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

    public CMessage tryGet(int pos)
    {
        if (pos < 0 || pos >= Messages.size()) return null;
        return Messages.get(pos);
    }

    public CMessage tryGetFromBack(int pos)
    {
        return tryGet(Messages.size() - pos - 1);
    }

    public int size()
    {
        return Messages.size();
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
        return AllAcks.contains(Long.toHexString(id));
    }

    public void remove(int index)
    {
        Messages.remove(index);
        fullSave();
    }

    public void insert(int index, CMessage item)
    {
        Messages.add(index, item);
        fullSave();
    }
}
