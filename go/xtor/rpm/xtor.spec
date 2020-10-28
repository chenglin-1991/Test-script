%global xtor_version %{xtor_package_version}
%global xtor_release %{xtor_package_release}


Summary: XTao xtor Agent and Cli
Name: xtor
Version: %{xtor_version} 
Release: %{xtor_release}
License: GPL
Buildroot: %{_tmppath}/%{name}-buildroot
Group: Applications/File
%define tarball %{name}-%{version}.tar.gz
Source0: %{tarball}


%description
This is distributed scheduler created by XTao to do Alamo OOB service

%package -n %{name}-server
Summary: xtorAgent server
Group: System Environment/Kernel
Provides: %{name}-server = %{version}-%{release}

%package -n %{name}-cli
Summary: xtorCli 
Group: System Environment/Kernel
Provides: %{name}-cli = %{version}-%{release}


%description -n %{name}-server
This is XTao xtor agent server

%description -n %{name}-cli
This is XTao xtor CLI utils

%prep
%setup -q


%build
%configure
make

%postun
%define __debug_install_post   \
         %{_rpmconfigdir}/find-debuginfo.sh %{?_find_debuginfo_opts} "%{_builddir}/%{?buildsubdir}"\
         %{nil}

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/opt/xtor
make DESTDIR=%{buildroot} MANDIR=%{_mandir} BINDIR=%{_sbindir} SYSTEMD_DIR=%{_unitdir} install
cd ../../../../
mkdir -p %{buildroot}/usr/lib/systemd/system
mkdir -p %{buildroot}/etc/xtorsvr
mkdir -p %{buildroot}/var/log/xtorsvr/
cp ./systemd/xtorsvr.service %{buildroot}/usr/lib/systemd/system/xtorsvr.service
cp ./conf/xtorsvr.conf %{buildroot}/etc/xtorsvr/
cd -

%clean
rm -rf %{buildroot}


%files -n %{name}-server
/opt/xtor/xtorsvr
/usr/lib/systemd/system/xtorsvr.service
/etc/xtorsvr
/etc/xtorsvr/xtorsvr.conf
/var/log/xtorsvr/

%files -n %{name}-cli
/usr/bin/xtorcli

%changelog
* Thu Aug 16 2019 Javen Wu <javen.wu@xtaotech.com> - initial version
- create the initial version

