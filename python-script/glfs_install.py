#!/usr/bin/python
#coding:UTF-8
import os
import sys
import threading
import multiprocessing
import time

path = ''

def install_glfs (hostname):
    print '==========%s=================' % hostname
    str = "ssh %s 'cd %s &&./autogen.sh && ./configure && make -j4 && make install' " % (hostname,path)
    os.system (str)

def f_install_glfs (hostname):
    print '==========%s=================' % hostname
    str = "ssh %s 'cd %s && make -j4 && make install' " % (hostname,path)
    os.system (str)

def exe_shell (hostname):
    print '==========%s=================' % hostname
    str = "ssh %s 'cd %s && %s'" % (hostname,path,sys.argv[-1])
    os.system (str)

if __name__ == '__main__':
    if (len (sys.argv) < 3):
        print "Useage ./install_glusterfs install/f-install/exe hostname/IP ['shell cmd string']"
        print 'install 完整的glusterfs安装'
        print 'f-install 简单的make和make install'
        print 'exe 执行一些shell命令'
        exit(-1)

    p = os.popen ('pwd').read()
    path = p[:-1]
    install = sys.argv[1]

    process = []
    arry = sys.argv[2:]
    print arry

    '''
    if install == 'install':
        for i in arry:
            t = threading.Thread(target=install_glfs,args=(i,))
            threads.append(t)
            #install_glfs (i)
    elif install == 'f-install':
        for i in arry:
            t = threading.Thread(target=f_install_glfs,args=(i,))
            threads.append(t)
            #f_install_glfs (i)
    elif install == 'exe':
        for i in arry:
            if i == sys.argv[-1]:
                break
            t = threading.Thread(target=exe_shell,args=(i,))
            threads.append(t)
            #exe_shell (i)
    print len(threads)
    for t in threads:
        t.start()

    for t in threads:
        t.join()
    print '-------end-------'
    '''
    if install == 'install':
        for i in arry:
            t = multiprocessing.Process(target=install_glfs,args=(i,))
            process.append(t)
            t.start()
            #install_glfs (i)
    elif install == 'f-install':
        for i in arry:
            t = multiprocessing.Process(target=f_install_glfs,args=(i,))
            process.append(t)
            t.start()
            #f_install_glfs (i)
    elif install == 'exe':
        for i in arry:
            if i == sys.argv[-1]:
                break
            t = threading.Thread(target=exe_shell,args=(i,))
            process.append(t)
            t.start()
            #exe_shell (i)
