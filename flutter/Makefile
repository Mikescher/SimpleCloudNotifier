

# Setup
#
# flutter config --jdk-dir "/usr/lib/jvm/default-runtime/bin"
# sudo archlinux-java set java-17-openjdk
#



run:
	flutter pub run build_runner build
	_JAVA_OPTIONS="" flutter run
	
run-android:
	ping -c1 10.10.10.177
	adb connect 10.10.10.177:5555
	flutter pub run build_runner build
	_JAVA_OPTIONS="" flutter run -d 10.10.10.177:5555
	
test:
	dart analyze

fix:
	dart fix --apply

gen:
	flutter pub run build_runner build
	
autoreload:
	@# run `make run` in another terminal (or another variant of flutter run)
	@_utils/autoreload.sh
	
icons:
	flutter pub run flutter_launcher_icons -f "flutter_launcher_icons.yaml"

clean:
	cd android && ./gradlew clean
	flutter clean