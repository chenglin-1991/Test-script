XTOR_VERSION=${VERSION}
XTOR_RELEASE=$(shell git rev-parse HEAD | cut -n -b -6)

RPM_BUILD_DIR=$(top_srcdir)/rpmbuild

SUBDIRS = xtorsvr xtorcli

AUTOMAKE_OPTIONS = foreign
EXTRA_DIST  = autogen.sh META

distclean-local::
	-$(RM) -R autom4te*.cache
	-$(RM) -R *.tar.gz

clean:
	$(RM) -R *.tar.gz
	$(RM) -rf ./rpmbuild/*

dist-hook:
	sed -i 's/Release:[[:print:]]*/Release:      $(RELEASE)/' \
		$(distdir)/META

ctags:
	$(RM) $(top_srcdir)/tags
	find $(top_srcdir) -name .git -prune -o -name '*.[hc]' | xargs ctags

etags:
	$(RM) $(top_srcdir)/TAGS
	find $(top_srcdir) -name .pc -prune -o -name '*.[hc]' | xargs etags -a

tags: ctags etags

pkg: @DEFAULT_PACKAGE@

rpm:    dist
	mkdir -p ./rpmbuild/SOURCES
	rm -rf ./rpmbuild/SOURCES/*
	cp -f *.tar.gz ./rpmbuild/SOURCES/
	echo "%_topdir `pwd`/rpmbuild">/`whoami`/.rpmmacros
	rpmbuild -ba --define "xtor_package_version $(XTOR_VERSION)" --define "xtor_package_release ${XTOR_RELEASE}" xtor.spec
	cp -r ./rpmbuild/RPMS/x86_64/* ./
