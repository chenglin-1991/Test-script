#!/usr/bin/python
import os
import sys
import re
import pexpect

if __name__ == '__main__':
    str1 = r'.+password:'
    child = pexpect.spawn('ssh -p 22 10.10.110.185')
    child.logfile = sys.stdout
    index = child.expect(str1)
    if index == 0:
        child.sendline('123456')
        index = child.expect(r'login:')
        print index
        #index = child.sendline('touch /root/fanjiaolong.txt')
        index = child.sendline('mkdir /root/mkmkmkm')
        index = child.expect('.*')
        print index
        child.sendline('ls /')
        child.expect('usr')

