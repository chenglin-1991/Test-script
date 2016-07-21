i=0
while ((i < 20))
do
	gluster v heal test-volume full
	sleep 1800
	sleep 1800
	((i ++))
done

ifup eno16777984
