#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import os
import time
import subprocess

def main():
    if (len(sys.argv) < 2):
        print "python force-set-dirty.py fuse_mount_dir_path scan_dirtry_log_file_path"
        print "see wiki for help:"
        print '                  http://192.168.12.50:3000/alamo/glusterfs-3.7.11/wiki/'
        return None

    dirty_dir_path = sys.argv[1]
    if sys.argv[1][-1] == '/':
        dirty_dir_path = sys.argv[1][:-1]

    log = open(sys.argv[2],'a+')

    file = open("dirty_list")
    for line in file:
        line=line.strip('\n')
        tmp_dirty_dir_path = "%s/%s" %(dirty_dir_path,line)
        for root,dir_path,file_path in os.walk(tmp_dirty_dir_path, topdown=False):
            for dir_name in dir_path:
                tmp_path = os.path.join(root,dir_name)
                cmd = "setfattr -n trusted.glusterfs.quota.dirty -v 0x3100 %s" %(tmp_path)
                tmp_path_string = tmp_path + '\n'

                info = subprocess.Popen(cmd, shell=True,
                       stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                info.wait()
                stderr1 = info.stderr.readlines()
                if stderr1:
                    print stderr1
                    continue
                else:
                    log.write(tmp_path_string)

                cmd = "stat %s" %(tmp_path)
                info = subprocess.Popen(cmd, shell=True,
                                        stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                info.wait()
                stderr1 = info.stderr.readlines()
                if stderr1:
                    print stderr1
    file.close()

if __name__ == '__main__':
    main()