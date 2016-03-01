#!/usr/bin/env bats

@test "[Charset] print some charset-specific text" {
  cmp <(printer < $BATS_TEST_DIRNAME/fixtures/charset 2>/dev/null) $BATS_TEST_DIRNAME/fixtures/charset.expected
}

@test "[Charset] fail for unsupported charsets" {
  run printer < $BATS_TEST_DIRNAME/fixtures/charset_unexisting

  [[ "$output" == *"Charset TOUTI not supported"* ]] || false

  [ "$status" -ne 0 ]
}

@test "[Charset] fail for unencodable unicode" {
  run printer < $BATS_TEST_DIRNAME/fixtures/charset_unencodable

  [[ "$output" == *"Couldn't encode to charset"* ]] || false

  [ "$status" -ne 0 ]
}
