#!/usr/bin/python

import sys

if __name__ == '__main__':
    highest = ''
    highest_int = int(sys.argv[1])
    highest_start_from = int(sys.argv[2])
    if not highest_int:
        raise Exception('Highest int should be provided as parameter')

    for x in range(highest_start_from, highest_start_from + highest_int + 1):
        highest += ',' + str(x)

    sys.stdout.write(highest.strip(','))
