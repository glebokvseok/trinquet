openssl genpkey -algorithm RSA -out ca.key -aes256
openssl req -key ca.key -new -x509 -out ca.pem -days 3650 -subj "/CN=test-stand.online"
openssl genpkey -algorithm RSA -out neo4j.key
openssl req -new -key neo4j.key -out neo4j.csr -subj "/CN=neo4j.test-stand.online"
openssl x509 -req -in neo4j.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out neo4j.crt -days 3650
