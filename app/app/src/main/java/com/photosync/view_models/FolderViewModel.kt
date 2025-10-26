package com.photosync.view_models

import android.database.sqlite.SQLiteConstraintException
import android.net.Uri
import androidx.core.net.toUri
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.photosync.database.Folder
import com.photosync.database.LocalDatabase
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch


class FolderViewModel(private var localDatabase: LocalDatabase) : ViewModel(){
    private val _folders = MutableStateFlow(listOf<String>());
    val folders = _folders.asStateFlow()
    private val folderDao = localDatabase.folderDao()

    private val _error = MutableStateFlow<String>("");
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

    public fun addFolderToSync(uri: Uri){
        viewModelScope.launch(Dispatchers.IO){
            try {
                folderDao.addFolder(Folder(uri.toString(), null))
                refreshFolders()
            } catch (e: SQLiteConstraintException){
                _error.value = e.toString();
            }
        }
    }
}
