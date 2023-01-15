package com.blackforestbytes;

import de.ralleytn.simple.json.JSONArray;
import de.ralleytn.simple.json.JSONFormatter;
import de.ralleytn.simple.json.JSONObject;

import java.io.ObjectInputStream;
import java.net.URI;
import java.nio.file.FileSystems;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

public class Main {
    @SuppressWarnings("unchecked")
    public static void main(String[] args) {
        if (args.length != 1) {
            System.err.println("call with ./androidExportConvert scn_export.dat");
            return;
        }

        try {

            var path = FileSystems.getDefault().getPath(args[0]).normalize().toAbsolutePath().toUri().toURL();

            ObjectInputStream stream = new ObjectInputStream(path.openStream());

            Map<String, ?> d1 = new HashMap<>((Map<String, ?>)stream.readObject());
            Map<String, ?> d2 = new HashMap<>((Map<String, ?>)stream.readObject());
            Map<String, ?> d3 = new HashMap<>((Map<String, ?>)stream.readObject());
            Map<String, ?> d4 = new HashMap<>((Map<String, ?>)stream.readObject());

            stream.close();

            JSONObject root = new JSONObject();

            var subConfig = new JSONObject();
            var subIAB = new JSONArray();
            var subCMessageList = new JSONArray();
            var subAcks = new JSONArray();
            var subQueryLog = new JSONArray();

            for (Map.Entry<String, ?> entry : d1.entrySet())
            {
                if (entry.getValue() instanceof String)  subConfig.put(entry.getKey(),    (String)entry.getValue());
                if (entry.getValue() instanceof Boolean) subConfig.put(entry.getKey(),   (Boolean)entry.getValue());
                if (entry.getValue() instanceof Float)   subConfig.put(entry.getKey(),     (Float)entry.getValue());
                if (entry.getValue() instanceof Integer) subConfig.put(entry.getKey(),       (Integer)entry.getValue());
                if (entry.getValue() instanceof Long)    subConfig.put(entry.getKey(),      (Long)entry.getValue());
                if (entry.getValue() instanceof Set<?>)  subConfig.put(entry.getKey(), ((Set<String>)entry.getValue()).toArray());
            }

            for (int i = 0; i < (Integer)d2.get("c"); i++) {
                var obj = new JSONObject();
                obj.put("key", d2.get("["+i+"]->key"));
                obj.put("value", d2.get("["+i+"]->value"));
                subIAB.add(obj);
            }

            for (int i = 0; i < (Integer)d3.get("message_count"); i++) {
                if (d3.get("message["+i+"].scnid") == null)
                    throw new Exception("ONF");

                var obj = new JSONObject();
                obj.put("timestamp", d3.get("message["+i+"].timestamp"));
                obj.put("title", d3.get("message["+i+"].title"));
                obj.put("content", d3.get("message["+i+"].content"));
                obj.put("priority", d3.get("message["+i+"].priority"));
                obj.put("scnid", d3.get("message["+i+"].scnid"));
                subCMessageList.add(obj);
            }

            subAcks.addAll(((Set<String>)d3.get("acks")).stream().map(p -> Long.decode("0x"+p)).toList());

            for (int i = 0; i < (Integer)d4.get("history_count"); i++) {
                if (d4.get("message["+(i+1000)+"].Name") == null)
                    throw new Exception("ONF");

                var obj = new JSONObject();
                obj.put("Level", d4.get("message["+(i+1000)+"].Level"));
                obj.put("Timestamp", d4.get("message["+(i+1000)+"].Timestamp"));
                obj.put("Name", d4.get("message["+(i+1000)+"].Name"));
                obj.put("URL", d4.get("message["+(i+1000)+"].URL"));
                obj.put("Response", d4.get("message["+(i+1000)+"].Response"));
                obj.put("ResponseCode", d4.get("message["+(i+1000)+"].ResponseCode"));
                obj.put("ExceptionString", d4.get("message["+(i+1000)+"].ExceptionString"));
                subQueryLog.add(obj);
            }

            root.put("config", subConfig);
            root.put("iab", subIAB);
            root.put("cmessagelist", subCMessageList);
            root.put("acks", subAcks);
            root.put("querylog", subQueryLog);

            System.out.println(new JSONFormatter().format(root.toString()));

        } catch (Exception e) {
            e.printStackTrace();
        }

    }
}