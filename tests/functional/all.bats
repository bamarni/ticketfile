#!/usr/bin/env bats

@test "[All] process a full-featured Ticketfile" {
  cmp <(printer < $BATS_TEST_DIRNAME/fixtures/Ticketfile 2>/dev/null) $BATS_TEST_DIRNAME/fixtures/Ticketfile.expected
}
