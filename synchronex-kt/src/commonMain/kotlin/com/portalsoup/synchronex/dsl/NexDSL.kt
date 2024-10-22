package com.portalsoup.synchronex.dsl

import com.portalsoup.synchronex.schema.Nex

abstract class JobBuilder<J: JobBuilder<J>> {
    internal lateinit var parent: J

    fun pre(lambda: J.() -> Unit) {
        if (!::parent.isInitialized) {
            throw RuntimeException("Didn't initialize parent value!")
        }
        parent.apply(lambda)
    }

    fun post(lambda: J.() -> Unit) {
        if (!::parent.isInitialized) {
            throw RuntimeException("Didn't initialize parent value!")
        }
        parent.apply(lambda)
    }
}

class NexBuilder: JobBuilder<NexBuilder>() {
    init {
        parent = this
    }

    lateinit var user: String

    private var batches = mutableListOf<NexBuilder>()
    private var filesList = mutableListOf<FileBuilder>()

    fun batch(name: String, lambda: NexBuilder.() -> Unit) {
        batches.add(NexBuilder().apply(lambda))
    }

    fun sync(dest: String, lambda: FileBuilder.() -> Unit) {
        val syncBuilder = FileBuilder(FileActions.SYNC).apply(lambda)
        syncBuilder.dest = dest

        filesList.add(syncBuilder)
    }

    fun build(): Nex = Nex(
        user = user,
        batches = batches.map { it.build() },
        files = filesList.map { it.build() }
    )
}


fun nex(lambda: NexBuilder.() -> Unit) = NexBuilder().apply(lambda)
