package com.photosync.view_models

import android.app.Application
import android.net.Uri
import androidx.core.net.toUri
import androidx.documentfile.provider.DocumentFile
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.photosync.api.ApiHandler
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

data class FolderStatus(
    val type: Type,
    val info: String
){
    enum class Type {
        Idle,
        Sync,
        Error
    }
}

class FolderViewModel(
    localDatabase: LocalDatabase,
    private var application: Application,
    private var apiHandler: ApiHandler
) : ViewModel(){

    private val logger = Logger.getLogger(this.javaClass.name)
    private val _folders = MutableStateFlow(listOf<Folder>())
    val folders = _folders.asStateFlow()
    private val folderDao = localDatabase.folderDao()
    private val _status = MutableStateFlow(FolderStatus(FolderStatus.Type.Idle, ""))
    val status = _status.asStateFlow()

    private fun refreshFolders(){
        val newFolders = mutableListOf<Folder>()
        for (folder in folderDao.getFolders()){
            newFolders.add(folder)
        }
        _folders.value = newFolders
    }

    init{
        viewModelScope.launch(Dispatchers.IO) {
            refreshFolders()
        }
    }

    private fun syncFile(file: DocumentFile, lastSync: Long?){
        val filename = file.uri.path.toString().substringAfterLast("/")
        val filepath = file.uri.path
        val fileLastModifiedUnix = file.lastModified()
        val fileLastModified = Instant.ofEpochMilli(fileLastModifiedUnix)
            .atZone(ZoneId.systemDefault())
            .format(DateTimeFormatter.ofPattern("uuuu.MM.dd HH:mm:ss"))
        val mimeType: String? = file.type
        if(lastSync != null && lastSync > fileLastModifiedUnix ){
            logger.info("Ignoring file<path={$filepath} lastModified={$fileLastModified}>")
            return
        }
        if(mimeType == null || mimeType != "image/jpeg"){
            logger.info("Ignoring file<path={$filepath} mimeType={$mimeType}>")
            return
        }
        _status.value = FolderStatus(FolderStatus.Type.Sync, filename)
        val inputStream = application.contentResolver.openInputStream(file.uri)
        if(inputStream == null)
        {
            throw Exception("$filepath not found")
        }
        val bytes = inputStream.use {it.readBytes()}
        logger.info("Uploading file<path={$filepath} size={${bytes.size}} lastModified={$fileLastModified}  mimeType={$mimeType}>")
        apiHandler.uploadFile(bytes, filename, fileLastModified)
    }

    private fun syncFolder(folder: DocumentFile, lastSync: Long?){
        for(file in folder.listFiles()){
            if(file.isDirectory){
                syncFolder(file, lastSync)
                continue
            }
            syncFile(file, lastSync)
        }
    }

    fun syncFolders(){
        _status.value = FolderStatus(FolderStatus.Type.Sync, "")
        viewModelScope.launch(Dispatchers.IO) {
            try {
                for (folder in folderDao.getFolders()) {
                    val directory = DocumentFile.fromTreeUri(
                        application.applicationContext,
                        folder.uri.toUri()
                    )
                    if (directory == null || !directory.isDirectory) {
                        throw Exception("Invalid folder")
                    }
                    val currentTime = System.currentTimeMillis()
                    syncFolder(directory, folder.lastSync)
                    folder.lastSync = currentTime
                    folderDao.updateFolder(folder)
                    _status.value = FolderStatus(FolderStatus.Type.Idle, "")
                }
            } catch (e: Exception){
                _status.value = FolderStatus(FolderStatus.Type.Error, e.toString())
            }
        }
    }

    fun addFolderToSync(uri: Uri){
        viewModelScope.launch(Dispatchers.IO){
            try {
                val folders = folderDao.getFolders()
                val uriStr = uri.toString()
                for(folder in folders){
                    if(uriStr.startsWith(folder.uri)){
                        return@launch
                    }
                    if(folder.uri.startsWith(uriStr)){
                        folderDao.deleteFolder(folder)
                    }
                }
                folderDao.addFolder(Folder(uriStr, null))
                refreshFolders()
            } catch (e: Exception){
                _status.value = FolderStatus(FolderStatus.Type.Error, e.toString())
            }
        }
    }

    fun deleteFolder(folder: Folder){
        viewModelScope.launch(Dispatchers.IO){
            folderDao.deleteFolder(folder)
            refreshFolders()
        }
    }

    fun resetStatus(){
        _status.value = FolderStatus(FolderStatus.Type.Idle, "")
    }
}
