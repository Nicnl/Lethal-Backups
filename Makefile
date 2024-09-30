BINARY_NAME=lethal_backups.exe

.PHONY: vue versioninfo go sign

all: all vue versioninfo go sign

vue:
	cd vue && \
	powershell -Command "if (Test-Path 'dist') { Remove-Item 'dist' -Force -Recurse }" && \
	npm run build

versioninfo:
	powershell -Command "if (Test-Path 'resource.syso') { Remove-Item 'resource.syso' -Force }"
	go generate

go:
	powershell -Command "if (Test-Path '${BINARY_NAME}') { Remove-Item '${BINARY_NAME}' -Force }"
	go build -ldflags "-w -s" -o ${BINARY_NAME} .

sign:
	"C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe" sign /f "MyCA.pfx" /fd SHA256 /t http://timestamp.digicert.com \
	${BINARY_NAME}
