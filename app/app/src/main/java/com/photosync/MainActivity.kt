package com.photosync

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.input.TextFieldLineLimits
import androidx.compose.foundation.text.input.rememberTextFieldState
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SecureTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import com.photosync.ui.theme.AppTheme
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.unit.dp

class MainActivity : ComponentActivity() {
    private var mainViewModel: MainViewModel? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        mainViewModel = MainViewModel(applicationContext)

        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AppTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    View(
                        innerPadding
                    )
                }
            }
        }
    }

    @Composable
    fun View(innerPadding: PaddingValues){
        val window by mainViewModel!!.window.collectAsState()
        Column(
            modifier = Modifier.padding(50.dp).fillMaxSize(),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
            content = {
                when (window) {
                    Window.Load -> {
                        mainViewModel!!.load()
                    }

                    Window.Login -> {
                        LoginForm()
                    }
                    Window.Sync -> {
                    }
                }
            }
        )
    }

    @Composable
    fun LoginForm() {
        val appSettings = mainViewModel!!.appSettings
        var initialServer = ""
        var initialUsername = ""
        if(appSettings != null){
            initialServer = appSettings.server
            initialUsername = appSettings.login
        }
        val server = rememberTextFieldState(initialText = initialServer)
        val username = rememberTextFieldState(initialText = initialUsername)
        val password = rememberTextFieldState(initialText = "")
        val loginStatus by mainViewModel!!.loginStatus.collectAsState()
        TextField(
            state = server,
            placeholder = { Text("server") },
            lineLimits = TextFieldLineLimits.SingleLine,
            modifier = Modifier.fillMaxWidth(),
            enabled = !loginStatus.isPending()
        )
        TextField(
            state = username,
            placeholder = { Text("login") },
            lineLimits = TextFieldLineLimits.SingleLine,
            modifier = Modifier.fillMaxWidth(),
            enabled = !loginStatus.isPending()
        )
        SecureTextField(
            state = password,
            placeholder = { Text("password") },
            modifier = Modifier.fillMaxWidth(),
            enabled = !loginStatus.isPending(),
        )
        Text(loginStatus.getError())
        Button(
            onClick = {
                mainViewModel!!.login(server.text.toString(), username.text.toString(), password.text.toString())
            },
            enabled = !loginStatus.isPending(),
            content = {
                Text("Login")
            },
            modifier = Modifier.fillMaxWidth(),
        )
    }
}
