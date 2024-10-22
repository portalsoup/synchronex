package com.portalsoup.synchronex.execution.state

import kotlinx.serialization.Serializable

@Serializable
data class ResourceState(
    val type: String,
    val id: String,
    val attributes: Map<String, String>
)