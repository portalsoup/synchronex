package com.portalsoup.synchronex.io

class FileIO: AbstractFileIO() {

    override fun copyFile(src: String, dest: String, replace: Boolean, owner: String) {
        // If destination already exists and replace == true
        //      move existing file to cache folder
        // else
        //      noop


    }

    override fun isFilePresent(path: String) {
        TODO("Not yet implemented")
    }

    override fun listChildren(path: String): String {
        TODO("Not yet implemented")
    }

    override fun createFile(path: String) {
        TODO("Not yet implemented")
    }

    override fun deleteFile(path: String) {
        TODO("Not yet implemented")
    }

    override fun appendToFile(path: String, append: String) {
        TODO("Not yet implemented")
    }
}