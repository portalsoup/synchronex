package com.portalsoup.synchronex.execution

import com.portalsoup.synchronex.execution.state.ResourceState

interface Job {
    fun validate(): Boolean
    fun execute()
    fun state(): ResourceState
}