cd ../tests
test_files=$(ls | grep .test.sh)
for test in $test_files; do
    source $test
done