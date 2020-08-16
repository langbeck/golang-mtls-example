CREDS_DIR=$(PWD)/creds

define define-ca
$(CREDS_DIR)/ca-$1-cert.pem: CN=$2
$(CREDS_DIR)/ca-$1-cert.pem:

certificates+=$(CREDS_DIR)/ca-$1-cert.pem
endef

define define-cert
$(CREDS_DIR)/$1-$2-cert.pem: CN=$3
$(CREDS_DIR)/$1-$2-cert.pem: CA_CERT=$(CREDS_DIR)/ca-$2-cert.pem
$(CREDS_DIR)/$1-$2-cert.pem: CA_KEY=$(CREDS_DIR)/ca-$2-key.pem

certificates+=$(CREDS_DIR)/$1-$2-cert.pem
endef

.PHONY: gen-certs
gen-certs: .gen-certs


$(eval $(call define-ca,server,Server CA))
$(eval $(call define-ca,client,Client CA))

$(eval $(call define-cert,tls,server,localhost))
$(eval $(call define-cert,tls,client,Client TLS))


.PHONY: .gen-certs
.gen-certs: $(certificates)


################################
## Generic generation targets ##
################################

ca-%-cert.pem: ca-%-key.pem
	openssl req -x509 -key "$(<)" -out "$(@)" -subj "/CN=$(CN)"


.SECONDEXPANSION:
%-cert.pem: %-csr.pem $$(CA_CERT) $$(CA_KEY)
	@mkdir -p $(@D)
	openssl x509 -req -in "$(<)" -out "$(@)" -CA "$(CA_CERT)" -CAkey "$(CA_KEY)" -CAcreateserial

%-csr.pem: %-key.pem
	@mkdir -p $(@D)
	openssl req -new -key "$(<)" -out "$(@)" -subj "/CN=$(CN)"

.SECONDARY:
%-key.pem:
	@mkdir -p $(@D)
	openssl genrsa -out $(@)
