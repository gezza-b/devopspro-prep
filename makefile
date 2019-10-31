
.PHONY: build clean deploy package


build:
	dep ensure -v
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/imghandler lambda/main.go
	chmod 755 bin/imghandler
clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	#sls deploy --verbose
	# Lambda deployment package
	$(ZIPFILE): clean lambda
		zip -9 -r $(ZIPFILE) $(OUTPUT)

package:
	aws cloudformation package --template-file $(TEMPLATE) --s3-bucket $(S3_BUCKET) --output-template-file $(PACKAGED_TEMPLATE)
# stat -f %A 