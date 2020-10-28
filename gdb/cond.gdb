set logging file gdb.log
set logging on
set pagination off
set print repeats 100

define print_frame
    b ngc_lookup
    b ngc_link
    b dht_lookup
    b dht_link
    b ngc_unlink

    commands 1
       p *loc
       c
    end

    commands 2
       p *newloc
       p *oldloc
       c
    end

    commands 3
       p *loc
       c
    end

    commands 4
       p *newloc
       p *oldloc
       c
    end

    commands 5
       p *loc
       c
    end

    c
