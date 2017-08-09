deps:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen

all: api
api: entity

entity:
	mockgen github.com/quintilesims/layer0/api/entity Provider > ../api/entity/mock_entity/mock_provider.go &
	mockgen github.com/quintilesims/layer0/api/entity Environment > ../api/entity/mock_entity/mock_environment.go &

.PHONY: deps api entity
