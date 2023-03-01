#!/bin/sh

while ! curl -f $PLAYWRIGHT_BASE_URL -o /dev/null; do sleep 3; done

exec "$@"
