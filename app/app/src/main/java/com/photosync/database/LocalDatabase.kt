package com.photosync.Database

import androidx.room.Database
import androidx.room.RoomDatabase

@Database(entities = [AppSettings::class], version = 1)
abstract class LocalDatabase : RoomDatabase(){
    abstract fun appSettingsDao(): AppSettingsDao
}