release:
	mkdir -p build
	mkdir -p build/db
	docker build -t srtl:latest .
	docker create -ti --name dummy srtl:latest bash
	docker cp dummy:/app/srtl ./build
	docker rm -f dummy

	cp -R templates/ build/
	cp -R static/ build/
	cp install.sh build/
	tar -czvf release.tar.gz ./build

build:
	go build -tags builtin_static -o srtl ./cmd/http

run:
	HOST=localhost PORT=8080 go run -tags builtin_static ./cmd/http

clean:
	rm -rv db/
	rm -rv build/
	rm -rv release.tar.gz