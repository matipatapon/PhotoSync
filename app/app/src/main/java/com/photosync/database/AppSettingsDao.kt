package com.photosync.database

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.Query
import androidx.room.Update

@Dao
interface AppSettingsDao{
    @Query("SELECT * FROM appsettings")
    fun getSettings(): AppSettings?

    @Update
    fun updateSettings(appSettings: AppSettings)

    @Insert
    fun insertSettings(appSettings: AppSettings)
}
