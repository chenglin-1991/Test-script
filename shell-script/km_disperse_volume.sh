gluster volume create test-volume disperse 6 redundancy 2 \
10.10.25.21:/root/node21-01/brick21-01/ \
10.10.25.22:/root/node22-01/brick22-01/ \
10.10.25.23:/root/node23-01/brick23-01/ \
10.10.25.21:/root/node21-02/brick22-02/ \
10.10.25.22:/root/node22-02/brick22-02/ \
10.10.25.23:/root/node23-02/brick23-02/ \
force

gluster volume start test-volume force 
./mount.sh
if [ $? == 0 ]
then
	echo "mount seccessful!"
fi

#digiocean volume create test-volume disperse 3 redundancy 1 \
#node21:/root/node21-01/brick21-01/ \
#node22:/root/node22-01/brick22-01/ \
#node23:/root/node23-01/brick23-01/ \
#force

#gluster volume create test-volume disperse 5 redundancy 2 \
#node21:/root/node21-01/brick21-01/ \
#node21:/root/node21-02/brick21-02/ \
#node22:/root/node22-01/brick22-01/ \
#node22:/root/node22-02/brick22-02/ \
#node23:/root/node23-02/brick23-02/ \
#node23:/root/node23-01/brick23-01/ \
#node101:/root/node101-01/brick101-01/ \
#node101:/root/node101-02/brick101-02/ \
#node102:/root/node102-01/brick102-01/ \
#node102:/root/node102-02/brick102-02/ \
#force

#digiocean volume create test-volume disperse-data 5 redundancy 1 \
#node21:/root/node21-01/brick21-01/ \
#node21:/root/node21-02/brick21-02/ \
#node22:/root/node22-01/brick22-01/ \
#node22:/root/node22-02/brick22-02/ \
#node23:/root/node23-01/brick23-01/ \
#node23:/root/node23-02/brick23-02/ \
#force

#digiocean volume start test-volume force 
