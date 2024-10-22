package com.portalsoup.synchronex.schema

import com.portalsoup.synchronex.execution.Job
import com.portalsoup.synchronex.execution.state.ResourceState
import kotlinx.serialization.Serializable

@Serializable
data class Nex(
    val user: String,
    val batches: List<Nex>,
    val files: List<File>
): Job {
    override fun validate(): Boolean {
        TODO("Not yet implemented")
    }

    override fun execute() {
        TODO("Not yet implemented")
    }

    override fun state(): ResourceState {
        TODO("Not yet implemented")
    }
}