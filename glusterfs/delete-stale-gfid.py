#!/usr/bin/python

import re
import sys
import time
import subprocess

if __name__ == '__main__':
    if len(sys.argv) < 5:
        print '''python %s diff_file_absolute_path volume_name mount_point_in_docker glusterfs_container_name''' %(sys.argv[0])
        print '''such as: python %s CANCER/XJ/NJPROJ2Nova_1903_201812/shell/up_ali/.ossutil_output xtvol1 /mnt/test glusterd''' %(sys.argv[0])
        exit(0)

    yes = ''
    s = "Are you share your mount point (%s) exist in your container %s" %(sys.argv[3], sys.argv[4])
    print('\033[1;31;40m%s\033[0m' % s)
    yes = raw_input('y/n:')
    if yes != 'y' and yes != 'yes' and yes != 'Yes':
        exit(0)

    cmd = 'gluster volume info %s|grep Brick|grep -v Number|grep -v Bricks' %(sys.argv[2])
    print cmd
    info = subprocess.Popen(cmd, shell=True,
                            stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    info.wait()
    stderr1 = info.stderr.readlines()
    if stderr1:
        print('\033[1;31;40m%s\033[0m' % stderr1)
        exit(1)

    stdout1 = info.stdout.readlines()
    for item in stdout1:
        list1 = item.split(":")
        hostname = list1[1].split()
        brick_name = list1[2].split()
        brick_name2 = ''.join(brick_name)
        if '@' in brick_name2:
            tmp_index = brick_name2.find('@')
            del brick_name[0]
            brick_name.append(brick_name2[:tmp_index])

        cmd = 'ssh %s ls -al %s/%s' %(hostname[0], brick_name[0], sys.argv[1])
        print cmd
        info = subprocess.Popen(cmd, shell=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        info.wait()
        stderr2 = info.stderr.readlines()
        if stderr2:
            print('\033[1;31;40m%s\033[0m' % stderr2)
            continue

        stdout2 = info.stdout.readlines()
        for dentry in stdout2:
            print('\033[1;34;40m%s\033[0m' % dentry)

    for item in stdout1:
        list1 = item.split(":")
        hostname = list1[1].split()
        brick_name = list1[2].split()
        brick_name2 = ''.join(brick_name)
        if '@' in brick_name2:
            tmp_index = brick_name2.find('@')
            del brick_name[0]
            brick_name.append(brick_name2[:tmp_index])

        cmd = 'ssh %s stat %s/%s' %(hostname[0], brick_name[0], sys.argv[1])
        print cmd
        info = subprocess.Popen(cmd, shell=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        info.wait()
        stderr2 = info.stderr.readlines()
        if stderr2:
            print('\033[1;31;40m%s\033[0m' % stderr2)
            continue

        stdout2 = info.stdout.readlines()
        for dentry in stdout2:
            print('\033[1;34;40m%s\033[0m' % dentry)

    gfid_dict = {}
    for item in stdout1:
        list1 = item.split(":")
        hostname = list1[1].split()
        brick_name = list1[2].split()
        brick_name2 = ''.join(brick_name)
        if '@' in brick_name2:
            tmp_index = brick_name2.find('@')
            del brick_name[0]
            brick_name.append(brick_name2[:tmp_index])

        cmd = 'ssh %s getfattr --absolute-names -dm . -e hex %s/%s |grep -E "gfid|dht"' %(hostname[0], brick_name[0], sys.argv[1])
        print cmd

        info = subprocess.Popen(cmd, shell=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        info.wait()
        stderr2 = info.stderr.readlines()
        if stderr2:
            print('\033[1;31;40m%s\033[0m' % stderr2)
            continue

        stdout3 = info.stdout.readlines()
        for xattr in stdout3:
            print xattr,
            if 'trusted.gfid' in xattr:
                if not gfid_dict.has_key(xattr):
                    gfid_dict[xattr] = 1
                else:
                    num = gfid_dict[xattr] + 1
                    gfid_dict[xattr] = num
        print ''

    print "************ all the gfid is this ************"
    for k,v in gfid_dict.items():
        s = "%s  %s" %(k, v)
        print('\033[1;34;40m%s\033[0m' % s)
        stdout3 = info.stdout.readlines()
        for xattr in stdout3:
            if k in xattr:
                if not gfid_dict.has_key(xattr):
                    gfid_dict[xattr] = 1
                else:
                    num = gfid_dict[xattr] + 1
                    gfid_dict[xattr] = num

    del_gfid_string = None
    yes = ''
    yes = raw_input("Is it continue y/n:")
    if yes == 'y' or yes == 'yes' or yes == 'Yes':
        del_gfid = ''
        del_gfid = raw_input("which gfid file will you delete (such as 0000001):")
        for g_num, v in gfid_dict.items():
            if del_gfid in g_num:
                index = re.split('=|\n', g_num)
                for ix in index:
                    if del_gfid in ix:
                        del_gfid_string = ix
                        break
                break

        if del_gfid_string is None:
            s = '**** we no find the gfid ****'
            print('\033[1;31;40m%s\033[0m' % s)
            exit(0)

        del_gfid_path = '.glusterfs/%s/%s/%s-%s-%s-%s' %(del_gfid_string[2:4],
                del_gfid_string[4:6], del_gfid_string[2:10],del_gfid_string[10:14],del_gfid_string[14:18],del_gfid_string[18:])

        for item in stdout1:
            list1 = item.split(":")
            hostname = list1[1].split()
            brick_name = list1[2].split()
            brick_name2 = ''.join(brick_name)
            if '@' in brick_name2:
                tmp_index = brick_name2.find('@')
                del brick_name[0]
                brick_name.append(brick_name2[:tmp_index])

            del_file_path = "ssh %s rmdir %s/%s" %(hostname[0], brick_name[0], sys.argv[1])
            del_gfid_path2 = "ssh %s sure_to_rm -f %s/%s" %(hostname[0], brick_name[0], del_gfid_path)

            cmd = 'ssh %s getfattr --absolute-names -dm . -e hex %s/%s |grep -E "gfid|dht"' %(hostname[0], brick_name[0], sys.argv[1])

            info = subprocess.Popen(cmd, shell=True,
                                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr2 = info.stderr.readlines()
            if stderr2:
                print('\033[1;31;40m%s\033[0m' % stderr2)
                continue

            stdout4 = info.stdout.readlines()
            for xattr in stdout4:
                if 'trusted.gfid' in xattr:
                    if del_gfid in xattr:
                        s = '************* the cmd will be run ****************'
                        print('\033[1;31;40m%s\033[0m' % s)
                        print ''
                        print('\033[1;34;40m%s\033[0m' % del_file_path)
                        print('\033[1;34;40m%s\033[0m' % del_gfid_path2)
                        print ''
                        s = '*********************** run end *****************'
                        print('\033[1;31;40m%s\033[0m' % s)
                        yes = ''
                        yes = raw_input("Do you sure to delete it y/n:")
                        if yes == 'y' or yes == 'yes' or yes == 'Yes':
                            info = subprocess.Popen(del_file_path, shell=True,
                                                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                            info.wait()
                            stderr2 = info.stderr.readlines()
                            if stderr2:
                                print('\033[1;31;40m%s\033[0m' % stderr2)
                                exit(0)

                            info = subprocess.Popen(del_gfid_path2, shell=True,
                                                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                            info.wait()
                            stderr2 = info.stderr.readlines()
                            if stderr2:
                                print('\033[1;31;40m%s\033[0m' % stderr2)
                                exit(0)
                        else:
                            exit(0)
        s = '************* delete complete *************\n'
        print('\033[1;34;40m%s\033[0m' % s)
        s = '********** we will heal it in mount_point_in_docker ******************\n'
        print('\033[1;34;40m%s\033[0m' % s)

        cmd = "docker exec -t %s stat %s/%s" %(sys.argv[4], sys.argv[3], sys.argv[1])
        s = cmd
        print('\033[1;34;40m%s\033[0m' % s)
        info = subprocess.Popen(cmd, shell=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        info.wait()
        stderr2 = info.stderr.readlines()
        if stderr2:
            print('\033[1;31;40m%s\033[0m' % stderr2)
            #exit(0)

        time.sleep(2)
        s = '*********** we will check %s in brick path again ***********' %(sys.argv[1])
        gfid_dict = {}
        print('\033[1;34;40m%s\033[0m' % s)

        for item in stdout1:
            list1 = item.split(":")
            hostname = list1[1].split()
            brick_name = list1[2].split()
            brick_name2 = ''.join(brick_name)
            if '@' in brick_name2:
                tmp_index = brick_name2.find('@')
                del brick_name[0]
                brick_name.append(brick_name2[:tmp_index])

            cmd = 'ssh %s ls -al %s/%s' %(hostname[0], brick_name[0], sys.argv[1])
            print cmd
            info = subprocess.Popen(cmd, shell=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr2 = info.stderr.readlines()
            if stderr2:
                print('\033[1;31;40m%s\033[0m' % stderr2)
                continue

            stdout2 = info.stdout.readlines()
            for dentry in stdout2:
                print('\033[1;34;40m%s\033[0m' % dentry)

        for item in stdout1:
            list1 = item.split(":")
            hostname = list1[1].split()
            brick_name = list1[2].split()
            brick_name2 = ''.join(brick_name)
            if '@' in brick_name2:
                tmp_index = brick_name2.find('@')
                del brick_name[0]
                brick_name.append(brick_name2[:tmp_index])

            cmd = 'ssh %s getfattr --absolute-names -dm . -e hex %s/%s |grep -E "gfid|dht"' %(hostname[0], brick_name[0], sys.argv[1])
            print cmd

            info = subprocess.Popen(cmd, shell=True,
                                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            info.wait()
            stderr2 = info.stderr.readlines()
            if stderr2:
                print('\033[1;31;40m%s\033[0m' % stderr2)
                continue
            stdout5 = info.stdout.readlines()
            for xattr in stdout5:
                print xattr,
                if 'trusted.gfid' in xattr:
                    if not gfid_dict.has_key(xattr):
                        gfid_dict[xattr] = 1
                    else:
                        num = gfid_dict[xattr] + 1
                        gfid_dict[xattr] = num
            print ''

        print "************ all the gfid is this ************"
        for k,v in gfid_dict.items():
            s = "%s  %s" %(k, v)
            print('\033[1;34;40m%s\033[0m' % s)
    else:
        exit(0)

