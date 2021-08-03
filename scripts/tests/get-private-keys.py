#!/usr/bin/python

import sys

if __name__ == '__main__':
    private_keys = ''
    private_keys_count = int(sys.argv[1])
    if not private_keys_count:
        raise Exception('Number of private keys should be provided as parameter')

    last_hex_number = 0x1002030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff
    for x in range(0, private_keys_count + 1):
        last_hex_number += x
        private_keys += ',' + str(hex(last_hex_number).replace('L', '').replace('0x', ''))

    sys.stdout.write(private_keys.strip(','))
