#!/usr/bin/env bats

@test "[Charset] print some charset-specific text" {
  cmp <(ticket < $BATS_TEST_DIRNAME/fixtures/charset 2>/dev/null) $BATS_TEST_DIRNAME/fixtures/charset.expected
}

@test "[Charset] fail for unsupported charsets" {
  run ticket < $BATS_TEST_DIRNAME/fixtures/charset_unexisting

  [[ "$output" == *"charset TOUTI not supported"* ]] || false

  [ "$status" -ne 0 ]
}

@test "[Charset] fail for unencodable unicode" {
  run ticket < $BATS_TEST_DIRNAME/fixtures/charset_unencodable

  [ "$status" -ne 0 ]
}
