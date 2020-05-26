VERSION = "v1.0.0"

PACKAGES = $(shell go list ./... | grep -v vendor)

build:
	go build -v -i -o terraform-provider-credstash

install: build
	mkdir -p ~/.terraform.d/plugins/darwin_amd64
	rm -f ~/.terraform.d/plugins/darwin_amd64/terraform-provider-credstash_$(VERSION)
	cp terraform-provider-credstash ~/.terraform.d/plugins/darwin_amd64/terraform-provider-credstash_$(VERSION)

test:
	go test $(TESTOPTS) $(PACKAGES)

testacc:
	AWS_REGION=eu-central-1 AWS_PROFILE=staging CREDSTASH_DYNAMODB_TABLE=credential-store AWS_KMS_KEY=alias/credstash TF_ACC=1 go test -v $(TESTOPTS) $(PACKAGES) -timeout 120m

release:
	GOOS=darwin go build -v -o terraform-provider-credstash_darwin_amd64
	GOOS=linux go build -v -o terraform-provider-credstash_linux_amd64
