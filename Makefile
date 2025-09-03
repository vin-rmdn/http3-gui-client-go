clean:
	rm -rf out
	mkdir out

build-macos:
	mkdir -p out/HTTP3\ GUI\ Client\ in\ Go.app/Contents/MacOS
	mkdir -p out/HTTP3\ GUI\ Client\ in\ Go.app/Contents/Resources

	go build -o out/HTTP3\ GUI\ Client\ in\ Go.app/Contents/MacOS/http3-gui-client-go main.go

	cp config.yml out/HTTP3\ GUI\ Client\ in\ Go.app/Contents/Resources/
	cp macos/Info.plist out/HTTP3\ GUI\ Client\ in\ Go.app/Contents/
	cp -r macos/http3-gui-client-go.icns out/HTTP3\ GUI\ Client\ in\ Go.app/Contents/Resources
