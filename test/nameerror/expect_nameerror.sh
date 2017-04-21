cp nameerror.gold nameerror.go
if goderive -dedup . ; then
    echo "expected conflicting function name error, even with deduplication on"
    rm ./derived.gen.go
    rm ./nameerror.go
    exit 1
else
    rm ./nameerror.go
    exit 0
fi
