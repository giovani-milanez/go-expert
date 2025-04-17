Requisitos sqlite3:
Important: because this is a CGO enabled package, you are required to set the environment variable CGO_ENABLED=1 and have a gcc compiler present within your path.

Por favor, instale um compilador GCC. Eu usei este:
https://sourceforge.net/projects/tdm-gcc/files/TDM-GCC%20Installer/tdm64-gcc-5.1.0-2.exe/download

Ative CGO:
go env -w CGO_ENABLED=1

Rodando:
go mod tidy
go run cmd/server/server.go
go run cmd/client/client.go