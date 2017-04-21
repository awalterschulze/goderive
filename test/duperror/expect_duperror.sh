cp duperror.gold duperror.go
if goderive -autoname . ; then
    echo "expected ambigious function name error, even with autoname on"
    rm ./derived.gen.go
    rm ./duperror.go
    exit 1
else
    rm ./duperror.go
    exit 0
fi
