package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/minimock -i UserService -o mocks -s "_minimock.go" -g
//go:generate ../../bin/minimock -i UserCacheService -o mocks -s "_minimock.go" -g
//go:generate ../../bin/minimock -i AuditService -o mocks -s "_minimock.go" -g
//go:generate ../../bin/minimock -i UserConsumerService -o mocks -s "_minimock.go" -g
//go:generate ../../bin/minimock -i UserProducerService -o mocks -s "_minimock.go" -g
