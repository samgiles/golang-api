#!/usr/bin/env sh
set -e

cleanup () {
  docker-compose -p inttests kill
  docker-compose -p inttests rm -f
}

trap 'cleanup; printf "Tests failed unexpectedly\n"' HUP INT QUIT PIPE TERM

cd "$(dirname "$0")"

docker-compose -p inttests build && docker-compose -p inttests up -d

if [ $? -ne 0 ] ; then
  printf "Docker Compose Failed\n"
  exit -1
fi

docker logs -f inttests_int_tests_1

TEST_EXIT_CODE=`docker wait inttests_int_tests_1`

if [ -z ${TEST_EXIT_CODE+x} ] || [ "$TEST_EXIT_CODE" -ne 0 ] ; then
  printf "Tests Failed with exit code: $TEST_EXIT_CODE\n"
else
  printf "Tests Passed\n"
fi

cleanup

exit $TEST_EXIT_CODE
