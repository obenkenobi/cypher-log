ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: certs
certs:
	cd $(ROOT_DIR)
	bash ./generate.sh

.PHONY: clean
clean:
	cd $(ROOT_DIR)
	rm *.pem

