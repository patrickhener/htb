#!/bin/bash
until ./devchat; do
    echo "Server 'devchat' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done

