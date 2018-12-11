package com.blackforestbytes.simplecloudnotifier.model;

import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;

import com.blackforestbytes.simplecloudnotifier.SCNApp;
import com.blackforestbytes.simplecloudnotifier.lib.collections.CollectionHelper;

import java.util.ArrayList;
import java.util.List;

public class QueryLog
{
    private final static int MAX_HISTORY_SIZE = 192;

    private static QueryLog _instance;
    public static QueryLog instance() { if (_instance == null) synchronized (QueryLog.class) { if (_instance == null) _instance = new QueryLog(); } return _instance; }

    private QueryLog(){ load(); }

    private final List<SingleQuery> history = new ArrayList<>();

    public synchronized void add(SingleQuery r)
    {
        history.add(r);
        while (history.size() > MAX_HISTORY_SIZE) history.remove(0);

        save();
    }

    public synchronized List<SingleQuery> get()
    {
        List<SingleQuery> r = new ArrayList<>(history);
        CollectionHelper.sort_inplace(r, (o1, o2) -> (-1) * o1.Timestamp.compareTo(o2.Timestamp));
        return r;
    }

    public synchronized void save()
    {
        SharedPreferences sharedPref = SCNApp.getContext().getSharedPreferences("QueryLog", Context.MODE_PRIVATE);
        SharedPreferences.Editor e = sharedPref.edit();

        e.clear();

        e.putInt("history_count", history.size());

        for (int i = 0; i < history.size(); i++) history.get(i).save(e, "message["+(i+1000)+"]");

        e.apply();
    }

    public synchronized void load()
    {
        try
        {
            Context c = SCNApp.getContext();
            SharedPreferences sharedPref = c.getSharedPreferences("QueryLog", Context.MODE_PRIVATE);
            int count = sharedPref.getInt("history_count", 0);
            for (int i=0; i < count; i++) history.add(SingleQuery.load(sharedPref, "message["+(i+1000)+"]"));

            CollectionHelper.sort_inplace(history, (o1, o2) -> (-1) * o1.Timestamp.compareTo(o2.Timestamp));
        }
        catch (Exception e)
        {
            Log.e("SC:QL:Load", e.toString());
        }
    }
}
