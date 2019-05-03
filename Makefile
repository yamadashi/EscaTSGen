DBNAME:=escaTSGen
ENV:=development

deps:
	which dep || go get -v -u github.com/golang/dep/cmd/dep
	dep ensure 
	go get github.com/rubenv/sql-migrate/...

test:
	go test -v ./...

run:
	go run base/base.go

migrate/init:
	mysql -u root -h localhost --protocol tcp -e "create database \`$(DBNAME)\`" -p

migrate/up:
	sql-migrate up -env=$(ENV)

migrate/down:
	sql-migrate down -env=$(ENV)

migrate/status:
	sql-migrate status -env=$(ENV)