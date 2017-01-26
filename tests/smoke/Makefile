deps:
	git clone https://github.com/sstephenson/bats.git
	cd bats
	sudo ./install.sh /usr/local

build:
	l0-setup apply --force smoketest 
	eval $(l0-setup endpoint -i)

test:
	bats admin.bats
	bats deploy.bats
	bats environment.bats
	bats job.bats
	bats load_balancer.bats
	bats service.bats
	bats task.bats

destroy:
	l0-setup destroy smoketest

fulltest: build deps test destroy

.PHONY: deps build test destroy fulltest
