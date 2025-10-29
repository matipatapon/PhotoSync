package com.photosync.database

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.Query
import androidx.room.Update

@Dao
interface FolderDao {
    @Query("SELECT * FROM folder")
    fun getFolders(): List<Folder>

    @Insert
    fun addFolder(folder: Folder)

    @Update
    fun updateFolder(folder: Folder)
}