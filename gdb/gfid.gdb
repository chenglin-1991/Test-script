set logging file gdb.log
set logging on
set pagination off
set print repeats 100

define print_frame
    bt
    set $var_num=232
    while $var_num > 5
        f $var_num
        set $var=gfid
        printf "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x\n",$var[0],$var[1],$var[2],$var[3],$var[4],$var[5],$var[6],$var[7],$var[8],$var[9],$var[10],$var[11],$var[12],$var[13],$var[14],$var[15]
        set $var_num=$var_num-1
