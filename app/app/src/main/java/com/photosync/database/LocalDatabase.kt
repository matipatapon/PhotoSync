package com.photosync.database

import androidx.room.Database
import androidx.room.RoomDatabase

@Database(entities = [AppSettings::class, Folder::class, UploadedFile::class], version = 1, exportSchema = false)
abstract class LocalDatabase : RoomDatabase(){
    abstract fun appSettingsDao(): AppSettingsDao
    abstract fun folderDao(): FolderDao
    abstract fun uploadedFileDao(): UploadedFileDao
}
