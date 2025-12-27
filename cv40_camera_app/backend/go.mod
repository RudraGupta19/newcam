module cv40-camera-backend

go 1.22.5

toolchain go1.24.6

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.0
	lt/client/go v0.0.0-00010101000000-000000000000
)

require (
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	golang.org/x/sys v0.22.0 // indirect
)

// Use the local Enciris LT Go SDK module located at repo-root/go
replace lt/client/go => ../../go
