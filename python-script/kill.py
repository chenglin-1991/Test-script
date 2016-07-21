#!/usr/bin/python
import os
import time

while (1):
    time.sleep (50)
    os.popen('gluster v start test-volume force')
    time.sleep (5)
    os.popen('gluster v heal test-volume full')
