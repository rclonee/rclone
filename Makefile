gen-certs:
	openssl genrsa -out key.pem 2048
	openssl req -new -key key.pem -out cert.csr
	openssl x509 -req -days 365 -in cert.csr -signkey key.pem -out cert.pem

https:
	rclone serve http ./SHARE1 --addr :443 --cert  ./cert.pem --key ./key.pem

webdav:
	rclone serve webdav ./SHARE1 --addr :443 --cert ./cert.pem --key ./key.pem

all: gen-certs https
