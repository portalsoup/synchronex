package com.portalsoup.synchronex.dsl

import com.portalsoup.synchronex.schema.File
import kotlinx.serialization.Serializable

@Serializable
enum class FileActions { SYNC }

class FileBuilder(val action: FileActions): JobBuilder<FileBuilder>() {
    var src: String = ""
    var dest: String = ""
    var user: String = ""
    var group: String = ""
    var mode: String = ""

    init {
        parent = this
    }

    fun build(): File = File(
        action = action,
        source = src,
        destination = dest,
        user = user,
        group = group,
        chmod = mode,
    )
}

class FilesBuilder {

    val files = mutableListOf<FileBuilder>()

    fun sync(dest: String, lambda: FileBuilder.() -> Unit) {
        val syncBuilder = FileBuilder(FileActions.SYNC).apply(lambda)
        syncBuilder.dest = dest

        files.add(syncBuilder)
    }
    fun build() = files
}