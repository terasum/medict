build:
	wails build

license:
	addlicense -c "Quan Chen <chenquan_act@163.com>" -l gpl3 -v -y 2023 -ignore frontend/**/* -ignore build/**/* -ignore .github/**/* frontend/src


.PHONY: build license