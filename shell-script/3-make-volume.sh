#rm -rf /usr/var/lib/digioceand/*
#service digioceand restart
#digiocean volume create test-volume replica 3 arbiter 1 node21:/root/node21-01/brick21-01/ node22:/root/node22-01/brick22-01/ node23:/root/node23-01/brick23-01/ force
gluster volume create test-volume replica 3 node21:/root/node21-01/brick21-01/ node22:/root/node22-01/brick22-01/ node23:/root/node23-01/brick23-01/ force

gluster volume start test-volume force 

#sleep 4

#mount.digioceanfs node59:/test-volume ~/mount2/
