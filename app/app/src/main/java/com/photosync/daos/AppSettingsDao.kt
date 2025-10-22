package com.photosync.daos

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.Query
import com.photosync.entities.AppSettings

@Dao
interface AppSettingsDao{
    @Query("SELECT * FROM appsettings")
    fun getSettings(): AppSettings?

    @Insert
    fun insertSettings(appSettings: AppSettings)

    @Delete
    fun deleteSettings(appSettings: AppSettings)

    fun updateSettings(appSettings: AppSettings){
        deleteSettings(appSettings)
        insertSettings(appSettings)
    }

}