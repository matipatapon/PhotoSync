package com.photosync.database

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity
data class Folder(
    @PrimaryKey var uri:String,
    @ColumnInfo var lastSync: Long?
);
