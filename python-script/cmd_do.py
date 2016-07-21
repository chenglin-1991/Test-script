#!/usr/bin/python
#coding:UTF-8
import os
import sys
import threading
import pexpect
import multiprocessing
import time

class glfs:
    __username = 'root'
    __passwd = '123456'
    __exe_install = "./autogen.sh && ./configure --build=x86_64-unknown-linux-gnu " \
                    "--host=x86_64-unknown-linux-gnu --target=x86_64-redhat-linux-gnu " \
                    "--program-prefix=--prefix=/usr --exec-prefix=/usr --bindir=/usr/bin " \
                    "--sbindir=/usr/sbin --sysconfdir=/etc --datadir=/usr/share"\
                    "--includedir=/usr/include --libdir=/usr/lib64 --libexecdir=/usr/libexec"\
                    "--localstatedir=/var --sharedstatedir=/var/lib --mandir=/usr/share/man " \
                    "--infodir=/usr/share/info --enable-systemtap=no " \
                    "--disable-experimental --disable-glupy && make -j4"

    def hostname_to_list (self, hostname_str):
        self.hostname_list = hostname_str.split(' ')

    def glfs_make_mult_process (self):
        for hostname in self.hostname_list:
            #p = multiprocessing.Process(target=self.__do_install,args=(hostname,))
            #p.start()
            self.__do_install (hostname)

    def __do_install(self,hostname):
        passwd = r'.+password:'
        yes = r'.+yes'
        login = r'login:'
        cmd_str = 'cd %s && %s' % (os.getcwd(),self.__exe_install)
        ssh_str = 'ssh %s' % hostname
        child = pexpect.spawn(ssh_str)
        child.logfile_send = sys.stdout
        index = child.expect([yes,passwd,login])
        print index
        if index == 0:
            child.sendline('yes')
            child.expect(passwd)
            child.sendline('123456')
            child.expect(login)
            child.sendline(cmd_str)
            while True:
                index = child.expect(['.*',pexpect.TIMEOUT])
                if index == 0:
                    break
                elif index == 1:
                    continue
        elif index == 1:
            child.sendline('123456')
            child.expect(login)
            child.sendline(cmd_str)
            while True:
                index = child.expect([pexpect.EOF,pexpect.TIMEOUT])
                print index
                if index == 0:
                    break
                elif index == 1:
                    continue
        elif index == 2:
            child.sendline(cmd_str)
            while True:
                index = child.expect(['.*',pexpect.TIMEOUT])
                if index == 0:
                    break
                elif index == 1:
                    continue

    def install_glfs (self):
        for hostname in self.hostname_list:
            #self.glfs_make_mult_process()
            self.__do_install(hostname)

def main():
    obglfs = glfs()
    while 1:
        print '''本脚本支持:
            1、代码的编译安装
            2、代码的快速编译安装
            3、在不同的机器上执行相同的命令，并返回结果
            '''
        #num = raw_input('请输入你的需求编号：')
        #num = int (num)
        num = 1

        if num == 1:
            #hostname_str = raw_input('请输入主机名/IP:')
            hostname_str = '10.10.110.185'
            obglfs.hostname_to_list (hostname_str)
            obglfs.install_glfs ()
            break
        elif num == 2:
            break
        elif num == 3:
            break;
        else:
            print '输入出错！'

if __name__ == '__main__':
    main()
