SHELL:=/bin/bash

deps:
	git clone https://github.com/sstephenson/bats.git
	sudo ./bats/install.sh /usr/local

test:
	bats admin.bats
	bats deploy.bats
	bats environment.bats
	bats job.bats
	bats load_balancer.bats
	bats service.bats
	bats task.bats

.PHONY: deps test
