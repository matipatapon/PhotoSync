package com.photosync.database

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity
data class UploadedFile(
    @PrimaryKey() val uri: String,
    @ColumnInfo(name = "modificationDate") val modificationDate: Long,
)
