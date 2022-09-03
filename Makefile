dev: libmdict
	wails dev -assetdir ./frontend/dist -devserver "localhost:3011"

build: libmdict
	wails build

fwatch:
	cd frontend && yarn watch

libmdict:
	@echo "building libmdict"
	rm -rf internal/libmdict/build
	mkdir -p internal/libmdict/build
	cd internal/libmdict/build && cmake .. && make -j4

clean:
	rm -rf internal/libmdict/cmake-build-release
	rm -rf internal/libmdict/cmake-build-debug
	rm -rf internal/libmdict/build
	@echo "make clean done"


.PHONY: watch all libmdict clean