.PHONY: certs

CERT_DIR ?= certs
CERT_KEY  = $(CERT_DIR)/tls.key
CERT_CRT  = $(CERT_DIR)/tls.crt

certs: $(CERT_KEY) $(CERT_CRT)

$(CERT_KEY) $(CERT_CRT):
	mkdir -p $(CERT_DIR)
	openssl req -x509 -newkey rsa:4096 -nodes \
		-keyout $(CERT_KEY) \
		-out $(CERT_CRT) \
		-days 365 \
		-subj "/CN=antivax"
