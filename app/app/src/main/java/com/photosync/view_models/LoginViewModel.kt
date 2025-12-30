package com.photosync.view_models

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.photosync.api.ApiHandler
import com.photosync.api.LoginStatus.*
import com.photosync.database.LocalDatabase
import com.photosync.database.AppSettings
import com.photosync.database.AppSettingsDao
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

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

class LoginViewModel(
    localDatabase: LocalDatabase,
    private var apiHandler: ApiHandler) : ViewModel() {
    private val settingsDao: AppSettingsDao = localDatabase.appSettingsDao()
    private val _loginStatus = MutableStateFlow(LoginStatus(error="", pending = false))
    private val _window = MutableStateFlow(Window.Load)
    val loginStatus: StateFlow<LoginStatus> = _loginStatus.asStateFlow()
    val window: StateFlow<Window> = _window.asStateFlow()
    var appSettings: AppSettings? = null

    fun load(){
        viewModelScope.launch(Dispatchers.IO) {
            appSettings = settingsDao.getSettings()
            _window.value = Window.Login
        }
    }

    fun login(server: String, username: String, password: String){
        _loginStatus.value = LoginStatus(error="", pending = true)
        viewModelScope.launch(Dispatchers.IO) {
            try{
                val loginStatus = apiHandler.login(server, username, password)
                when (loginStatus) {
                    SUCCESS -> {
                        _loginStatus.value = LoginStatus(error="", pending = false)
                        settingsDao.clearSettings()
                        settingsDao.insertSettings(AppSettings(server, username))
                        _window.value = Window.Sync
                    }
                    INVALID_CREDENTIALS -> {
                        _loginStatus.value = LoginStatus(error="Invalid credentials", pending = false)
                    }
                    ERROR -> {
                        _loginStatus.value = LoginStatus(error="Something went wrong", pending = false)
                    }
                }
            }
            catch(e: Exception){
                _loginStatus.value = LoginStatus(error="$e", pending = false)
            }
        }
    }
}