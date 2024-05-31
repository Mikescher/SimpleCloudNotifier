package com.blackforestbytes.simplecloudnotifier

import io.flutter.embedding.android.FlutterActivity

class MainActivity: FlutterActivity() {
    onCreate() {
        GoogleApiAvailability.makeGooglePlayServicesAvailable() 
    }
    onResume() {
        GoogleApiAvailability.makeGooglePlayServicesAvailable() 
    }
}
