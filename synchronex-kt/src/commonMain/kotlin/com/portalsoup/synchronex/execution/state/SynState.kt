package com.portalsoup.synchronex.execution.state

import kotlinx.serialization.Serializable

/*
 * Inspired by .tfstate.  Rather than just inline running the validation on each item and then conditionally executing
 * it, record the expected state of each resource applied.
 *
 * Make plan an explicit subcommand to run a comparison check between the expected and actual and print differences
 * Then running apply simply does the work, updating in place resources that change, and leaving alone existing
 * untouched resources
 */
@Serializable
data class SynState(
    val state: List<ResourceState>
)

