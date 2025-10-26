package com.photosync.database

import androidx.room.Database
import androidx.room.RoomDatabase

@Database(entities = [AppSettings::class, Folder::class], version = 1)
abstract class LocalDatabase : RoomDatabase(){
    abstract fun appSettingsDao(): AppSettingsDao
    abstract fun folderDao(): FolderDao
}
