dev:
	wails dev --loglevel info
build:
	wails build -devtools

license:
	addlicense -c "Quan Chen <chenquan_act@163.com>" -l gpl3 -v -y 2023 -ignore frontend/**/* -ignore build/**/* -ignore .github/**/* frontend/src

create-dmg:
	create-dmg --volname "Medict" --volicon "build/assets/darwin/dmg_icon.icns" --background "build/assets/darwin/dmg_bg.png" --window-size 512 360 --icon-size 100 --icon "Medict.app" 100 185  --hide-extension "Medict.app" --app-drop-link 388 185 "Medict_3.0.1_Darwin_x86_64.dmg" "build/bin"



.PHONY: build license
