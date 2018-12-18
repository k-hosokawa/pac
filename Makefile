INSTALL_SCRIPT_URL = "https://raw.githubusercontent.com/k-hosokawa/pac/master/install"

.PHONY: all
all:
	echo "make test or make test-on-container"

.PHONY: test
test:
	sudo docker build -t k-hosokawa/pac .
	sudo docker run k-hosokawa/pac make test-on-container

.PHONY: test-on-container
test-on-container:
	curl -sL $(INSTALL_SCRIPT_URL) | bash
	/usr/local/bin/pac go install
	/usr/local/bin/pac src install
