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

class MainActivity : ComponentActivity() {
    private val mainViewModel: MainViewModel = MainViewModel()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AppTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    LoginForm(
                        innerPadding
                    )
                }
            }
        }
    }

    @Composable
    fun LoginForm(innerPadding: PaddingValues) {
        val server = rememberTextFieldState(initialText = "")
        val username = rememberTextFieldState(initialText = "")
        val password = rememberTextFieldState(initialText = "")
        val loginStatus by mainViewModel.loginState.collectAsState()
        var errorMsg = ""
        if(loginStatus == LoginState.ERROR){
            errorMsg = "Something went wrong"
        } else if(loginStatus == LoginState.INVALID_CREDENTIALS){
            errorMsg = "Invalid credentials"
        }
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
                    modifier = Modifier.fillMaxWidth()
                )
                TextField(
                    state = username,
                    placeholder = { Text("login") },
                    lineLimits = TextFieldLineLimits.SingleLine,
                    modifier = Modifier.fillMaxWidth()
                )
                SecureTextField(
                    state = password,
                    placeholder = { Text("password") },
                    modifier = Modifier.fillMaxWidth(),
                )
                Text(errorMsg)
                Button(
                    onClick = {
                        mainViewModel.login(server.text.toString(), username.text.toString(), password.text.toString())
                        password.setTextAndPlaceCursorAtEnd("haha")
                    },
                    enabled = loginStatus != LoginState.WORKING,
                    content = {
                        Text("Login")
                    },
                    modifier = Modifier.fillMaxWidth(),
                )
            })
    }
}
