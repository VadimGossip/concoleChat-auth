package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/minimock -i UserRepository -o mocks -s "_minimock.go" -g
//go:generate ../../bin/minimock -i UserCacheRepository -o mocks -s "_minimock.go" -g
//go:generate ../../bin/minimock -i AuditRepository -o mocks -s "_minimock.go" -g
