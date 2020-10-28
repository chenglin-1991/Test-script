# -*- encoding: utf-8 -*-

import json
import os
import string
import time
import random
import mpi4py


def do_read(content_info, file_name):
    fd = os.open(file_name, os.O_RDONLY)
    if fd == -1:
        return - 1, 0

    os.lseek(fd, content_info["offset"], os.SEEK_SET)
    read_str = os.read(fd, len(content_info["buf"]))
    if len(read_str) != len(content_info["buf"]):
        print("Good buf len={}, Bad buf len={}".format(len(content_info["buf"], len(read_str))))
        print("Good Buf: {}".format(json.dumps(content_info, indent=4)))
        print("Bad Buf: {}".format(read_str))
        return -1, 0

    os.close(fd)
    return 0, read_str


def do_write(buf, file_name):
    fd = os.open(file_name, os.O_RDWR | os.O_CREAT)
    if fd == -1:
        return - 1, 0

    offset = random.randint(0, MAX_OFFSET)
    os.lseek(fd, offset, os.SEEK_SET)
    ret = os.write(fd, buf)
    if ret != len(buf):
        return -1, 0
    os.close(fd)
    return 0, offset


def do_op(rank, write_rank, file_name):
    content_info = None
    if rank == write_rank:
        content_size = random.randint(1, CONTENT_MAX_SIZE)
        content_str = ''.join(random.choice(string.ascii_letters)
                              for i in range(content_size))
        print("write rank {}, start writing {} ...".format(rank, file_name))
        ret, offset = do_write(content_str, file_name)
        if ret == -1:
            return - 1
        content_info = {
            "offset": offset,
            "buf": content_str
        }
        print("write rank {}, writing {} done...".format(rank, file_name))
        content_info = COMM.bcast(content_info, root=write_rank)
    else:
        content_info = COMM.bcast(content_info, root=write_rank)
        print("read rank {}, start reading {} ...".format(rank, file_name))
        ret, read_str = do_read(content_info, file_name)
        if ret == -1:
            return - 1
        print("read rank {},reading {} done...".format(rank, file_name))
        if read_str == content_info["buf"]:
            print("Check success!")
        else:
	    print("Good buf len={}, Bad buf len={}".format(len(content_info["buf"]), len(read_str)))
	    print("Good Buf {}".format(content_info['buf']))
	    print("Bad Buf: {}".format(read_str))
            return -1


def main(file_index):
    random.seed(time.time())
    rank = COMM.Get_rank()

    for i in range(3):
        write_rank = None
        if rank == 0:
            size = COMM.Get_size()
            #if size != 3:
                #print("should be three client..")
                #return
            file_name = "/mnt/testl/alamo3-stripe%d" % file_index
            write_rank = random.randint(0, size-1)
            #print("write rank is {}, file_name={}".format(
            #    (write_rank), (file_name)))
        else:
            file_name = None
        file_name = COMM.bcast(file_name, root=0)
        write_rank = COMM.bcast(write_rank, root=0)
        ret = do_op(rank, write_rank, file_name)
        if ret == -1:
            return - 1

        COMM.barrier()

    return 0


if __name__ == "__main__":
    global CONTENT_MAX_SIZE
    global MAX_OFFSET
    global COMM

    time1  = time.time()

    localTime = time.localtime(time.time())

    localTimeStrs = time.strftime("%Y-%m-%d %H:%M:%S",localTime)
    print localTimeStrs

    from mpi4py import MPI
    CONTENT_MAX_SIZE = 1000
    MAX_OFFSET = 100000
    COMM = MPI.COMM_WORLD
    if COMM.Get_rank() == 0:
        file_index_list = random.sample(range(1, 100), 10)
    else:
        file_index_list = [None for i in range(10)]
    #main(1)
    for i in range(10000000):
        for j in range(2000):
            print("normal running index {}".format(j))
            ret = main(j)
            if ret:
                print("index={}, ret={}".format(j, ret))
                break
        if ret:
            break
