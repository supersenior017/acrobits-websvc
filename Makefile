.PHONY: docker

NAME           := acrobits-websvc

prefix         ?= /usr/local
exec_prefix    ?= $(prefix)
bindir         ?= $(exec_prefix)/bin
sysconfdir     ?= $(prefix)/etc
srcdir         ?= ./src

confname       := acrobits-websvc.json
targetdir      := ./target
target         := $(targetdir)/$(NAME)
bindestdir     := $(DESTDIR)$(bindir)
sysconfdestdir := $(DESTDIR)$(sysconfdir)

all: build

build:
	go build -o $(target) $(srcdir)/

installdirs:
	install -d $(bindestdir)/ $(sysconfdestdir)/

install: installdirs
	install $(target) $(bindestdir)/
	install -b -m 0644 ./config/$(confname) $(sysconfdestdir)/

uninstall:
	rm -f $(bindestdir)/$(NAME)
	rm -f $(sysconfdestdir)/$(confname)

clean:
	rm -rf $(targetdir)/

docker:
	docker build -t $(NAME) -f ./docker/Dockerfile "."
