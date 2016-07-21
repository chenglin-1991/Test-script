umount /cluster2/test/
#digiocean volume stop test-volume 
#digiocean volume delete test-volume
gluster volume stop test-volume 
gluster volume delete test-volume 

./2-fdisk.sh
if [ $? -eq 0 ]
then
	echo "====================ssh  node21 ok!==========================="
fi

./ssh.sh
#./mk_disperse_volume.sh
