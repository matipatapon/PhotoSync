package com.photosync.database

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import androidx.room.Update

@Dao
interface UploadedFileDao {
    @Query("SELECT * FROM uploadedfile WHERE uri = :uri AND modificationDate = :modificationDate")
    fun getUploadedFile(uri: String, modificationDate: Long): UploadedFile?

    @Insert
    fun addUploadedFile(uploadedFile: UploadedFile)

    @Query("DELETE FROM uploadedfile WHERE folderId = :folderId")
    fun removeAllUploadedFilesInGivenFolderFromCache(folderId: Long)
}
