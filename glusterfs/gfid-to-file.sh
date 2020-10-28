#!/bin/bash
 
if [[ "$#" < "2" || "$#" > "3" ]]; then
  cat <<END
Glusterfs GFID resolver -- turns a GFID into a real file path
 
Usage: $0 <brick-path> <gfid> [-q]
  <brick-path> : the path to your glusterfs brick (required)
  
  <gfid> : the gfid you wish to resolve to a real path (required)
  
  -q : quieter output (optional)
       with this option only the actual resolved path is printed.
       without this option $0 will print the GFID, 
       whether it identifies a file or directory, and the resolved
       path to the real file or directory.
 
Theory:
The .glusterfs directory in the brick root has files named by GFIDs
If the GFID identifies a directory, then this file is a symlink to the
actual directory.  If the GFID identifies a file then this file is a
hard link to the actual file.
END
exit
fi

BRICK="$1"
 
GFID="$2"
GP1=`cut -c 1-2 <<<"$GFID"`
GP2=`cut -c 3-4 <<<"$GFID"`
GFIDPRE="$BRICK"/.glusterfs/"$GP1"/"$GP2"
GFIDPATH="$GFIDPRE"/"$GFID"
 
if [[ "$#" == "2" ]]; then
  echo -ne "$GFID\t==\t"
fi

 
if [[ -h "$GFIDPATH" ]]; then
  if [[ "$#" == "2" ]]; then
    echo -ne "Directory:\t"
  fi
  DIRPATH="$GFIDPRE"/`readlink "$GFIDPATH"`
  echo $(cd $(dirname "$DIRPATH"); pwd -P)/$(basename "$DIRPATH")
else
  if [[ "$#" == "2" ]]; then
    echo -ne "File:\t"
  fi
  INUM=`ls -i "$GFIDPATH" | cut -f 1 -d \ `  
  find "$BRICK" -inum "$INUM" ! -path \*.glusterfs/\*
fi
