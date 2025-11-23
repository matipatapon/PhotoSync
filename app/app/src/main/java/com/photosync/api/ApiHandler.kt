package com.photosync.api

import okhttp3.MediaType
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.MultipartBody
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import java.util.logging.Logger
import kotlin.math.log

enum class LoginStatus{
    SUCCESS,
    INVALID_CREDENTIALS,
    ERROR
}

enum class UploadStatus{
    SUCCESS,
    ERROR
}

class ApiHandler {
    private val client: OkHttpClient = OkHttpClient()
    private var server: String? = null
    private val logger = Logger.getLogger(this.javaClass.name)
    private var token: String? = null

    fun login(server: String, username: String, password: String): LoginStatus{
        try {
            this.server = server
            val payload = """
                {
                    "username": "$username",
                    "password": "$password"
                }
            """.trimIndent()
            val requestBody = payload.toRequestBody("application/json; charset=utf-8".toMediaTypeOrNull())
            val request = Request.Builder()
                .url("$server/v1/login")
                .post(requestBody)
                .build()
            val response = client.newCall(request).execute()
            val responseCode = response.code
            if(responseCode == 200){
                token = response.body.string()
                return LoginStatus.SUCCESS
            } else if(responseCode == 401){
                return LoginStatus.INVALID_CREDENTIALS
            }
        } catch(e: Exception){
            logger.warning("Error occurred: '${e.toString()}'")
        }
        return LoginStatus.ERROR
    }

    fun uploadFile(file: ByteArray, filename: String, lastModified: String): UploadStatus{
        try {
            val requestBody: MultipartBody = MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("filename", filename)
                .addFormDataPart("modification_date", lastModified)
                .addFormDataPart(
                    "file",
                    "",
                    file.toRequestBody("image/jpeg".toMediaTypeOrNull())
                )
                .build()
            val request = Request.Builder()
                .header("Authorization", token!!)
                .url("$server/v1/upload")
                .post(requestBody)
                .build()
            val response = client.newCall(request).execute()
            val responseCode = response.code
            if (responseCode == 200 || responseCode == 402) {
                logger.info("Uploaded file<filename={$filename} lastModified={$lastModified}>")
                return UploadStatus.SUCCESS
            } else{
                logger.warning("Failed to upload file<filename={$filename} lastModified={$lastModified}>")
            }
        } catch (e: Exception){
            logger.warning("Error occurred: '${e.toString()}'")
        }
        return UploadStatus.ERROR
    }
}