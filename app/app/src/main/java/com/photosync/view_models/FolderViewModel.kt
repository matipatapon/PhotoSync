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
import okhttp3.MediaType
import okhttp3.MultipartBody
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import java.time.Instant
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.util.logging.Logger


class FolderViewModel(
    localDatabase: LocalDatabase,
    private var application: Application
) : ViewModel(){

    private val client: OkHttpClient = OkHttpClient()
    private val logger = Logger.getLogger(this.javaClass.name)
    private val _folders = MutableStateFlow(listOf<String>())
    val folders = _folders.asStateFlow()
    private val folderDao = localDatabase.folderDao()
    private val appSettingsDao = localDatabase.appSettingsDao()

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

    private fun uploadFile(file: ByteArray, filename: String, lastModified: String, mimeType: String, token: String){
        val requestBody: MultipartBody = MultipartBody.Builder()
            .setType(MultipartBody.FORM)
            .addFormDataPart("filename", filename)
            .addFormDataPart("modification_date", lastModified)
            .addFormDataPart("file", "", RequestBody.create(MediaType.parse("image/jpeg"),file))
            .build()
        val server = appSettingsDao.getSettings()!!.server
        val request = Request.Builder()
            .header("Authorization", token)
            .url("$server/v1/upload")
            .post(requestBody)
            .build()
        val response = client.newCall(request).execute()
        val responseCode = response.code()
    }

    private fun syncFile(file: DocumentFile, token: String){
        val fileUri = file.uri
        val filename = file.uri.path.toString().substringAfterLast("/")
        val inputStream = application.contentResolver.openInputStream(file.uri)
        if(inputStream == null)
        {
            throw Exception("$fileUri not found")
        }
        val bytes = inputStream.use {it.readBytes()}
        val fileLastModified = Instant.ofEpochMilli(file.lastModified())
            .atZone(ZoneId.systemDefault())
            .format(DateTimeFormatter.ofPattern("uuuu.MM.dd HH:mm:ss"))
        val mimeType: String? = file.type
        if(mimeType == null || mimeType != "image/jpeg"){
            return
        }
        logger.info("Uploading file<name={$filename} size={${bytes.size}} lastModified={$fileLastModified} path={${fileUri.path}} mimeType={$mimeType}>")
        uploadFile(bytes, filename, fileLastModified, mimeType, token)
    }

    private fun syncFolder(folder: DocumentFile, lastSync: Long?, token: String){
        var folderLastSynced = ""
        val folderLastModified = Instant.ofEpochMilli(folder.lastModified())
            .atZone(ZoneId.systemDefault())
            .format(DateTimeFormatter.ofPattern("uuuu.MM.dd HH:mm:ss"))
        if(lastSync != null) {
            folderLastSynced = Instant.ofEpochMilli(folder.lastModified())
                .atZone(ZoneId.systemDefault())
                .format(DateTimeFormatter.ofPattern("uuuu.MM.dd HH:mm:ss"))
            if(lastSync > folder.lastModified()){
                logger.info("Skipping folder<path={${folder.uri.path}} lastModified={$folderLastModified} lastSynced={$folderLastSynced}>")
                return
            }
        }
        logger.info("Syncing folder<path={${folder.uri.path}} lastModified={$folderLastModified} lastSynced={$folderLastSynced}>> ")
        for(file in folder.listFiles()){
            if(file.isDirectory){
                syncFolder(file, lastSync, token)
                continue
            }
            syncFile(file, token)
        }
    }

    fun syncFolders(token: String){
        viewModelScope.launch(Dispatchers.IO) {
            for (folder in folderDao.getFolders()) {
                val directory = DocumentFile.fromTreeUri(
                    application.applicationContext,
                    folder.uri.toUri()
                )
                if (directory != null && directory.isDirectory) {
                    syncFolder(directory, folder.lastSync, token)
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
