package com.photosync.databases

import androidx.room.Database
import androidx.room.RoomDatabase
import com.photosync.daos.AppSettingsDao
import com.photosync.entities.AppSettings

@Database(entities = [AppSettings::class], version = 1)
abstract class LocalDatabase : RoomDatabase(){
    abstract fun appSettingsDao(): AppSettingsDao
}