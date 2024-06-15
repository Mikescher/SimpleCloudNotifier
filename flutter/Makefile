

run:
	flutter pub run build_runner build
	flutter run
	
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