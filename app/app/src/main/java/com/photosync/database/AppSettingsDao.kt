package com.photosync.database

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import androidx.room.Update

@Dao
interface AppSettingsDao{
    @Query("SELECT * FROM appsettings")
    fun getSettings(): AppSettings?

    @Insert
    fun insertSettings(appSettings: AppSettings)

    @Query("DELETE FROM appsettings")
    fun clearSettings()
}
