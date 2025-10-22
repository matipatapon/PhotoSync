package com.photosync

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
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
    Login,
    Sync
}

class LoginStatus(private var error: String, private var pending: Boolean) {
    public fun getError(): String {
        return this.error
    }

    public fun isPending(): Boolean{
        return this.pending
    }
}

class MainViewModel : ViewModel() {
    private final val client: OkHttpClient = OkHttpClient();
    private val _loginStatus = MutableStateFlow<LoginStatus>(LoginStatus(error="", pending = false));
    private val _window = MutableStateFlow<Window>(Window.Login);
    val loginStatus: StateFlow<LoginStatus> = _loginStatus.asStateFlow()
    val window: StateFlow<Window> = _window.asStateFlow()
    var token: String? = null;

    fun login(server: String, username: String, password: String){
        viewModelScope.launch(Dispatchers.IO) {
            _loginStatus.value = LoginStatus(error="", pending = true)
            // http://192.168.68.60:8080/v1/login
            val payload = """
                {
                    "username": "${username}",
                    "password": "${password}"
                }
            """.trimIndent()
            try {
            val request = Request.Builder()
                .url("$server/v1/login")
                .post(RequestBody.create(MediaType.parse("application/json; charset=utf-8"), payload))
                .build();
                val response = client.newCall(request).execute()
                val responseCode = response.code()
                if(responseCode == 401){
                    _loginStatus.value = LoginStatus(error="Invalid credentials", pending = false)
                } else if(responseCode != 200){
                    _loginStatus.value = LoginStatus(error="Something went wrong", pending = false)
                } else{
                    token = response.body().toString()
                    _loginStatus.value = LoginStatus(error="", pending = false)
                    _window.value = Window.Sync
                }
            } catch(e: Exception){
                _loginStatus.value = LoginStatus(error="$e", pending = false)
            }
        }
    }

}