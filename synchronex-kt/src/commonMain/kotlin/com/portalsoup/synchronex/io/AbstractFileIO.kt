package com.portalsoup.synchronex.io

abstract class AbstractFileIO {

    abstract fun copyFile(src: String, dest: String, replace: Boolean = false, owner: String)
    abstract fun isFilePresent(path: String)
    abstract fun listChildren(path:String): String
    abstract fun createFile(path: String)
    abstract fun deleteFile(path: String)
    abstract fun appendToFile(path: String, append: String)
}