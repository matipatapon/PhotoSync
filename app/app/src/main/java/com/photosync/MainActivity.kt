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
import androidx.compose.foundation.text.input.setTextAndPlaceCursorAtEnd
import androidx.compose.foundation.text.input.setTextAndSelectAll
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SecureTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.mutableStateOf
import androidx.compose.ui.Modifier
import com.photosync.ui.theme.AppTheme
import androidx.compose.runtime.remember
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.unit.dp
import androidx.room.Room
import com.photosync.daos.AppSettingsDao
import com.photosync.databases.LocalDatabase

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
        if(window == Window.Login){
            LoginForm()
        }
        else if(window == Window.Sync){
        }
    }

    @Composable
    fun LoginForm() {
        val server = rememberTextFieldState(initialText = "")
        val username = rememberTextFieldState(initialText = "")
        val password = rememberTextFieldState(initialText = "")
        val loginStatus by mainViewModel!!.loginStatus.collectAsState()
        Column(
            Modifier
                .fillMaxSize()
                .padding(50.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
            content = {
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
                        password.setTextAndPlaceCursorAtEnd("haha")
                    },
                    enabled = !loginStatus.isPending(),
                    content = {
                        Text("Login")
                    },
                    modifier = Modifier.fillMaxWidth(),
                )
            })
    }
}
