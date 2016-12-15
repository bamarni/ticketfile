#!/usr/bin/env bats

@test "[All] process a full-featured Ticketfile as ESCP/POS" {
  cmp <(ticket < $BATS_TEST_DIRNAME/fixtures/Ticketfile 2>/dev/null) $BATS_TEST_DIRNAME/fixtures/Ticketfile.expected
}

@test "[All] process a full-featured Ticketfile as HTML" {
  cmp <(ticket -html < $BATS_TEST_DIRNAME/fixtures/Ticketfile 2>/dev/null) $BATS_TEST_DIRNAME/fixtures/Ticketfile.html.expected
}
