import 'package:flutter/material.dart';
import 'package:toastification/toastification.dart';

class Toaster {
  // https://payamzahedi.com/toastification/

  static const autoCloseDuration = Duration(seconds: 4);
  static const alignment = Alignment.topCenter;
  static const animationDuration = Duration(milliseconds: 200);
  static final borderRadius = BorderRadius.circular(4.0);

  static void simple(String title) {
    toastification.show(
      type: ToastificationType.success,
      style: ToastificationStyle.simple,
      title: Text(title),
      description: Text(title),
      autoCloseDuration: autoCloseDuration,
      borderRadius: borderRadius,
      closeButtonShowType: CloseButtonShowType.none,
      alignment: alignment,
      animationDuration: animationDuration,
      pauseOnHover: false,
      applyBlurEffect: true,
      closeOnClick: true,
      showProgressBar: false,
    );
  }

  static void success(String title, String message) {
    toastification.show(
      type: ToastificationType.success,
      style: ToastificationStyle.flatColored,
      title: Text(title),
      description: Text(message),
      autoCloseDuration: autoCloseDuration,
      borderRadius: borderRadius,
      closeButtonShowType: CloseButtonShowType.none,
      alignment: alignment,
      animationDuration: animationDuration,
      pauseOnHover: false,
      applyBlurEffect: true,
      closeOnClick: true,
      showProgressBar: false,
    );
  }

  static void info(String title, String message) {
    toastification.show(
      type: ToastificationType.info,
      style: ToastificationStyle.flatColored,
      title: Text(title),
      description: Text(message),
      autoCloseDuration: autoCloseDuration,
      borderRadius: borderRadius,
      closeButtonShowType: CloseButtonShowType.none,
      alignment: alignment,
      animationDuration: animationDuration,
      pauseOnHover: false,
      applyBlurEffect: true,
      closeOnClick: true,
      showProgressBar: false,
    );
  }

  static void warn(String title, String message) {
    toastification.show(
      type: ToastificationType.warning,
      style: ToastificationStyle.flatColored,
      title: Text(title),
      description: Text(message),
      autoCloseDuration: autoCloseDuration,
      borderRadius: borderRadius,
      closeButtonShowType: CloseButtonShowType.none,
      alignment: alignment,
      animationDuration: animationDuration,
      pauseOnHover: false,
      applyBlurEffect: true,
      closeOnClick: true,
      showProgressBar: false,
    );
  }

  static void error(String title, String message) {
    toastification.show(
      type: ToastificationType.error,
      style: ToastificationStyle.flatColored,
      title: Text(title),
      description: Text(message),
      autoCloseDuration: autoCloseDuration,
      borderRadius: borderRadius,
      closeButtonShowType: CloseButtonShowType.none,
      alignment: alignment,
      animationDuration: animationDuration,
      pauseOnHover: false,
      applyBlurEffect: true,
      closeOnClick: true,
      showProgressBar: false,
    );
  }
}
