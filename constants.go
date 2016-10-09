package gosmile

const token_prefix_small_int = 0xC0
const token_prefix_integer = 0x24
const token_prefix_fp = 0x28

const token_misc_integer = 0x00
const token_misc_float_32 = 0x00
const token_misc_float_64 = 0x01

const token_byte_int_32 = token_prefix_integer + token_misc_integer
const token_byte_float_32 = token_prefix_fp | token_misc_float_32
