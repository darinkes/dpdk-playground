include ../nff-go/mk/include.mk

all:
	go build -tags '$(GO_BUILD_TAGS)' -v .
