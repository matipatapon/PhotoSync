package com.photosync.database

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity
data class AppSettings(
    @PrimaryKey val uid: Int,
    @ColumnInfo(name = "server") val server: String,
    @ColumnInfo(name = "login") val login: String,
)