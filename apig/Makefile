.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler handler/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --stage dev --aws-profile futuagro --env development --verbose

prod-deploy:
	@make clean build
	sls deploy --stage prod --aws-profile futuagro --env production --verbose

offline:
	@make clean build
	sls offline --stage dev

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
