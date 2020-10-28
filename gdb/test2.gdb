set logging file frame.log
set logging on
set pagination off
set print repeats 100

define print_frame
    bt
    f 10
    set $stub=((call_stub_t *)((char *)(async)-(unsigned long)(&((call_stub_t *)0)->async)))
    set $frame=$stub->frame
    p ($frame)->this->name
    p ($frame)->wind_from
    set $head=&($frame->root->myframes)
    set $tmp_head=$head
    set $next_frame=$head->next
    while $head != $next_frame
        set $num=1
        set $tmp_frame=((call_frame_t *)((char *)(($tmp_head)->next)-(unsigned long)(&((call_frame_t *)0)->frames)))
        p ($tmp_frame)->this->name
        p ($tmp_frame)->wind_from
        set $tmp_head=$tmp_head->next
        set $next_frame=$next_frame->next
