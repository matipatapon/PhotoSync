package com.photosync.view_models

import android.app.Application
import android.database.sqlite.SQLiteConstraintException
import android.net.Uri
import androidx.core.net.toUri
import androidx.documentfile.provider.DocumentFile
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.photosync.database.Folder
import com.photosync.database.LocalDatabase
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.time.Instant
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.util.logging.Logger


class FolderViewModel(
    localDatabase: LocalDatabase,
    private var application: Application
) : ViewModel(){

    private val logger = Logger.getLogger(this.javaClass.name)
    private val _folders = MutableStateFlow(listOf<String>())
    val folders = _folders.asStateFlow()
    private val folderDao = localDatabase.folderDao()

    private val _error = MutableStateFlow("")
    val error = _error.asStateFlow()

    private fun refreshFolders(){
        val newFolders = mutableListOf<String>()
        for (folder in folderDao.getFolders()){
            newFolders.add(folder.uri.toUri().path.toString())
        }
        _folders.value = newFolders
    }

    init{
        viewModelScope.launch(Dispatchers.IO) {
            refreshFolders()
        }
    }

    private fun syncFile(file: DocumentFile){
        val fileUri = file.uri
        val fileName = file.uri.path.toString().substringAfterLast("/")
        val inputStream = application.contentResolver.openInputStream(file.uri)
        if(inputStream == null)
        {
            throw Exception("$fileUri not found")
        }
        val bytes = inputStream.use {it.readBytes()}
        val fileLastModified = Instant.ofEpochMilli(file.lastModified())
            .atZone(ZoneId.systemDefault())
            .format(DateTimeFormatter.ofPattern("uuuu.MM.dd HH:mm:ss"))
        logger.info("File<name={$fileName} size={${bytes.size}} lastModified={$fileLastModified} path={${fileUri.path}}>")
    }

    private fun syncFolder(folder: DocumentFile){
        logger.info("Folder<path={${folder.uri.path}}>")
        for(file in folder.listFiles()){
            if(file.isDirectory){
                syncFolder(file)
                continue
            }
            syncFile(file)
        }
    }

    fun syncFolders(){
        viewModelScope.launch(Dispatchers.IO) {
            for (folder in folderDao.getFolders()) {
                val directory = DocumentFile.fromTreeUri(
                    application.applicationContext,
                    folder.uri.toUri()
                )
                if (directory != null && directory.isDirectory) {
                    syncFolder(directory)
                }
            }
        }
    }

    fun addFolderToSync(uri: Uri){
        viewModelScope.launch(Dispatchers.IO){
            try {
                folderDao.addFolder(Folder(uri.toString(), null))
                refreshFolders()
            } catch (e: SQLiteConstraintException){
                _error.value = e.toString()
            }
        }
    }
}
