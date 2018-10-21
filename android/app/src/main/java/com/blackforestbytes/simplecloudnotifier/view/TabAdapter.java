package com.blackforestbytes.simplecloudnotifier.view;

import androidx.fragment.app.Fragment;
import androidx.fragment.app.FragmentManager;
import androidx.fragment.app.FragmentStatePagerAdapter;

public class TabAdapter extends FragmentStatePagerAdapter {

    public NotificationsFragment tab1 = new NotificationsFragment();
    public AccountFragment       tab2 = new AccountFragment();
    public SettingsFragment      tab3 = new SettingsFragment();

    public TabAdapter(FragmentManager fm)
    {
        super(fm);
    }

    @Override
    public Fragment getItem(int position) {

        switch (position) {
            case 0:
                return tab1;
            case 1:
                return tab2;
            case 2:
                return tab3;
            default:
                return null;
        }
    }

    @Override
    public CharSequence getPageTitle(int position)
    {
        switch (position)
        {
            case 0:  return "Notifications";
            case 1:  return "Account";
            case 2:  return "Settings";
            default: return null;
        }
    }

    @Override
    public int getCount() {
        return 3;
    }
}