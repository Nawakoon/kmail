fail=false
total_tests=0
fail_tests=0
pass_tests=0
skip_tests=0

function run_test() {
    total_tests=$((total_tests+1))
    if [ "$2" ]; then
        pass_tests=$((pass_tests+1))
        echo "✓ pass\t $1"
    else
        fail_tests=$((fail_tests+1))
        echo "✕ fail\t $1"
        fail=true
    fi
}

function skip_test() {
    total_tests=$((total_tests+1))
    skip_tests=$((skip_tests+1))
    echo "- skip\t $1"
    echo "\b * $2"
}

function summary() {
    if [ "$fail" = true ]; then
        echo "failed\t $fail_tests/$total_tests test(s)"
        exit 1
    fi
    if [ "$skip_tests" -gt 0 ]; then
        echo "skipped\t $skip_tests test(s)"
    fi
}