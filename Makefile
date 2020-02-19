GOPATH	= $(CURDIR)
BINDIR	= $(CURDIR)/bin

PROGRAMS = vault-keepass

depend:
	env GOPATH=$(GOPATH) go get -u github.com/atotto/clipboard
	env GOPATH=$(GOPATH) go get -u github.com/sirupsen/logrus
	env GOPATH=$(GOPATH) go get -u golang.org/x/crypto/ssh/terminal

build:
	env GOPATH=$(GOPATH) go install $(PROGRAMS)

destdirs:
	mkdir -p -m 0755 $(DESTDIR)/usr/bin

strip: build
	strip --strip-all $(BINDIR)/vault-keepass

install: strip destdirs install-bin

install-bin:
	install -m 0755 $(BINDIR)/vault-keepass $(DESTDIR)/usr/bin

clean:
	/bin/rm -f bin/vault-keepass

distclean: clean
	rm -rf src/github.com/
	rm -rf src/golang.org/

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: depend build strip install

