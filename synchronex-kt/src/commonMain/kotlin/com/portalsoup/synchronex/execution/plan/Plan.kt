package com.portalsoup.synchronex.execution.plan

/*
 * Inspired by .tfstate.  Rather than just inline running the validation on each item and then conditionally executing
 * it, build the plan from the state file, which can be used to preserve the case where tracked files are removed and
 * require "untracked" managing to clean up on a subsequent run
 */
class Plan {
}