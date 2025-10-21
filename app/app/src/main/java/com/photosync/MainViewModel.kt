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
    WORKING,
    SUCCESS,
    INVALID_CREDENTIALS,
    ERROR
}

class MainViewModel : ViewModel() {
    private final val client: OkHttpClient = OkHttpClient();
    private val _loginState = MutableStateFlow<LoginState?>(null)
    val loginState: StateFlow<LoginState?> = _loginState.asStateFlow()
    var token: String? = null;

    fun login(server: String, username: String, password: String){
        _loginState.value = LoginState.WORKING
        viewModelScope.launch(Dispatchers.IO) {
            // http://192.168.68.60:8080/v1/login
            val payload = """
                {
                    "username": "${username}",
                    "password": "${password}"
                }
            """.trimIndent()
            val request = Request.Builder()
                .url("$server/v1/login")
                .post(RequestBody.create(MediaType.parse("application/json; charset=utf-8"), payload))
                .build();
            try {
                val response = client.newCall(request).execute()
                val responseCode = response.code()
                if(responseCode == 401){
                    _loginState.value = LoginState.INVALID_CREDENTIALS
                } else if(responseCode != 200){
                    _loginState.value = LoginState.ERROR
                } else{
                    token = response.body().toString()
                    _loginState.value = LoginState.SUCCESS
                }
            } catch(e: Exception){
                _loginState.value = LoginState.ERROR
            }

        }
    }

}