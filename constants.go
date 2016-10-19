package gosmile

const token_prefix_small_int = 0xC0
const token_prefix_integer = 0x24
const token_prefix_fp = 0x28

const token_prefix_tiny_ascii = 0x40
const token_prefix_small_ascii = 0x60
const token_prefix_tiny_unicode = 0x80
const token_prefix_short_unicode = 0xA0

const token_misc_integer = 0x00
const token_misc_float_32 = 0x00
const token_misc_float_64 = 0x01
const token_misc_long_text_ascii = 0xE0
const token_misc_long_text_unicode = 0xE4

const token_byte_int_32 = token_prefix_integer + token_misc_integer
const token_byte_float_32 = token_prefix_fp | token_misc_float_32
const token_byte_float_64 = token_prefix_fp | token_misc_float_64

const token_byte_long_string_ascii = token_misc_long_text_ascii

const token_literal_empty_string = 0x20
const token_literal_null = 0x21
const token_literal_false = 0x22
const token_literal_true = 0x23

const max_short_value_string_bytes = 64

const int_marker_end_of_string = 0xFC
const byte_marker_end_of_string = int_marker_end_of_string
