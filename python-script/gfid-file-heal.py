#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import os
import time
import subprocess

def main():
    if (len(sys.argv) < 6):
        print "python gfid.py target_dir brick_path log_file_path stat_file_path stat_dir_path err_log"
        return None

    brick_path = sys.argv[2]
    if sys.argv[2][-1] == '/':
        brick_path = sys.argv[2][:-1]

    file_log = None
    file_log = open(sys.argv[3],'a')
    if file_log == None:
        print 'open log_file_path %s failed!' %(sys.argv[3])

    stat_file = open(sys.argv[4],'a')
    if stat_file == None:
        print 'open log_file_path %s failed!' %(sys.argv[4])

    stat_dir = open(sys.argv[5],'w')
    if stat_dir == None:
        print 'open log_file_path %s failed!' %(sys.argv[5])

    err_log = open(sys.argv[6],'a')
    if err_log == None:
        print 'open log_file_path %s failed!' %(sys.argv[6])


    file_log.write('-------------dir gfid link file check------------\n')

    i = 0
    parent_item = None
    for root,dir_path,file_path in os.walk(sys.argv[1]):
        if '.glusterfs' in dir_path:
            index = dir_path.index('.glusterfs')
            del dir_path[index]

        for dir_name in dir_path:
            print ("root={}".format(root))
            tmp_path = os.path.join(root,dir_name)
            index= tmp_path.rfind('/')
            parent_path = tmp_path[:index]

            cmd = "getfattr -dm . -e hex --no-dereference %s | grep gfid= | cut -d '=' -f 2" %(parent_path)
            info = subprocess.Popen(cmd, shell=True,
                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr1 = info.stderr.readlines()
            if not (stderr1 and len(stderr1) == 1 and "Removing leading '/' from absolute path names" in stderr1[0]):
                if "No such file or directory" in stderr1[0]:
                    continue
                err_log.write("{}\n".format(' '.join(stderr1)))
                print stderr1
                print cmd
                return None
            print '---------------start------------------'
            stdout1 = info.stdout.readlines()
            for item in stdout1:
                parent_item = item[:-1]
                result = "%s/.glusterfs/%s/%s/%s-%s-%s-%s-%s" \
                        %(brick_path,parent_item[2:4],parent_item[4:6],parent_item[2:10],\
                        parent_item[10:14],parent_item[14:18],parent_item[18:22],parent_item[22:])
                parent_str =  'parent_path---> %s parent_gfid---> %s' %(parent_path,result)
                print parent_str

            cmd = "getfattr -dm . -e hex --no-dereference %s | grep gfid= | cut -d '=' -f 2" %(tmp_path)
            #cmd = cmd.encode ('utf8')
            print cmd
            info = subprocess.Popen(cmd, shell=True,
                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr1 = info.stderr.readlines()
            if not (stderr1 and len(stderr1) == 1 and "Removing leading '/' from absolute path names" in stderr1[0]):
                if "No such file or directory" in stderr1[0]:
                    continue
                err_log.write("{}\n".format(' '.join(stderr1)))
                print stderr1
                print cmd
                return None

            stdo = info.stdout.readlines()
            for item in stdo:
                item = item[:-1]
                result = "%s/.glusterfs/%s/%s/%s-%s-%s-%s-%s" %(brick_path,item[2:4],item[4:6],item[2:10],item[10:14],item[14:18],item[18:22],item[22:])
                file_str =  "file_path---> %s file_gfid---> %s" %(tmp_path,result)
                print file_str
                cmd = "ls -l %s/.glusterfs/%s/%s/%s-%s-%s-%s-%s" %(brick_path,item[2:4],item[4:6],item[2:10],item[10:14],item[14:18],item[18:22],item[22:])
                info = subprocess.Popen(cmd, shell=True,
                        stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                info.wait()
                stde = info.stderr.readlines()
                if stde:
                    print stde
                    i= i + 1
                    file_log.write('---------find err--------- %s -\n' %i)
                    file_log.write(parent_str+ '\n')
                    file_log.write(file_str + '\n')
                    file_log.write(str(stde) + '\n')

                    stat_dir.write("{}\n".format(tmp_path))

                    # path1 = "%s/.glusterfs/%s" %(brick_path,item[2:4])
                    # if os.path.exists(path1) == False:
                    #     os.makedirs(path1)

                    # path2 = "%s/.glusterfs/%s/%s" %(brick_path,item[2:4],item[4:6])
                    # if os.path.exists(path2) == False:
                    #     os.makedirs(path2)

                    # cmd = ("ln %s %s" %(tmp_path, result))

                    # file_log.write('exe cmd is ' + cmd + '\n')

                    # info = subprocess.Popen(cmd, shell=True,
                    #                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                    # info.wait()
                    # stderr1 = info.stderr.readlines()
                    # if stderr1:
                    #     print stderr1
                    #     file_log.write(stderr1)
                    #     file_log.write('---------heal err----------\n')
                    file_log.write('---------handle err----------\n')
                    print '---------------end------------------'
                continue
                stdout1 = info.stdout.readlines()
                print str(stdout1)
            print '---------------end------------------'
            file_log.flush()
        for file_name in file_path:
            tmp_path = os.path.join(root,file_name)
            index= tmp_path.rfind('/')
            parent_path = tmp_path[:index]

            cmd = "getfattr -dm . -e hex --no-dereference %s | grep gfid= | cut -d '=' -f 2" %(parent_path)
            info = subprocess.Popen(cmd, shell=True,
                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr1 = info.stderr.readlines()
            if not (stderr1 and len(stderr1) == 1 and "Removing leading '/' from absolute path names" in stderr1[0]):
                if "No such file or directory" in stderr1[0]:
                    continue
                err_log.write("{}\n".format(' '.join(stderr1)))
                print stderr1
                print cmd
                return None
            print '---------------start------------------'
            stdout1 = info.stdout.readlines()
            for item in stdout1:
                parent_item = item[:-1]
                result = "%s/.glusterfs/%s/%s/%s-%s-%s-%s-%s" \
                        %(brick_path,parent_item[2:4],parent_item[4:6],parent_item[2:10],\
                        parent_item[10:14],parent_item[14:18],parent_item[18:22],parent_item[22:])
                parent_str =  'parent_path---> %s parent_gfid---> %s' %(parent_path,result)
                print parent_str

            cmd = "getfattr -dm . -e hex --no-dereference %s | grep gfid= | cut -d '=' -f 2" %(tmp_path)
            #cmd = cmd.encode ('utf8')
            print cmd
            info = subprocess.Popen(cmd, shell=True,
                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr1 = info.stderr.readlines()
            if not (stderr1 and len(stderr1) == 1 and "Removing leading '/' from absolute path names" in stderr1[0]):
                if "No such file or directory" in stderr1[0]:
                    continue
                err_log.write("{}\n".format(' '.join(stderr1)))
                print stderr1
                print cmd
                return None

            stdo = info.stdout.readlines()
            for item in stdo:
                item = item[:-1]
                result = "%s/.glusterfs/%s/%s/%s-%s-%s-%s-%s" %(brick_path,item[2:4],item[4:6],item[2:10],item[10:14],item[14:18],item[18:22],item[22:])
                file_str =  "file_path---> %s file_gfid---> %s" %(tmp_path,result)
                print file_str
                cmd = "ls -l %s/.glusterfs/%s/%s/%s-%s-%s-%s-%s" %(brick_path,item[2:4],item[4:6],item[2:10],item[10:14],item[14:18],item[18:22],item[22:])
                info = subprocess.Popen(cmd, shell=True,
                        stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                info.wait()
                stde = info.stderr.readlines()
                if stde and "No such file" in stde[0]:
                    print stde
                    i= i + 1
                    file_log.write('---------find err--------- %s -\n' %i)
                    file_log.write(parent_str+ '\n')
                    file_log.write(file_str + '\n')
                    file_log.write(str(stde) + '\n')

                    stat_file.write("{}\n".format(tmp_path))

                    path1 = "%s/.glusterfs/%s" %(brick_path,item[2:4])
                    if os.path.exists(path1) == False:
                        os.makedirs(path1)

                    path2 = "%s/.glusterfs/%s/%s" %(brick_path,item[2:4],item[4:6])
                    if os.path.exists(path2) == False:
                        os.makedirs(path2)

                    cmd = ("ln %s %s" %(tmp_path, result))

                    file_log.write('exe cmd is ' + cmd + '\n')

                    info = subprocess.Popen(cmd, shell=True,
                                            stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                    info.wait()
                    stderr1 = info.stderr.readlines()
                    if stderr1:
                        print stderr1
                        file_log.write(stderr1)
                        file_log.write('---------heal err----------\n')
                    file_log.write('---------handle err----------\n')
                    print '---------------end------------------'
                elif stde:
                    err_log.write("{}\n".format(" ".join(stde)))
                continue
                stdout1 = info.stdout.readlines()
                print str(stdout1)
            print '---------------end------------------'
            file_log.flush()

if __name__ == '__main__':
    main()