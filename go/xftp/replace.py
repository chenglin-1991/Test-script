#!/usr/bin/python

import sys
import json

if __name__ == '__main__':
    print sys.argv
    json_data = {}

    with open('./xftp.json')as fp:
        json_data = json.load(fp)


    for k,v in json_data.items():
        print k,v

    release = sys.argv[1]
    with open('info.json','w') as fp:
        json_data['release'] = release
        json.dump(json_data,fp)
    print 'end'
