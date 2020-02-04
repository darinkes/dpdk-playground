SUBDIRS= sniffer

all: nff-go $(SUBDIRS)

nff-go:
	make -C libbpf/src all
	sudo make -C libbpf/src install
	sudo sh -c "echo /usr/lib64 > /etc/ld.so.conf.d/usrlib64.conf"
	sudo ldconfig
	cd nff-go && go mod download && make -j$$(nproc)

$(SUBDIRS):
	$(MAKE) --print-directory -C $@

.PHONY: all nff-go $(SUBDIRS)
