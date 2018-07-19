# Secure Tranport Protocol

Package stp implements the client and server communication protocol
that is meant to be simple to use while still being secure.


# Key/Cert Generation

These are the commands used to build the test keys and certs.

	openssl genrsa -des3 -out ca.key 4096
	openssl req -x509 -new -nodes -key ca.key -sha256 -days 10240 -subj "/C=US/ST=VA/L=Langley/O=Internal Affairs/CN=localhost" -out ca.crt
	openssl genrsa -out client.key 2048
	openssl req -new -key client.key  -subj "/C=US/ST=VA/L=Langley/O=Internal Affairs/CN=localhost" -out client.csr
	openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 10240 -sha256

