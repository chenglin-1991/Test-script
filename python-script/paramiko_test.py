#!/usr/bin/python
# -*- coding:UTF-8 -*-
import sys
import os
import paramiko


ip_list = ['10.10.25.21']

def remote_exec(cmd, name='root', passwd='123456', port=22):
    try:
        for host in ip_list:
            ssh = paramiko.SSHClient()
            ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            ssh.connect(host, port, name, passwd)
            str1 = 'cd %s && %s' % (os.getcwd(), cmd)
            print str1
            stdin,stdout,stderr = ssh.exec_command(str1)
            print stdout.read()  
            #print stdin.read()  
            #print stderr.read()  
            ssh.close()
    except Exception, e:
        print e
        return -2
    return 0

def put_file(localfile, remotefile, port=22, name='root', passwd='123456'):
    try:
        if not localfile or not remotefile:
            print error
            return -1
        for host in ip_list:
            con = paramiko.Transport((host, port))
            con.connect(username=name, password=passwd)
            sftp = paramiko.SFTPClient.from_transport(con)
            sftp.put(localfile, remotefile)
            con.close()
    except Exception, e:
        print e
        return -2
    return 0
            
def get_file(remotefile, localfile, port=22, name='root', passwd='123456'):
    try:
        if not localfile or not remotefile:
            print error
            return -1
        for host in ip_list:
            con = paramiko.Transport((host, port))
            con.connect(username=name, password=passwd)
            sftp = paramiko.SFTPClient.from_transport(con)
            if os.path.exists(localfile):
                tmp_localfile = localfile + '_' + host
                sftp.get(remotefile, tmp_localfile)
            else:
                sftp.get(remotefile, localfile)
            con.close()
    except Exception, e:
        print e
        return -2
    return 0

def useage():
    print 'Useage:'
    print '\tpython %s [Options]'%(sys.argv[0]) 
    print 'Options list:'
    print '\tremote_exec'
    print "\t\texample:  python %s 'cat /etc/hosts'"%(sys.argv[0])
    print '\tput'
    print '\t\texample:  python %s localfile  remotefile'%(sys.argv[0])
    print '\tget'
    print '\t\texample:  python %s remotefile  localfile'%(sys.argv[0])

if __name__ == '__main__':

    if len(sys.argv) < 2:
        useage()
        exit(-1)
    if sys.argv[1] == 'remote_exec':
        if not sys.argv[2]:
            print error
            exit(-1)
        remote_exec(sys.argv[2])
    elif sys.argv[1] == 'put':
        if not sys.argv[2] or not sys.argv[3]:
            print error
            exit(-1)
        put_file(sys.argv[2], sys.argv[3])
    elif sys.argv[1] == 'get':
        if not sys.argv[2] or not sys.argv[3]:
            print error
            rexit(-1)
        get_file(sys.argv[2], sys.argv[3])
    exit(-1)
