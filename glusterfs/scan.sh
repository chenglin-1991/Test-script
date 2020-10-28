for i in {1..10}; do ssh xt${i} "grep 'gfid diff' /var/log/glusterfs/nfs.log"; done | perl -lne "print /(?<=-dht: ).*(?=: gfid diff)/g" | sort -u
