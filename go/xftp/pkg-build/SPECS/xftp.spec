Name: xftp
Version: 0.0.1
Release: 4dd46d0
License: xtao.com
Summary: Create binary rpm package with xftp

%description
a cluster ftp server

%prep

%build

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/usr/lib/systemd/system
mkdir -p %{buildroot}/etc/xftp
cp /home/goworker/src/xftp/xftp %{buildroot}/usr/local/bin
cp /home/goworker/src/xftp/rpm/xftp.service %{buildroot}/usr/lib/systemd/system
cp /home/goworker/src/xftp/rpm/xftp.conf %{buildroot}/etc/xftp


%files
/usr/local/bin/xftp
/usr/lib/systemd/system/xftp.service
/etc/xftp/xftp.conf


%clean

%changelog

