bldNum = $(if $(BLD_NUM),$(BLD_NUM),999)
version = $(if $(VERSION),$(VERSION),1.0.0)
productVersion = $(version)-$(bldNum)
ARTIFACTS = build/artifacts/

# This allows the container tags to be explicitly set.
DOCKER_USER = couchbase
DOCKER_TAG = v1

# What exact revision is this?
GIT_REVISION := $(shell git rev-parse HEAD)

# Set this to, for example beta1, for a beta release.
# This will affect the "-v" version strings and docker images.
# This is analogous to revisions in DEB and RPM archives.
revision = $(if $(REVISION),$(REVISION),)

# These are propagated into each binary so we can tell for sure the exact build
# that a binary came from.
# LDFLAGS = "-s -w -X github.com/couchbaselabs/workbench-prototype/pkg/version.version=$(version) -X github.com/couchbaselabs/workbench-prototype/pkg/version.revision=$(revision) -X github.com/couchbaselabs/workbench-prototype/pkg/version.buildNumber=$(bldNum) -X github.com/couchbaselabs/workbench-prototype/pkg/version.gitRevision=$(GIT_REVISION)"

# Hardcode version values for testing
# TEST_LDFLAGS = "-X github.com/couchbaselabs/workbench-prototype/pkg/version.version=1 -X github.com/couchbaselabs/workbench-prototype/pkg/version.revision=2 -X github.com/couchbaselabs/workbench-prototype/pkg/version.buildNumber=3 -X github.com/couchbaselabs/workbench-prototype/pkg/version.gitRevision=456"

.PHONY: all agent-build agent-cross agent-dist build fluent-bit-build fluent-bit-cross lint container container-public container-lint container-scan dist fluent-bit-cross test-dist container-clean clean

all: clean build lint agent-build container container-lint container-scan dist test-dist

build: container

fluent-bit-build:
	tools/build-fluent-bit.sh

agent-build: fluent-bit-build
	CGO_ENABLED=0 go build -o ./build -tags netgo ./agent/cmd/cbhealthagent


# TODO (CMOS-302): support other distributions
build/fluent-bit-linux-amd64:
	$(eval OS := linux)
	$(eval ARCH := amd64)
	$(eval export OS ARCH)
	./tools/build-fluent-bit-docker.sh

fluent-bit-cross: build/fluent-bit-linux-amd64

agent-cross: OS ?= linux
agent-cross: ARCH ?= amd64
agent-cross: build/fluent-bit-$(OS)-$(ARCH)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -o ./build/cbhealthagent-$(OS)-$(ARCH) -tags netgo ./agent/cmd/cbhealthagent

agent-dist:
	-rm build/cbhealthagent.zip
	-rm -r build/etc
	tools/prepare-fluent-bit-config.sh
	cd build; zip -r cbhealthagent.zip * -x cbmultimanager

image-artifacts:
	mkdir -p $(ARTIFACTS)
	cp -rv cmd server ui $(ARTIFACTS)
	cp -v Dockerfile* go.* .*ignore *.sh LICENSE $(ARTIFACTS)

# This target (and only this target) is invoked by the production build job.
# This job will archive all files that end up in the dist/ directory.
dist: image-artifacts
	rm -rf dist
	mkdir -p dist
	tar -C $(ARTIFACTS) -czvf dist/couchbase-cluster-manager-image_$(productVersion).tgz .
	rm -rf $(ARTIFACTS)

lint: container-lint
	docker run --rm -i -v ${PWD}:/work -w /work golangci/golangci-lint:v1.42.1 golangci-lint run -v
	tools/shellcheck.sh
	tools/licence-lint.sh
	go run ./tools/validate-checker-docs.go

# NOTE: This target is only for local development. While we use this Dockerfile
# (for now), the actual "docker build" command is located in the Jenkins job.
container:
# Set the same variables the official builds set (https://github.com/couchbase/build-tools/blob/7a44c105cf8768a7f758e80968b357eb37c08fc0/k8s-microservice/jenkins/util/build-k8s-images.sh#L82-L86)
	docker build -f Dockerfile --build-arg PROD_VERSION=$(version) --build-arg PROD_BUILD=$(bldNum) -t ${DOCKER_USER}/cluster-manager:${DOCKER_TAG} .

container-lint:
	docker run --rm -i hadolint/hadolint < Dockerfile

# We are using a fuller fat base image than probably necessary so use the Dive checks as information only.
container-scan: container
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy \
		--severity "HIGH,CRITICAL" --ignore-unfixed --exit-code 1 --no-progress ${DOCKER_USER}/cluster-manager:${DOCKER_TAG}
	-docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -e CI=true wagoodman/dive \
		${DOCKER_USER}/cluster-manager:${DOCKER_TAG}

# This target pushes the containers to a public repository.
# A typical one liner to deploy to the cloud would be:
# 	make container-public -e DOCKER_USER=couchbase DOCKER_TAG=2.0.0
container-public: container
	docker push ${DOCKER_USER}/cluster-manager:${DOCKER_TAG}

# Special target to verify the internal release pipeline will work as well
# Take the archive we would make and extract it to a local directory to then run the docker builds on
test-dist: dist
	rm -rf test-dist/
	mkdir -p test-dist/
	tar -xzvf dist/couchbase-cluster-manager-image_$(productVersion).tgz -C test-dist/
	docker build -f test-dist/Dockerfile --build-arg PROD_VERSION=$(version) --build-arg PROD_BUILD=$(bldNum) test-dist/ -t ${DOCKER_USER}/cluster-manager-test-dist:${DOCKER_TAG}

# Remove our images then remove dangling ones to prevent any caching
container-clean:
	docker rmi -f ${DOCKER_USER}/cluster-manager:${DOCKER_TAG} \
				  ${DOCKER_USER}/cluster-manager-test-dist:${DOCKER_TAG}
	docker image prune --force

clean: container-clean
	rm -rf $(ARTIFACTS) dist/ test-dist/
