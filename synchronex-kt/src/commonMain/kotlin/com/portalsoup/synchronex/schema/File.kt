package com.portalsoup.synchronex.schema

import com.portalsoup.synchronex.dsl.FileActions
import com.portalsoup.synchronex.execution.Job
import com.portalsoup.synchronex.execution.state.ResourceState
import kotlinx.serialization.Serializable

@Serializable
data class File(
    val action: FileActions,
    val source: String,
    val destination: String,
    val user: String,
    val group: String,
    val chmod: String
): Job {
    override fun validate(): Boolean {
        TODO("Not yet implemented")
        // verify source path as file
    }

    override fun execute() {
        TODO("Not yet implemented")
    }

    override fun state(): ResourceState {
        return ResourceState(
            "file",
            id = destination,
            attributes = mapOf(
                "source" to source,
                "destination" to destination,
                "user" to user,
                "group" to group,
                "chmod" to chmod,
            )
        )
    }
}