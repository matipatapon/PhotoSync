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

enum class LoginState{
    TYPING, WORKING, SUCCESS, ERROR
}

class MainViewModel : ViewModel() {
    private final val client: OkHttpClient = OkHttpClient();
    private val _loginState = MutableStateFlow(LoginState.TYPING)
    val loginState: StateFlow<LoginState> = _loginState.asStateFlow()

    fun login(server: String, username: String, password: String){
        viewModelScope.launch(Dispatchers.IO) {
            _loginState.value = LoginState.WORKING
            // http://192.168.68.60:8080/v1/login
            val json = """
                {
                    "username": "user",
                    "password": "password"
                }
            """.trimIndent()
            val request = Request.Builder()
                .url(server)
                .post(RequestBody.create(MediaType.parse("application/json; charset=utf-8"), json))
                .build();
            try {
                val response = client.newCall(request).execute()
                _loginState.value = LoginState.SUCCESS
            } catch(e: Exception){
                _loginState.value = LoginState.ERROR
            }

        }
    }

}