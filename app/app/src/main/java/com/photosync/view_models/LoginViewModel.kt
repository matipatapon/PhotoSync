package com.photosync.view_models

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.photosync.database.LocalDatabase
import com.photosync.database.AppSettings
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import okhttp3.MediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody

enum class Window {
    Load,
    Login,
    Sync
}

class LoginStatus(private var error: String, private var pending: Boolean) {
    fun getError(): String {
        return this.error
    }

    fun isPending(): Boolean{
        return this.pending
    }
}

class LoginViewModel(private var localDatabase: LocalDatabase) : ViewModel() {
    private val client: OkHttpClient = OkHttpClient()
    private val _loginStatus = MutableStateFlow(LoginStatus(error="", pending = false))
    private val _window = MutableStateFlow(Window.Load)
    val loginStatus: StateFlow<LoginStatus> = _loginStatus.asStateFlow()
    val window: StateFlow<Window> = _window.asStateFlow()
    var token: String? = null
    var appSettings: AppSettings? = null

    fun load(){
        viewModelScope.launch(Dispatchers.IO) {
            val dao = localDatabase.appSettingsDao()
            appSettings = dao.getSettings()
            _window.value = Window.Login
        }
    }

    fun login(server: String, username: String, password: String){
        viewModelScope.launch(Dispatchers.IO) {
            _loginStatus.value = LoginStatus(error="", pending = true)
            val payload = """
                {
                    "username": "$username",
                    "password": "$password"
                }
            """.trimIndent()
            try {
                val request = Request.Builder()
                    .url("$server/v1/login")
                    .post(RequestBody.create(MediaType.parse("application/json; charset=utf-8"), payload))
                    .build()
                val response = client.newCall(request).execute()
                val responseCode = response.code()
                if(responseCode == 401){
                    _loginStatus.value = LoginStatus(error="Invalid credentials", pending = false)
                } else if(responseCode != 200){
                    _loginStatus.value = LoginStatus(error="Something went wrong", pending = false)
                } else{
                    val dao = localDatabase.appSettingsDao()
                    val currentSettings = dao.getSettings()
                    val newAppSettings = AppSettings(1, server, username)
                    if(currentSettings == null){
                        dao.insertSettings(newAppSettings)
                    }else{
                        dao.updateSettings(newAppSettings)
                    }
                    token = response.body().string()

                    appSettings = newAppSettings
                    _loginStatus.value = LoginStatus(error="", pending = false)
                    _window.value = Window.Sync
                }
            } catch(e: Exception){
                _loginStatus.value = LoginStatus(error="$e", pending = false)
            }
        }
    }

}