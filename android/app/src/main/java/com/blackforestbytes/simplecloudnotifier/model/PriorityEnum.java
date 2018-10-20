package com.blackforestbytes.simplecloudnotifier.model;

public enum PriorityEnum
{
    LOW(0),
    NORMAL(1),
    HIGH(2);

    public final int ID;

    PriorityEnum(int id) { ID = id; }

    public static PriorityEnum parseAPI(String v) throws Exception
    {
        for (PriorityEnum p : values())
        {
            if (String.valueOf(p.ID).equals(v.trim())) return p;
        }
        throw new Exception("Invalid value for <PriorityEnum> : '"+v+"'");
    }

    public static PriorityEnum parseAPI(int v)
    {
        for (PriorityEnum p : values())
        {
            if (p.ID == v) return p;
        }
        return PriorityEnum.NORMAL;
    }
}
